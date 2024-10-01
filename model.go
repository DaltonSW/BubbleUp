package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	InfoLevel = iota
	WarningLevel
	ErrorLevel

	InfoSymbol    = ""
	WarningSymbol = "󱈸"
	ErrorSymbol   = "󰬅"

	InfoColor    = lipgloss.Color("#00FF00")
	WarningColor = lipgloss.Color("#FFFF00")
	ErrorColor   = lipgloss.Color("#FF0000")
)

var (
	Symbols = map[NotifLevel]string{
		InfoLevel:    InfoSymbol,
		WarningLevel: WarningSymbol,
		ErrorLevel:   ErrorSymbol,
	}

	Colors = map[NotifLevel]lipgloss.Color{
		InfoLevel:    InfoColor,
		WarningLevel: WarningColor,
		ErrorLevel:   ErrorColor,
	}
)

type NotifLevel int

type NotifMsg struct {
	msg   string
	level NotifLevel
	dur   time.Duration
}

func newNotif(msg string, lvl NotifLevel, dur time.Duration) *notif {
	notifColor := Colors[lvl]
	notifSymbol := Symbols[lvl]

	notifStyle := lipgloss.NewStyle().Foreground(notifColor).Padding(1).
		Border(lipgloss.RoundedBorder()).BorderForeground(notifColor)

	return &notif{
		message:   msg,
		level:     lvl,
		deathTime: time.Now().Add(dur),
		symbol:    notifSymbol,
		style:     notifStyle,
	}

}

type notif struct {
	message   string
	level     NotifLevel
	deathTime time.Time
	symbol    string
	style     lipgloss.Style
}

func (n *notif) render(width int) string {
	return n.style.MaxWidth(width).Render(fmt.Sprintf("%v %v", n.symbol, n.message))
}

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

func (m Model) Render(content string) string {
	if !m.active {
		return content
	}

	notifString := m.activeNotif.render(m.width)

	// log.Println(notifString)

	notifSplit := strings.Split(notifString, "\n")
	contentSplit := strings.Split(content, "\n")

	// log.Println(notifSplit)
	// log.Println(contentSplit)

	outSplit := []string{}

	for i := range len(notifSplit) {
		notifLine := notifSplit[i]
		contentLine := contentSplit[i]
		outLine := ""

		// log.Println(notifLen)

		for _, ch := range notifLine {
			outLine += string(ch)
		}

		outLine = outLine + contentLine[len(notifLine):]

		outSplit = append(outSplit, outLine)
	}

	outSplit = append(outSplit, contentSplit[len(notifSplit):]...)

	return strings.Join(outSplit, "\n")
}

func (m Model) Notify(msg string, level NotifLevel, dur time.Duration) {

	if m.activeNotif != nil {
		return
	}

	m.activeNotif = newNotif(msg, level, dur)
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
