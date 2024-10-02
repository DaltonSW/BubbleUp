package bubbleup

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/ansi"
)

type AlertModel struct {
	activeAlert *alert

	width int
}

// TODO: Set defaults for position and duration
func New() *AlertModel {
	return &AlertModel{
		activeAlert: nil,
	}
}

func (m AlertModel) Init() tea.Cmd {
	return tickCmd()
}

func (m AlertModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.activeAlert != nil && m.activeAlert.deathTime.Before(time.Time(msg)) {
			m.activeAlert = nil
		}

		return m, tickCmd()
	case AlertMsg:
		m.activeAlert = newNotif(msg.msg, msg.level, msg.dur)
	}

	return m, nil
}

func (m AlertModel) View() string {
	return ""
}

// Used the following code as a reference:
//
//	https://github.com/charmbracelet/lipgloss/pull/102/commits/a075bfc9317152e674d661a2cdfe58144306e77a
func (m AlertModel) Render(content string) string {
	if m.activeAlert == nil {
		return content
	}

	notifString := m.activeAlert.render()

	notifSplit, _ := getLines(notifString)
	contentSplit, _ := getLines(content)

	notifHeight := len(notifSplit)
	contentHeight := len(contentSplit)

	// posX, posY := 0, 0

	var builder strings.Builder

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

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
