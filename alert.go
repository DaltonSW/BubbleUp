package bubbleup

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// Alert keys for the included alert types.
const (
	InfoKey  = "Info"
	WarnKey  = "Warn"
	ErrorKey = "Error"
	DebugKey = "Debug"
)

// Symbols used by the included alert types.
// To use the NerdFont symbols, you must be using a NerdFont,
// which can be obtained from https://www.nerdfonts.com/.
// If you want to use the default non-NerdFont symbols, pass
// false into the useNerdFont parameter when creating your alert model.
const (
	InfoNerdSymbol  = " "
	WarnNerdSymbol  = "󱈸 "
	ErrorNerdSymbol = "󰬅 "
	DebugNerdSymbol = "󰃤 "

	InfoASCIIPrefix    = "(i)"
	WarningASCIIPrefix = "(!)"
	ErrorASCIIPrefix   = "[!!]"
	DebugASCIIPrefix   = "(?)"

	InfoUnicodePrefix    = "\u24D8 " // Trailing space is intentional
	WarningUnicodePrefix = "\u26A0"
	ErrorUnicodePrefix   = "\u2718"
	DebugUnicodePrefix   = "\u003F"

	// Deprecated: use InfoASCIIPrefix instead.
	InfoUniPrefix = InfoASCIIPrefix

	// Deprecated: use WarningASCIIPrefix instead.
	WarningUniPrefix = WarningASCIIPrefix

	// Deprecated: use ErrorASCIIPrefix instead.
	ErrorUniPrefix = ErrorASCIIPrefix

	// Deprecated: use DebugASCIIPrefix instead.
	DebugUniPrefix = DebugASCIIPrefix
)

// Defaults used by the notification rendering.
const (
	DefaultLerpIncrement = 0.18
)

// Colors used by the included alert types.
const (
	InfoColor  = "#00FF00"
	WarnColor  = "#FFFF00"
	ErrorColor = "#FF0000"
	DebugColor = "#FF00FF"
	BackColor  = "#000000"
)

// Constant colors and stylings used for included alert types.
// Ignoring errors because we are using hardcoded values
var (
	infoColor, _  = colorful.Hex(InfoColor)
	warnColor, _  = colorful.Hex(WarnColor)
	errorColor, _ = colorful.Hex(ErrorColor)
	debugColor, _ = colorful.Hex(DebugColor)
	backColor, _  = colorful.Hex(BackColor)

	baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
)

var parsedColors = map[string]colorful.Color{
	InfoColor:  infoColor,
	WarnColor:  warnColor,
	ErrorColor: errorColor,
	DebugColor: debugColor,
	BackColor:  backColor,
}

// alertMsg is the tea.Msg used to activate a notification
type alertMsg struct {
	alertKey string
	msg      string
	dur      time.Duration

	// TODO:
	// animation: how the notification should appear and disappear
	// style: Mimic nvim.notify's style options perhaps?
}

func (m AlertModel) newNotify(key, msg string, dur time.Duration) *alert {
	if msg == "" || key == "" {
		return nil
	}

	alertDef, ok := m.alertTypes[key]

	if !ok {
		return nil
	}

	foreColor, ok := parsedColors[alertDef.ForeColor]
	if !ok {
		// Can safely discard error because we validated the color
		// when registering the alert defition
		foreColor, _ = colorful.Hex(alertDef.ForeColor)
	}

	return &alert{
		message:     msg,
		deathTime:   time.Now().Add(dur),
		prefix:      alertDef.Prefix,
		foreColor:   foreColor,
		style:       alertDef.Style,
		width:       m.width,
		minWidth:    m.minWidth,
		curLerpStep: 0.3,
		position:    m.position,
	}

}

// alert represents an instance of an actual alert, including
// all information needed to render and destroy itself
type alert struct {
	message   string
	deathTime time.Time
	prefix    string
	foreColor colorful.Color
	style     lipgloss.Style
	width     int
	minWidth  int

	curLerpStep float64
	position    Position
}

// render will render the given alert based on its values
// Returns the string representation of the alert, ready to be
// overlayed onto the main content.
func (n *alert) render() string {
	newColor := backColor.BlendLab(n.foreColor, n.curLerpStep)
	lipColor := lipgloss.Color(newColor.Hex())

	// Calculate actual width based on minWidth setting
	actualWidth := n.width // default to max/fixed width

	if n.minWidth > 0 {
		// Dynamic mode: measure message width
		messageText := fmt.Sprintf("%v %v", n.prefix, n.message)

		// Get the width of the message text itself
		messageWidth := lipgloss.Width(messageText)

		// Account for extra space needed, determined imperically
		messageWidth += 3

		// Clamp between min and max
		if messageWidth < n.minWidth {
			actualWidth = n.minWidth
		} else if messageWidth > n.width {
			actualWidth = n.width
		} else {
			actualWidth = messageWidth
		}
	}

	newStyle := baseStyle.
		Foreground(lipColor).
		BorderForeground(lipColor).
		Width(actualWidth).
		Padding(0, 1)

	// Compute width available for text inside border+padding.
	textWidth := actualWidth - 2
	if textWidth < 1 {
		textWidth = 1
	}

	content := hangingWrap(n.prefix, n.message, textWidth)
	return newStyle.Render(content)
}

// Region: Model stuff

// AlertDefinition is all the information needed to register a new alert type.
type AlertDefinition struct {
	// (Req) Unique key used to refer to an alert type
	Key string

	// (Req) Hex code of the color you want your alert to be
	ForeColor string

	// (Opt) lipgloss.Style used to render the alert
	Style lipgloss.Style

	// (Opt) String used to prefix the alert message
	Prefix string

	// DefaultDur time.Duration
	// DefaultPos
	// Default
}

// NewAlertCmd will construct and return the tea.Cmd needed to trigger
// an alert. This should be called in your Update() function, and the
// returned tea.Cmd should be batched into your return.
func (m AlertModel) NewAlertCmd(alertType, message string) tea.Cmd {
	return func() tea.Msg {
		return alertMsg{alertKey: alertType, msg: message, dur: time.Second * m.duration}
	}
}

// RegisterNewAlertType will registery a new alert type based on the provided
// AlertDefintion. This can also be used to overwrite the provided defaults
// by providing an AlertDefintion with one of the default keys.
func (m AlertModel) RegisterNewAlertType(definition AlertDefinition) {
	_, err := colorful.Hex(definition.ForeColor)
	if err != nil {
		log.Fatal(err)
		return
	}

	if m.alertTypes == nil {
		m.alertTypes = make(map[string]AlertDefinition)
	}

	m.alertTypes[definition.Key] = definition
}

var unicodePrefixes = map[string]string{
	"Info":  InfoUnicodePrefix,
	"Warn":  WarningUnicodePrefix,
	"Error": ErrorUnicodePrefix,
	"Debug": DebugUnicodePrefix,
}

// Registers all the alert types that ship with BubbleUp by out of the box.
func (m AlertModel) registerDefaultAlertTypes() {
	var (
		infoPref  string
		warnPref  string
		errPref   string
		debugPref string
	)

	if m.useNerdFont {
		infoPref = InfoNerdSymbol
		warnPref = WarnNerdSymbol
		errPref = ErrorNerdSymbol
		debugPref = DebugNerdSymbol
	} else {
		infoPref = InfoASCIIPrefix
		warnPref = WarningASCIIPrefix
		errPref = ErrorASCIIPrefix
		debugPref = DebugASCIIPrefix
	}

	infoDef := AlertDefinition{
		Key:       "Info",
		Prefix:    infoPref,
		ForeColor: InfoColor,
	}

	m.RegisterNewAlertType(infoDef)

	warnDef := AlertDefinition{
		Key:       "Warn",
		Prefix:    warnPref,
		ForeColor: WarnColor,
	}

	m.RegisterNewAlertType(warnDef)

	errorDef := AlertDefinition{
		Key:       "Error",
		Prefix:    errPref,
		ForeColor: ErrorColor,
	}

	m.RegisterNewAlertType(errorDef)

	debugDef := AlertDefinition{
		Key:       "Debug",
		Prefix:    debugPref,
		ForeColor: DebugColor,
	}

	m.RegisterNewAlertType(debugDef)
}
