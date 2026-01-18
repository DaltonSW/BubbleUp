package bubbleup

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/ansi"
)

// AlertModel maintains a list of alert types, and facilitates the display and
// clearing of alerts, as well as registering custom alerts and overrides to defaults
// Note: duration is measured in seconds
// Note: width behavior depends on minWidth:
//   - minWidth == 0 (default): width is fixed width
//   - minWidth > 0: width is max width, minWidth is minimum, actual width varies with message length
type AlertModel struct {
	useNerdFont      bool
	useUnicodePrefix bool
	allowEscToClose  bool
	alertTypes       map[string]AlertDefinition
	activeAlert      *alert
	width            int
	minWidth         int
	duration         time.Duration
	position         Position
}

// TODO: Set defaults for duration

// NewAlertModel creates and returns a new AlertModel, initialized with default alert types
func NewAlertModel(width int, useNerdFont bool, duration time.Duration) *AlertModel {
	model := &AlertModel{
		activeAlert: nil,
		width:       width,
		minWidth:    0,
		useNerdFont: useNerdFont,
		alertTypes:  make(map[string]AlertDefinition),
		duration:    duration,
		position:    TopLeftPosition,
	}

	model.registerDefaultAlertTypes()

	return model
}

// WithPosition returns a new AlertModel with the specified position.
// This is an immutable operation that returns a copy with the updated position.
func (m AlertModel) WithPosition(pos Position) AlertModel {
	m.position = pos
	return m
}

// WithMinWidth returns a new AlertModel with dynamic width enabled.
// When minWidth > 0, the notification width will vary between minWidth and width (max)
// based on the actual message length. This is an immutable operation.
// If min > width, it will be clamped to width.
func (m AlertModel) WithMinWidth(min int) AlertModel {
	if min > m.width {
		min = m.width // clamp to max
	}
	m.minWidth = min
	return m
}

// WithUnicodePrefix switches the AlertModule to use Unicode fonts
func (m AlertModel) WithUnicodePrefix() AlertModel {
	m.useNerdFont = false
	m.useUnicodePrefix = true //
	alertTypes := make(map[string]AlertDefinition, len(m.alertTypes))
	for name, alertType := range m.alertTypes {
		alertType.Prefix = unicodePrefixes[alertType.Key]
		alertTypes[name] = alertType
	}
	m.alertTypes = alertTypes
	return m
}

func (m AlertModel) WithAllowEscToClose() AlertModel {
	m.allowEscToClose = true
	return m
}

// Init required as part of BubbleTea Model interface
func (m AlertModel) Init() tea.Cmd {
	return nil
}

// Update takes in a message and returns an associated command to drive model
// functionality. First alertMsg starts the ticking command that causes alert
// refreshing Implemented as part of BubbleTea Model interface
func (m AlertModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case alertMsg:
		m.activeAlert = m.newAlert(msg.alertKey, msg.msg, msg.dur)
		return m, tickCmd() // Start ticking when new alert appears

	case tickMsg: // Check to see if it's time to clear the alert
		if m.activeAlert == nil {
			// No alert, don't tick
			break
		}
		if m.activeAlert.deathTime.Before(time.Time(msg)) {
			// Alert expired, stop ticking
			m.activeAlert = nil
			break
		}
		// Keep ticking while alert is active
		m.activeAlert.curLerpStep += DefaultLerpIncrement
		if m.activeAlert.curLerpStep > 1 {
			m.activeAlert.curLerpStep = 1
		}
		return m, tickCmd()

	case tea.KeyMsg:
		if m.activeAlert == nil {
			break
		}
		if msg.String() != "esc" {
			break
		}
		if !m.allowEscToClose {
			break
		}
		m.activeAlert = nil

	}

	return m, nil
}

// HasActiveAlert allows other models to tell if there is an active already and
// avoid processing an esc key used to clear an alert
func (m AlertModel) HasActiveAlert() bool {
	return m.activeAlert != nil
}

// View doesn't do anything, and it should never be called directly
// Implemented as part of BubbleTea Model interface
func (m AlertModel) View() string {
	return ""
}

// Render takes in the main view content and overlays the model's active alert.
// This function expects you build the entirety of your view's content before calling
// this function. It's recommended for this to be the final call of your model's View().
// Returns a string representation of the content with overlayed alert.
func (m AlertModel) Render(content string) string {
	if m.activeAlert == nil {
		return content
	}

	notifString := m.activeAlert.render()
	notifSplit, notifWidth := getLines(notifString)
	contentSplit, contentWidth := getLines(content)
	notifHeight := len(notifSplit)
	contentHeight := len(contentSplit)

	var builder strings.Builder
	for i := range contentHeight {
		if i > 0 {
			builder.WriteByte('\n')
		}

		line := m.buildLineForPosition(
			contentSplit[i],
			notifSplit,
			i,
			notifHeight,
			contentHeight,
			notifWidth,
			contentWidth,
		)
		builder.WriteString(line)
	}

	return builder.String()
}

// buildLineForPosition determines how to overlay notification on content line based on position
func (m AlertModel) buildLineForPosition(contentLine string, notifLines []string, lineIdx, notifHeight, contentHeight, notifWidth, contentWidth int) string {
	// Determine if notification should appear on this line
	var notifIdx int
	var showNotif bool

	switch m.activeAlert.position {
	case TopLeftPosition, TopCenterPosition, TopRightPosition:
		showNotif = lineIdx < notifHeight
		notifIdx = lineIdx
	case BottomLeftPosition, BottomCenterPosition, BottomRightPosition:
		startLine := contentHeight - notifHeight
		if startLine < 0 {
			startLine = 0
		}
		if lineIdx >= startLine {
			showNotif = true
			notifIdx = lineIdx - startLine
		}
	}

	if !showNotif {
		return contentLine
	}

	notifLine := notifLines[notifIdx]

	// Position-specific overlay logic
	switch m.activeAlert.position {
	case TopLeftPosition, BottomLeftPosition:
		return notifLine + cutLeft(contentLine, notifWidth)

	case TopRightPosition, BottomRightPosition:
		// Calculate overlay position based on max content width for consistent alignment
		keepWidth := contentWidth - notifWidth
		if keepWidth < 0 {
			// Notification is wider than content - just show notification
			return notifLine
		}

		// Check if this specific line is shorter than the overlay position
		contentLineWidth := ansi.PrintableRuneWidth(contentLine)
		if contentLineWidth < keepWidth {
			// Pad the line to reach the overlay position
			padding := strings.Repeat(" ", keepWidth-contentLineWidth)
			return contentLine + padding + notifLine
		}

		// Line is long enough - cut and overlay
		return cutRight(contentLine, keepWidth) + notifLine

	case TopCenterPosition, BottomCenterPosition:
		return m.overlayCenter(contentLine, notifLine, notifWidth, contentWidth)

	default:
		// Fallback to left
		return notifLine + cutLeft(contentLine, notifWidth)
	}
}

// overlayCenter overlays notification at center of content line
// If content is too short, just overlay at position 0
func (m AlertModel) overlayCenter(contentLine, notifLine string, notifWidth, contentWidth int) string {
	leftPad := (contentWidth - notifWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	contentLen := ansi.PrintableRuneWidth(contentLine)

	// If content line is shorter than where notification should start, just overlay at 0
	if contentLen < leftPad {
		return notifLine + cutLeft(contentLine, notifWidth)
	}

	// Extract left portion (before notification)
	left := cutRight(contentLine, leftPad)

	// Extract right portion (after notification)
	rightStart := leftPad + notifWidth
	var right string
	if rightStart < contentLen {
		right = cutLeft(contentLine, rightStart)
	}

	return left + notifLine + right
}

// Timer stuff

// TickMsg is the message that tells the model to assess active alert lifespan.
type tickMsg time.Time

// tickCmd returns a tea Command to initiate a tick.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
