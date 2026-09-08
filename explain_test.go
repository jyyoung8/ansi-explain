package main

import "testing"

func TestExplainSGR256Color(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'm', Params: []int{38, 5, 208}}
	exp := Explain(tok)
	want := "foreground color 208 (256-color)"
	if exp.Detail != want {
		t.Fatalf("detail = %q, want %q", exp.Detail, want)
	}
}

func TestExplainSGRTruecolor(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'm', Params: []int{48, 2, 10, 20, 30}}
	exp := Explain(tok)
	want := "background RGB(10,20,30)"
	if exp.Detail != want {
		t.Fatalf("detail = %q, want %q", exp.Detail, want)
	}
}

func TestExplainSGRExtendedColorThenMoreParams(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'm', Params: []int{38, 5, 208, 1}}
	exp := Explain(tok)
	want := "foreground color 208 (256-color), bold"
	if exp.Detail != want {
		t.Fatalf("detail = %q, want %q", exp.Detail, want)
	}
}

func TestExplainDCS(t *testing.T) {
	tok := Token{Kind: TokenDCS, Text: []byte("\x1bPq...\x1b\\")}
	exp := Explain(tok)
	if exp.Summary != "DCS (device control string)" {
		t.Fatalf("summary = %q, want DCS summary", exp.Summary)
	}
}

func TestExplainDECPrivateModeEnable(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'h', Params: []int{1049}, Private: true}
	exp := Explain(tok)
	if exp.Summary != "enable DEC private mode" {
		t.Fatalf("summary = %q, want %q", exp.Summary, "enable DEC private mode")
	}
	want := "alternate screen buffer with cursor save"
	if exp.Detail != want {
		t.Fatalf("detail = %q, want %q", exp.Detail, want)
	}
}

func TestExplainDECPrivateModeDisableMultipleParams(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'l', Params: []int{25, 2004}, Private: true}
	exp := Explain(tok)
	if exp.Summary != "disable DEC private mode" {
		t.Fatalf("summary = %q, want %q", exp.Summary, "disable DEC private mode")
	}
	want := "cursor visibility, bracketed paste mode"
	if exp.Detail != want {
		t.Fatalf("detail = %q, want %q", exp.Detail, want)
	}
}

func TestExplainDECPrivateModeUnknown(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'h', Params: []int{9999}, Private: true}
	exp := Explain(tok)
	if exp.Detail != "mode 9999" {
		t.Fatalf("detail = %q, want %q", exp.Detail, "mode 9999")
	}
}

func TestExplainStandardModeIsNotTreatedAsPrivate(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'h', Params: []int{4}}
	exp := Explain(tok)
	if exp.Summary != "set mode" {
		t.Fatalf("summary = %q, want %q", exp.Summary, "set mode")
	}
	if exp.Detail != "params=[4]" {
		t.Fatalf("detail = %q, want %q", exp.Detail, "params=[4]")
	}
}

// TestExplainVimAltScreenSequence exercises the block of sequences vim
// sends on startup and on exit to take over and release the terminal: enter
// the alternate screen, hide the cursor, then the reverse on the way out.
// Real captures (from `tmux -CC` or `script`) follow this exact shape.
func TestExplainVimAltScreenSequence(t *testing.T) {
	enter := []byte("\x1b[?1049h\x1b[?25l")
	tokens := Tokenize(enter)
	if len(tokens) != 2 {
		t.Fatalf("got %d tokens, want 2: %+v", len(tokens), tokens)
	}
	if got := Explain(tokens[0]).Detail; got != "alternate screen buffer with cursor save" {
		t.Fatalf("first detail = %q", got)
	}
	if got := Explain(tokens[1]).Detail; got != "cursor visibility" {
		t.Fatalf("second detail = %q", got)
	}

	leave := []byte("\x1b[?25h\x1b[?1049l")
	tokens = Tokenize(leave)
	if len(tokens) != 2 {
		t.Fatalf("got %d tokens, want 2: %+v", len(tokens), tokens)
	}
	if got := Explain(tokens[0]).Summary; got != "enable DEC private mode" {
		t.Fatalf("first summary = %q", got)
	}
	if got := Explain(tokens[1]).Summary; got != "disable DEC private mode" {
		t.Fatalf("second summary = %q", got)
	}
}

func TestExplainSGRExtendedColorMissingArgs(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'm', Params: []int{38, 5}}
	exp := Explain(tok)
	want := "code 38, code 5"
	if exp.Detail != want {
		t.Fatalf("detail = %q, want %q", exp.Detail, want)
	}
}
