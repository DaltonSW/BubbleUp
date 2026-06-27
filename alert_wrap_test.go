package bubbleup

import (
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/muesli/reflow/ansi"
)

// renderTestAlert renders an alert and returns the inner text lines, with the
// border and the single padding cell on each side stripped off.
func renderTestAlert(t *testing.T, prefix, msg string, width, minWidth int) []string {
	t.Helper()
	a := &alert{
		message:     msg,
		deathTime:   time.Now().Add(time.Minute),
		prefix:      prefix,
		foreColor:   errorColor,
		width:       width,
		minWidth:    minWidth,
		curLerpStep: 1.0,
		position:    TopLeftPosition,
	}

	lines := strings.Split(a.render(), "\n")

	// Every line of a rendered box must share the same total width.
	boxWidth := ansi.PrintableRuneWidth(lines[0])
	for i, ln := range lines {
		if w := ansi.PrintableRuneWidth(ln); w != boxWidth {
			t.Fatalf("ragged box: line %d width=%d, want %d", i, w, boxWidth)
		}
	}

	// Drop the top and bottom border rows, then strip the side border + padding
	// ("│ " ... " │") from each remaining row, leaving the wrapped text.
	inner := make([]string, 0, len(lines)-2)
	for _, ln := range lines[1 : len(lines)-1] {
		s := stripANSI(ln)
		s = strings.TrimPrefix(s, "│ ")
		s = strings.TrimSuffix(s, " │")
		inner = append(inner, strings.TrimRight(s, " "))
	}
	return inner
}

func stripANSI(s string) string {
	var b strings.Builder
	var isAnsi bool
	for _, c := range s {
		if c == ansi.Marker {
			isAnsi = true
		}
		if isAnsi {
			if ansi.IsTerminator(c) {
				isAnsi = false
			}
			continue
		}
		b.WriteRune(c)
	}
	return b.String()
}

// TestHangingIndentPreserved guards against the width-wrapping regression where
// lipgloss re-wraps overflow lines (because the text area was mis-measured) and
// the re-wrapped fragment loses its hanging indent. Wider ASCII prefixes such as
// "[!!]" and "(!)" are the ones that triggered it.
func TestHangingIndentPreserved(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		msg    string
	}{
		{"error-ascii", ErrorASCIIPrefix, "This is an error message that is longer to show dynamic width and wrapping"},
		{"warn-ascii", WarningASCIIPrefix, "Another long warning to demonstrate width variation when the text is super long so it will wrap on three (3) lines"},
		{"info-ascii", InfoASCIIPrefix, "This is an error message that is longer to show dynamic width and wrapping"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inner := renderTestAlert(t, c.prefix, c.msg, 50, 15)

			if len(inner) < 2 {
				t.Fatalf("expected the message to wrap onto multiple lines, got %d: %q", len(inner), inner)
			}

			// The message column is where text begins after "<prefix> ".
			indentW := lipgloss.Width(c.prefix + " ")
			indent := strings.Repeat(" ", indentW)

			if !strings.HasPrefix(inner[0], c.prefix+" ") {
				t.Errorf("first line should start with the prefix, got %q", inner[0])
			}
			for i := 1; i < len(inner); i++ {
				if !strings.HasPrefix(inner[i], indent) {
					t.Errorf("continuation line %d lost its hanging indent (want %d leading spaces): %q",
						i, indentW, inner[i])
				}
			}
		})
	}
}

// TestLongTokenKeepsIndent guards the case where a single word is wider than the
// available text column. It must be hard-wrapped to the column width while every
// resulting line keeps the hanging indent (no overflow, no orphaned fragments).
func TestLongTokenKeepsIndent(t *testing.T) {
	const longWord = "supercalifragilisticexpialidocioussupercalifragilistic"
	inner := renderTestAlert(t, ErrorASCIIPrefix, "aaaa "+longWord+" bbb", 30, 30)

	indentW := lipgloss.Width(ErrorASCIIPrefix + " ")
	indent := strings.Repeat(" ", indentW)
	for i := 1; i < len(inner); i++ {
		if !strings.HasPrefix(inner[i], indent) {
			t.Errorf("continuation line %d lost its hanging indent: %q", i, inner[i])
		}
	}
}

// TestDynamicWidthFitsShortMessage verifies a short message that fits within the
// max width is laid out on a single line, i.e. the box reserves exactly the
// border+padding frame and no more.
func TestDynamicWidthFitsShortMessage(t *testing.T) {
	inner := renderTestAlert(t, InfoASCIIPrefix, "Short message", 50, 15)
	if len(inner) != 1 {
		t.Fatalf("short message should fit on one line, got %d lines: %q", len(inner), inner)
	}
	want := InfoASCIIPrefix + " Short message"
	if inner[0] != want {
		t.Errorf("got %q, want %q", inner[0], want)
	}
}
