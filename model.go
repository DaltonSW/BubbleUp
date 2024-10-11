package bubbleup

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/ansi"
)

// AlertModel maintains a list of alert types, and facilitates the display and
// clearing of alerts, as well as registering custom alerts and overrides to defaults
type AlertModel struct {
	useNerdFont bool
	alertTypes  map[string]AlertDefinition
	activeAlert *alert
	width       int
}

// TODO: Set defaults for position and duration

// NewAlertModel creates and returns a new AlertModel, initialized with default alert types
func NewAlertModel(width int, useNerdFont bool) *AlertModel {
	model := &AlertModel{
		activeAlert: nil,
		width:       width,
		useNerdFont: useNerdFont,
		alertTypes:  make(map[string]AlertDefinition),
	}

	model.registerDefaultAlertTypes()

	return model
}

// Init starts the ticking command that causes alert refreshing
// Implemented as part of BubbleTea Model interface
func (m AlertModel) Init() tea.Cmd {
	return tickCmd()
}

// Update takes in a message and returns an associated command to drive model functionality
// Implemented as part of BubbleTea Model interface
func (m AlertModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg: // Check to see if it's time to clear the alert
		if m.activeAlert != nil {
			if m.activeAlert.deathTime.Before(time.Time(msg)) {
				m.activeAlert = nil
			} else {
				m.activeAlert.curLerpStep += DefaultLerpIncrement
				if m.activeAlert.curLerpStep > 1 {
					m.activeAlert.curLerpStep = 1
				}
			}
		}

		return m, tickCmd()
	case alertMsg:
		m.activeAlert = m.newNotif(msg.alertKey, msg.msg, msg.dur)
	}

	return m, nil
}

// View doesn't do anything, and it should never be called directly
// Implemented as part of BubbleTea Model interface
func (m AlertModel) View() string {
	return ""
}

// RenderAlert takes in the main view content and overlays the model's active alert.
// This function expects you build the entirety of your view's content before calling
// this function. It's recommended for this to be the final call of your model's View().
// Returns a string representation of the content with overlayed alert.
func (m AlertModel) Render(content string) string {
	if m.activeAlert == nil {
		return content
	}

	notifString := m.activeAlert.render()

	notifSplit, _ := getLines(notifString)
	contentSplit, _ := getLines(content)

	notifHeight := len(notifSplit)
	contentHeight := len(contentSplit)

	var builder strings.Builder

	// NOTE: The current implementation here assumes the notification is
	// in the top left. It will need to be adapted to handle other positions

	for i := range contentHeight {
		if i > 0 {
			// End previous line with a newline
			builder.WriteByte('\n')
		}

		if i >= notifHeight { // If we're past the notifation, render normally
			builder.WriteString(contentSplit[i])
		} else {
			// Add notification line
			notifLine := notifSplit[i]
			builder.WriteString(notifLine)

			notifLen := ansi.PrintableRuneWidth(notifLine)
			remainingContent := cutLeft(contentSplit[i], notifLen)
			builder.WriteString(remainingContent)
		}
	}

	return builder.String()
}

// Timer stuff

// tickMsg is the message that tells the model to assess active alert lifespan.
type tickMsg time.Time

// tickCmd returns a tea Command to initiate a tick.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
