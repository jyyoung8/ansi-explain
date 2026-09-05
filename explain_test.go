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

func TestExplainSGRExtendedColorMissingArgs(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'm', Params: []int{38, 5}}
	exp := Explain(tok)
	want := "code 38, code 5"
	if exp.Detail != want {
		t.Fatalf("detail = %q, want %q", exp.Detail, want)
	}
}
