package bubbleup

import (
	"bytes"
	"strings"

	"charm.land/lipgloss/v2"
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

// cutRight keeps printable characters from the left, up to keepWidth cells.
// ANSI escape sequences are preserved. Complement to cutLeft().
func cutRight(s string, keepWidth int) string {
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
				b.Write(ab.Bytes())
				ab.Reset()
			}
			continue
		}

		w = runewidth.RuneWidth(c)
		if pos+w > keepWidth {
			break
		}

		b.WriteRune(c)
		pos += w
	}

	// Reset to avoid color bleed
	if b.Len() > 0 && !bytes.HasSuffix(b.Bytes(), []byte("[0m")) {
		b.WriteByte(ansi.Marker)
		b.WriteString("[0m")
	}

	return b.String()
}

// hangingWrap lays out "<prefix> <message>" within textWidth cells, wrapping the
// message with a hanging indent so continuation lines align under the message
// rather than under the prefix.
//
// The prefix occupies a fixed-width column and the message is rendered as a
// second column wrapped to the remaining width; joining them horizontally yields
// the indent for free, since lipgloss pads the blank rows beneath the prefix.
// Letting lipgloss own the wrapping also means long, unbreakable tokens are hard
// wrapped to the column width instead of overflowing.
func hangingWrap(prefix, msg string, textWidth int) string {
	prefix = prefix + " "
	indentW := lipgloss.Width(prefix)
	avail := textWidth - indentW
	if avail < 1 {
		// Degenerate case: no room for the message beside the prefix. Hand the
		// whole thing back and let the caller's box wrap it.
		return prefix + msg
	}

	message := lipgloss.NewStyle().Width(avail).Render(msg)
	return lipgloss.JoinHorizontal(lipgloss.Top, prefix, message)
}
