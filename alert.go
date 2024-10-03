package bubbleup

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// InfoUniSymbol    = ""
	// WarningUniSymbol = ""
	// ErrorUniSymbol   = ""
	// DebugUniSymbol   = "🐛 "

	infoColor  = lipgloss.Color("#00FF00")
	warnColor  = lipgloss.Color("#FFFF00")
	errorColor = lipgloss.Color("#FF0000")
	debugColor = lipgloss.Color("#FF00FF")

	NotifWidth = 40
)

// func InfoAlertCmd(message string) tea.Cmd {
// 	return func() tea.Msg {
// 		return AlertMsg{msg: message, level: InfoLevel, dur: time.Second * 2}
// 	}
// }

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
		message:   msg,
		deathTime: time.Now().Add(dur),
		symbol:    alertDef.Symbol,
		style:     alertDef.Style,
		width:     NotifWidth,
	}

}

type alert struct {
	message   string
	deathTime time.Time
	symbol    string
	style     lipgloss.Style
	width     int
}

func (n *alert) render() string {
	return n.style.Render(fmt.Sprintf("%v %v", n.symbol, n.message))
}

// Region: Model stuff

type AlertDefinition struct {
	Key    string
	Style  lipgloss.Style
	Symbol string
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
	m.alertTypes[definition.Key] = definition

}

func (m AlertModel) registerDefaultAlertTypes() {
	infoDef := AlertDefinition{
		Key:    "Info",
		Symbol: infoNerdSymbol,
		Style:  lipgloss.NewStyle().Foreground(infoColor).Border(lipgloss.RoundedBorder()).BorderForeground(infoColor),
	}

	m.RegisterNewAlertType(infoDef)

	warnDef := AlertDefinition{
		Key:    "Warn",
		Symbol: warnNerdSymbol,
		Style:  lipgloss.NewStyle().Foreground(warnColor).Border(lipgloss.RoundedBorder()).BorderForeground(warnColor),
	}

	m.RegisterNewAlertType(warnDef)

	errorDef := AlertDefinition{
		Key:    "Error",
		Symbol: errorNerdSymbol,
		Style:  lipgloss.NewStyle().Foreground(errorColor).Border(lipgloss.RoundedBorder()).BorderForeground(errorColor),
	}

	m.RegisterNewAlertType(errorDef)

	debugDef := AlertDefinition{
		Key:    "Debug",
		Symbol: debugNerdSymbol,
		Style:  lipgloss.NewStyle().Foreground(debugColor).Border(lipgloss.RoundedBorder()).BorderForeground(debugColor),
	}

	m.RegisterNewAlertType(debugDef)
}
