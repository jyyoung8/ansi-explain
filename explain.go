package main

import "fmt"

// Explanation is a short, human-readable description of one Token.
type Explanation struct {
	Summary string
	Detail  string
}

var sgrNames = map[int]string{
	0: "reset", 1: "bold", 2: "faint", 3: "italic", 4: "underline",
	7: "reverse video", 22: "normal intensity", 24: "underline off",
	30: "foreground black", 31: "foreground red", 32: "foreground green",
	33: "foreground yellow", 34: "foreground blue", 35: "foreground magenta",
	36: "foreground cyan", 37: "foreground white", 39: "default foreground",
	40: "background black", 41: "background red", 42: "background green",
	43: "background yellow", 44: "background blue", 45: "background magenta",
	46: "background cyan", 47: "background white", 49: "default background",
}

var csiFinalNames = map[byte]string{
	'A': "cursor up", 'B': "cursor down", 'C': "cursor forward", 'D': "cursor back",
	'H': "cursor position", 'J': "erase in display", 'K': "erase in line",
	'm': "select graphic rendition (SGR)",
}

var simpleFinalNames = map[byte]string{
	'7': "save cursor position", '8': "restore cursor position",
	'c': "reset terminal",
	'M': "reverse index (move up, scroll if needed)",
}

// Explain describes a single Token in plain language. It reads no global or
// terminal state, so the same Token always produces the same Explanation.
func Explain(t Token) Explanation {
	switch t.Kind {
	case TokenText:
		return Explanation{Summary: fmt.Sprintf("%d bytes of plain text", len(t.Text))}
	case TokenCSI:
		return explainCSI(t)
	case TokenOSC:
		return Explanation{Summary: "OSC (operating system command)", Detail: fmt.Sprintf("%q", t.Text)}
	case TokenSimple:
		if name, ok := simpleFinalNames[t.Final]; ok {
			return Explanation{Summary: name}
		}
		return Explanation{Summary: fmt.Sprintf("ESC %c (unrecognized)", t.Final)}
	default:
		return Explanation{Summary: "unrecognized escape sequence", Detail: fmt.Sprintf("%q", t.Text)}
	}
}

func explainCSI(t Token) Explanation {
	name, ok := csiFinalNames[t.Final]
	if !ok {
		return Explanation{Summary: fmt.Sprintf("CSI ... %c (unrecognized)", t.Final), Detail: fmt.Sprintf("%q", t.Text)}
	}
	if t.Final != 'm' {
		return Explanation{Summary: name, Detail: fmt.Sprintf("params=%v", t.Params)}
	}
	if len(t.Params) == 0 {
		return Explanation{Summary: name, Detail: sgrNames[0]}
	}
	parts := make([]string, 0, len(t.Params))
	for _, p := range t.Params {
		if n, ok := sgrNames[p]; ok {
			parts = append(parts, n)
		} else {
			parts = append(parts, fmt.Sprintf("code %d", p))
		}
	}
	return Explanation{Summary: name, Detail: joinComma(parts)}
}

func joinComma(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ", "
		}
		out += p
	}
	return out
}
