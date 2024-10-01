package main

import (
	"fmt"
	"time"

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

	NotifWidth = 40
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

//package bubbleup

type NotifLevel int

type NotifMsg struct {
	msg   string
	dur   time.Duration
	level NotifLevel

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

func newNotif(msg string, lvl NotifLevel, dur time.Duration) *notif {
	notifColor := Colors[lvl]
	notifSymbol := Symbols[lvl]

	notifStyle := lipgloss.NewStyle().Foreground(notifColor).Width(NotifWidth).
		Border(lipgloss.RoundedBorder()).BorderForeground(notifColor)

	return &notif{
		message:   msg,
		level:     lvl,
		deathTime: time.Now().Add(dur),
		symbol:    notifSymbol,
		style:     notifStyle,
		width:     NotifWidth,
	}

}

type notif struct {
	message   string
	level     NotifLevel
	deathTime time.Time
	symbol    string
	style     lipgloss.Style
	width     int
}

func (n *notif) render() string {
	return n.style.Render(fmt.Sprintf("%v %v", n.symbol, n.message))
}
