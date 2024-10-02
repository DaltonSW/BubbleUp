package bubbleup

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/ansi"
)

type Model struct {
	activeNotif *notif
	active      bool

	width int
}

func New() *Model {
	return &Model{
		active:      false,
		activeNotif: nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.active && m.activeNotif.deathTime.Before(time.Time(msg)) {
			m.active = false
			m.activeNotif = nil
		}

		return m, tickCmd()
	case NotifMsg:
		// log.Println("Notif Msg received ", msg)
		m.activeNotif = newNotif(msg.msg, msg.level, msg.dur)
		m.active = true
	}

	return m, nil
}

func (m Model) View() string {
	return ""
}

// Used the following code as a reference:
//
//	https://github.com/charmbracelet/lipgloss/pull/102/commits/a075bfc9317152e674d661a2cdfe58144306e77a
func (m Model) Render(content string) string {
	if !m.active {
		return content
	}

	notifString := m.activeNotif.render()

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

func (m Model) Notify(msg string, level NotifLevel, dur time.Duration) {

	m.activeNotif = newNotif(msg, level, dur)
	m.active = true

}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
