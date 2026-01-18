package bubbleup

import (
	"bytes"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/wordwrap"
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

// hangingWrap wraps text with a prefix to provide hanging indents
func hangingWrap(prefix, msg string, textWidth int) string {
	prefix = prefix + " "
	indentW := lipgloss.Width(prefix)
	avail := textWidth - indentW
	if avail < 1 {
		// Degenerate case: not enough room for message; just show prefix and message raw.
		return prefix + msg
	}

	// Wrap message to the available width.
	// wordwrap.WrapString wraps on spaces; it will still break long tokens if needed.
	wrapped := wordwrap.String(msg, avail)

	// Add hanging indent to subsequent lines.
	indent := strings.Repeat(" ", indentW)
	lines := strings.Split(wrapped, "\n")
	for i := 1; i < len(lines); i++ {
		lines[i] = indent + lines[i]
	}

	return prefix + strings.Join(lines, "\n")
}
