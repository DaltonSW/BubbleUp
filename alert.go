package bubbleup

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	InfoAlertKey  = "Info"
	WarnAlertKey  = "Warn"
	ErrorAlertKey = "Error"
	DebugAlertKey = "Debug"

	infoNerdSymbol  = " "
	warnNerdSymbol  = "󱈸 "
	errorNerdSymbol = "󰬅 "
	debugNerdSymbol = "󰃤 "

	InfoUniSymbol    = "(i)"
	WarningUniSymbol = "(!)"
	ErrorUniSymbol   = "[!]"
	DebugUniSymbol   = "(?)"

	NotifWidth = 40

	LerpIncrement = 0.18
)

// Constant colors used for included alert types
var (
	infoColor, _  = colorful.Hex("#00FF00")
	warnColor, _  = colorful.Hex("#FFFF00")
	errorColor, _ = colorful.Hex("#FF0000")
	debugColor, _ = colorful.Hex("#FF00FF")

	backColor, _ = colorful.Hex("#000000")

	baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Width(NotifWidth)
)

type AlertLevel int

type AlertMsg struct {
	alertKey string
	msg      string
	dur      time.Duration

	// TODO:
	// animation: how the notification should appear and disappear
	// location: where on the screen it should appear
	// style: Mimic nvim.notify's style options perhaps?
}

// TODO:
// type CustomNotifMsg struct {
// 	msg    string
// 	dur    time.Duration
// 	symbol string
// 	color  lipgloss.Color
// }

func (m AlertModel) newNotif(key, msg string, dur time.Duration) *alert {
	if msg == "" || key == "" {
		return nil
	}

	alertDef, ok := m.alertTypes[key]

	if !ok {
		return nil
	}

	return &alert{
		message:     msg,
		deathTime:   time.Now().Add(dur),
		symbol:      alertDef.Symbol,
		foreColor:   alertDef.ForeColor,
		style:       alertDef.Style,
		width:       NotifWidth,
		curLerpStep: 0.3,
	}

}

type alert struct {
	message   string
	deathTime time.Time
	symbol    string
	foreColor colorful.Color
	style     lipgloss.Style
	width     int

	curLerpStep float64

	// animation
	// location
}

func (n *alert) render() string {
	newColor := backColor.BlendLab(n.foreColor, n.curLerpStep)
	lipColor := lipgloss.Color(newColor.Hex())
	newStyle := baseStyle.Foreground(lipColor).BorderForeground(lipColor)
	return newStyle.Render(fmt.Sprintf("%v %v", n.symbol, n.message))
}

// Region: Model stuff

type AlertDefinition struct {
	Key       string
	ForeColor colorful.Color
	Style     lipgloss.Style
	Symbol    string
	// DefaultDur time.Duration
	// DefaultPos
	// Default
}

func (m AlertModel) NewAlertCmd(alertType, message string) tea.Cmd {
	return func() tea.Msg {
		return AlertMsg{alertKey: alertType, msg: message, dur: time.Second * 2}
	}
}

func (m AlertModel) RegisterNewAlertType(definition AlertDefinition) {
	if m.alertTypes == nil {
		m.alertTypes = make(map[string]AlertDefinition)
	}

	m.alertTypes[definition.Key] = definition
}

func (m AlertModel) registerDefaultAlertTypes() {
	infoDef := AlertDefinition{
		Key:       "Info",
		Symbol:    infoNerdSymbol,
		ForeColor: infoColor,
	}

	m.RegisterNewAlertType(infoDef)

	warnDef := AlertDefinition{
		Key:       "Warn",
		Symbol:    warnNerdSymbol,
		ForeColor: warnColor,
	}

	m.RegisterNewAlertType(warnDef)

	errorDef := AlertDefinition{
		Key:       "Error",
		Symbol:    errorNerdSymbol,
		ForeColor: errorColor,
	}

	m.RegisterNewAlertType(errorDef)

	debugDef := AlertDefinition{
		Key:       "Debug",
		Symbol:    debugNerdSymbol,
		ForeColor: debugColor,
	}

	m.RegisterNewAlertType(debugDef)
}
