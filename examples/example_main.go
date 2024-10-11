package main

// Only two imports that are required
import (
	tea "github.com/charmbracelet/bubbletea"
	"go.dalton.dog/bubbleup"
)

import (
	"log"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type testModel struct {
	content string
	alert   bubbleup.AlertModel
}

func main() {
	// Do whatever you need to do to get your
	// model looking pretty and initialized.
	content := getTestContent()

	// Create a new alert model and embed it
	// within your model and you're good to go.

	m := testModel{
		content: content,

		// width = 80, useNerdFont = true
		alert: *bubbleup.NewAlertModel(80, true),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err, _ := p.Run(); err != nil {
		return
	}

}

func (m testModel) Init() tea.Cmd {
	// Be sure to return the result of the alert models' Init()
	// If you need to also return one or more commands,
	// be sure to use tea.Batch() to bundle them together.
	return m.alert.Init()
}

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Creating a new alert is as simple as calling NewAlertCmd()
	// with a key and a message. The formatting and stylings will be
	// handled by the AlertDefition types. Below are the included
	// alert types, but you can also create your own custom ones!
	// Check out AlertModel.RegisterNewAlertType()
	var alertCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			alertCmd = m.alert.NewAlertCmd(bubbleup.InfoKey, "New info alert.")
		case "w":
			alertCmd = m.alert.NewAlertCmd(bubbleup.WarnKey, "New warning alert...")
		case "e":
			alertCmd = m.alert.NewAlertCmd(bubbleup.ErrorKey, "New error alert!")
		case "d":
			alertCmd = m.alert.NewAlertCmd(bubbleup.DebugKey, "New debug alert?")
		case "q":
			return m, tea.Quit
		}

	}

	// Be sure to pass any received messages to the alert
	// model, and appropriately use the return values.
	// Reassign your stored alert with the updated alert,
	// and return the given command, either alone or via tea.Batch().
	outAlert, outCmd := m.alert.Update(msg)
	m.alert = outAlert.(bubbleup.AlertModel)

	return m, tea.Batch(outCmd, alertCmd)
}

func (m testModel) View() string {
	// Do any View stuff you need to like normal, and
	// call your alert's Render function to render any active
	// alerts over your content. Note: The alert model's View()
	// function is empty and is not meant to be called.
	return m.alert.Render(m.content)
}

func getTestContent() string {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	outStyle := lipgloss.NewStyle().Width(width-2).Height(height-2).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#00FFFF")).
		Align(lipgloss.Center, lipgloss.Center)

	return outStyle.Render("This is a test string, wow look at it go!")
}
