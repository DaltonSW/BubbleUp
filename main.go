package main

import (
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func main() {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	outStyle := lipgloss.NewStyle().Width(width-2).Height(height-2).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#00FFFF")).
		Align(lipgloss.Center, lipgloss.Center)

	content := outStyle.Render("This is a test string, wow look at it go!")

	m := mainModel{
		content: content,
		notif:   *New(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err, _ := p.Run(); err != nil {
		return
		// log.Fatal(err)
	}

}

type mainModel struct {
	notif   Model
	content string
}

func (m mainModel) Init() tea.Cmd {
	return m.notif.Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var outNotif tea.Model
	var outCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			outMsg := NotifMsg{msg: "Info notification!", level: InfoLevel, dur: time.Second * 2}
			outNotif, outCmd = m.notif.Update(outMsg)
			m.notif = outNotif.(Model)
		case "w":
			outMsg := NotifMsg{msg: "Warning notification!", level: WarningLevel, dur: time.Second * 2}
			outNotif, outCmd = m.notif.Update(outMsg)
			m.notif = outNotif.(Model)
		case "e":
			outMsg := NotifMsg{msg: "Error notification!", level: ErrorLevel, dur: time.Second * 2}
			outNotif, outCmd = m.notif.Update(outMsg)
			m.notif = outNotif.(Model)
		case "q":
			return m, tea.Quit
		}

	case tickMsg:
		outNotif, outCmd = m.notif.Update(msg)
		m.notif = outNotif.(Model)
	}

	return m, outCmd
}

func (m mainModel) View() string {
	// return m.content
	return m.notif.Render(m.content)
}
