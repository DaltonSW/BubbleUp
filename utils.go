package bubbleup

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
)

// Obtained from https://github.com/charmbracelet/lipgloss/blob/master/get.go
func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := ansi.PrintableRuneWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}

// Obtained from https://github.com/charmbracelet/lipgloss/pull/102/commits/a075bfc9317152e674d661a2cdfe58144306e77a
// cutLeft cuts printable characters from the left.
func cutLeft(s string, cutWidth int) string {
	var (
		pos    int
		isAnsi bool
		ab     bytes.Buffer
		b      bytes.Buffer
	)
	for _, c := range s {
		var w int
		if c == ansi.Marker || isAnsi {
			isAnsi = true
			ab.WriteRune(c)
			if ansi.IsTerminator(c) {
				isAnsi = false
				if bytes.HasSuffix(ab.Bytes(), []byte("[0m")) {
					ab.Reset()
				}
			}
		} else {
			w = runewidth.RuneWidth(c)
		}

		if pos >= cutWidth {
			if b.Len() == 0 {
				if ab.Len() > 0 {
					b.Write(ab.Bytes())
				}
				if pos-cutWidth > 1 {
					b.WriteByte(' ')
					continue
				}
			}
			b.WriteRune(c)
		}
		pos += w
	}
	return b.String()
}

func lerpColor(start, end lipgloss.TerminalColor, step float64) lipgloss.Color {
	startR, startG, startB, _ := start.RGBA()
	endR, endG, endB, _ := end.RGBA()

	r := lerp(startR, endR, step)
	g := lerp(startG, endG, step)
	b := lerp(startB, endB, step)

	hex := fmt.Sprintf("#%02X%02X%02X", r, g, b)
	return lipgloss.Color(hex)

}

func lerp(a, b uint32, step float64) int {
	return int(float64(a) + step*(float64(b)-float64(a)))
}
