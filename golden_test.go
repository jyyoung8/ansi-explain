package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "regenerate golden files in testdata/ from current output")

// vimSessionCapture is shaped like what `script` or `tmux -CC` records from a
// real vim session: enter the alternate screen, hide the cursor, redraw with
// a syntax-highlighted line, set the window title over OSC, then unwind all
// of it on quit. The plain text in between (source lines, the status line)
// is left as-is since run() only reports the escape sequences.
var vimSessionCapture = []byte(
	"\x1b[?1049h" + // enter alternate screen, save cursor
		"\x1b[?25l" + // hide cursor while drawing
		"\x1b[H\x1b[2J" + // home cursor, clear screen
		"\x1b]0;vim README.md\x07" + // set window title
		"\x1b[1;1H" + // position cursor for the first line
		"\x1b[38;5;208m\x1b[1m" + // syntax highlight color, bold
		"func main() {" +
		"\x1b[0m" + // reset attributes
		"\x1b[?25h\x1b[?25l" + // cursor blinks back on then off for the redraw below
		"\x1b[24;1H" + // move to the status line
		"\x1b[0m" +
		"\x1b[?25h" + // cursor back on
		"\x1b[?1049l", // leave alternate screen, restore cursor
)

func TestGoldenVimSession(t *testing.T) {
	var out bytes.Buffer
	if err := run(nil, bytes.NewReader(vimSessionCapture), &out); err != nil {
		t.Fatalf("run: %v", err)
	}

	goldenPath := filepath.Join("testdata", "vim_session.golden")
	if *update {
		if err := os.WriteFile(goldenPath, out.Bytes(), 0o644); err != nil {
			t.Fatalf("writing golden file: %v", err)
		}
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden file: %v", err)
	}
	if out.String() != string(want) {
		t.Fatalf("output does not match %s\n--- got ---\n%s\n--- want ---\n%s", goldenPath, out.String(), want)
	}
}
