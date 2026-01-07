package main

// Only two imports that are required
import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.dalton.dog/bubbleup"
	"golang.org/x/term"
)

// testModel is this example app's struct that is neeeed to be
// passed to tea.NewProgram() in order to run a BubbleTea app.
type testModel struct {
	content    string              // Our instruction content to display
	fontChoice string              // Current icon font: Unicode, NerdFont, or ASCII
	alert      bubbleup.AlertModel // Model that implements our BubbleUp alert
	ExitKey    tea.KeyMsg          // Track program exit keys so we can switch fonts
	KeyPressed bool
	PressedKey tea.KeyMsg
}

// main is an example app that let's you see to use Bubble up.
// The part you need to understand is really two lines of code long
// Look for BEGIN and END comments.
// Everything else in this file is to provide an example app.
func main() {
	// Start by getting the screen size.
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	// Create a model for the program, start with Unicode fonts
	m := testModel{fontChoice: "Unicode"}

	// Loop until 'q' is pressed. User can press 'N', 'U', or 'A' to switch fonts
	for {
		// Do whatever you need to do to get your
		// model looking pretty and initialized.
		m.content = getTestContent(m, width, height)

		// We need either UseNextFont==true or useNerdFont==false because
		// NewAlertModel()'s 2nd parameter expects to indicate if it should
		// use NerdFont (true) or ASCII (false). If we want Unicode instead
		// we'll override with a WithUnicodePrefix() after NewAlertModel().
		useNerdFont := m.fontChoice == "NerdFont"

		// BEGIN here to understand how to create and use Bubble Up

		// Create a new alert model and embed it within your program's model
		//  width = 50 (max), minWidth = 10, duration = 10
		// useNerdFont = default is false, but example app user can change.
		m.alert = bubbleup.NewAlertModel(50, useNerdFont, 10).
			WithMinWidth(15).     // Dynamic width: alerts will size 15-50 chars
			WithAllowEscToClose() // Allow <esc> to close an alert before timeout
		// based on message length

		if m.fontChoice == "Unicode" {
			// We don't pass in Unicode as an option for the 2nd parameter
			// because this option was added after the original development
			// and we wanted to maintain backward compatibility where the
			// original choice was either NerdFont or not. So here we override
			// the choice of NerdFont above and change to Unicode options
			// because NerdFont only works if the user has it installed and
			// Unicode should work pretty much everywhere.
			m.alert = m.alert.WithUnicodePrefix()
		}

		// Also see the required BubbleTea methods Init(), Update() and
		// View() to understand how to use Bubble up.

		// END of what you need to know from this file to use Bubble Up

		// The rest of this is just a convenient way to allow for running
		// the program and switching the font. This code is based on an
		// earier example and rather thar write a well-structured TEA app
		// that you might feel like you need to understand, we just hacked
		// it for you to see the different icon fonts used by Bubble Up.
		// You do NOT need to understand the following to learn how to use
		// BubbleUp.
		p := tea.NewProgram(m, tea.WithAltScreen())
		result, err := p.Run()
		switch {
		case err != nil:
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		default:
			m = result.(testModel)
			switch m.ExitKey.String() {
			case "N", "n":
				m.fontChoice = "NerdFont"
			case "A", "a":
				m.fontChoice = "ASCII"
			case "U", "u":
				m.fontChoice = "Unicode"
			case "esc", "q":
				fallthrough
			default:
				return
			}
		}
	}
}

// Init is the required BubbleTea method required to initialize a BubbleTea app
func (m testModel) Init() tea.Cmd {
	// Be sure to return the result of the alert models' Init()
	// If you need to also return one or more commands,
	// be sure to use tea.Batch() to bundle them together.
	return m.alert.Init()
}

// Update is the required BubbleTea method used to process key and other events
func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Creating a new alert is as simple as calling NewAlertCmd()
	// with a key and a message. The formatting and stylings will be
	// handled by the AlertDefition types. Below are the included
	// alert types, but you can also create your own custom ones!
	// Check out AlertModel.RegisterNewAlertType()
	//
	// This example demonstrates the WithPosition() method to change where
	// notifications appear on screen. There are 6 positions available:
	// TopLeftPosition, TopCenterPosition, TopRightPosition, BottomLeftPosition,
	// BottomCenterPosition, BottomRightPosition.
	//
	// This example also allows exiting on a current icon font selection so the alert
	// can be recreated to use the newly-selected icon font from a list of choices
	// that include Unicode, NerdFont and ASCII.
	var alertCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.PressedKey = msg
		m.KeyPressed = true
		switch msg.String() {
		case "i":
			m.alert = m.alert.WithPosition(bubbleup.TopLeftPosition)
			alertCmd = m.alert.NewAlertCmd(bubbleup.InfoKey, "Short message")
		case "w":
			m.alert = m.alert.WithPosition(bubbleup.TopCenterPosition)
			alertCmd = m.alert.NewAlertCmd(bubbleup.WarnKey, "Medium length message")
		case "e":
			m.alert = m.alert.WithPosition(bubbleup.TopRightPosition)
			alertCmd = m.alert.NewAlertCmd(bubbleup.ErrorKey, "This is an error message that is longer to show dynamic width and wrapping")
		case "d":
			m.alert = m.alert.WithPosition(bubbleup.BottomLeftPosition)
			alertCmd = m.alert.NewAlertCmd(bubbleup.DebugKey, "Shortest")
		case "I":
			m.alert = m.alert.WithPosition(bubbleup.BottomCenterPosition)
			alertCmd = m.alert.NewAlertCmd(bubbleup.InfoKey, "Medium message here")
		case "W":
			m.alert = m.alert.WithPosition(bubbleup.BottomRightPosition)
			alertCmd = m.alert.NewAlertCmd(bubbleup.WarnKey, "Another long warning to demonstrate width variation when the text is super long so it will wrap on three (3) lines")
		case "q", "U", "u", "N", "n", "A", "a":
			m.ExitKey = msg
			return m, tea.Quit
		case "esc":
			if !m.alert.HasActiveAlert() {
				m.ExitKey = msg
				return m, tea.Quit
			}
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

// View is the required BubbleTea method used to render output to the screen
func (m testModel) View() string {
	// Do any View stuff you need to like normal, and
	// call your alert's Render function to render any active
	// alerts over your content. Note: The alert model's View()
	// function is empty and is not meant to be called.
	content := m.alert.Render(m.content)
	if m.KeyPressed {
		legend := fmt.Sprintf("Keypress: <%s>", m.PressedKey.String())
		quit := "q - Quit"
		legendStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
		content = strings.Replace(content,
			quit+strings.Repeat(" ", 40),
			quit+strings.Repeat(" ", 10)+"["+legendStyle.Render(legend)+"]"+strings.Repeat(" ", 30-len(legend)-2),
			-1,
		)
	}
	return content
}

// getFontMenuChoices sets up icon font switcher menu.
// You do NOT need to understand it to learn how to use BubbleUp.
func getFontMenuChoices(model testModel) string {
	switch model.fontChoice[0] {
	case 'N', 'n': // NerdFont
		return "U)nicode or A)SCII"
	case 'A', 'a': // ASCII
		return "N)erdFont or U)nicode"
	case 'U', 'u': // Unicode
		fallthrough
	default:
		return "N)erdFont or A)SCII"
	}
}

// getTestContent sets up the content for the example app.
// You do NOT need to understand it to learn how to use BubbleUp.
func getTestContent(model testModel, width, height int) string {

	// Initialize so color styles
	style := lipgloss.NewStyle()
	greenStyle := style.Foreground(lipgloss.Color("#00FF00"))
	yellowStyle := style.Foreground(lipgloss.Color("#FFFF00"))

	// Compose current "Icon Font" indicator and mini-font selector menu
	fontMenu := fmt.Sprintf(`Icon Font: %s; To change: %s`,
		yellowStyle.Render(model.fontChoice),
		getFontMenuChoices(model),
	)

	// Compose the instructions menu, add the mini-font selector menu
	instructions := fmt.Sprintf(`Press keys to test different positions:

  i - Info (Top-Left) - Short Message
  w - Warning (Top-Center) - medium length message
  e - Error (Top-Right) - Longer error message 
  d - Debug (Bottom-Left) - Shortest
  I - Info (Bottom-Center) - Medium message
  W - Warning (Bottom-Right) - A really long 3 line message 

  q - Quit

Note: Alerts are dynamic in width — 15-50 chars — and their
      width is derived from their message's render length.
      %s`,
		fontMenu,
	)

	// In rare cases an width won't be specified, line when running in terminal
	// emulation mode like within a JetBrains IDE, so provide a reasonable default.
	if width == 0 {
		width = 8 + linesWidth(instructions)
	}

	// Create and style the welcome banner
	banner := "Welcome to the BubbleUp Demo!"
	bannerStyle := greenStyle.
		PaddingTop(1).
		PaddingLeft(8).
		PaddingBottom(1).
		Width(width-15).
		Align(lipgloss.Center, lipgloss.Center)

	// Inner content box (no border)
	instructionsStyle := lipgloss.NewStyle().
		PaddingTop(1).
		PaddingLeft(8).
		PaddingBottom(2).
		Align(lipgloss.Left, lipgloss.Top)

	// Combine the banner with the instructions
	innerContent :=
		bannerStyle.Render(banner) +
			instructionsStyle.Render(instructions)

	// Outer box with border
	outerStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#00FFFF"))

	return outerStyle.Render(innerContent)
}

// linesWidth is just a simple helper app
// You do NOT need to understand it to learn how to use BubbleUp.
func linesWidth(lines string) (width int) {
	for _, line := range strings.Split(lines, "\n") {
		width = max(width, len(line))
	}
	return width
}
