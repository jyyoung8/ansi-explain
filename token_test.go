package main

import (
	"reflect"
	"testing"
)

func TestTokenizePlainText(t *testing.T) {
	tokens := Tokenize([]byte("hello"))
	if len(tokens) != 1 || tokens[0].Kind != TokenText {
		t.Fatalf("got %+v, want single text token", tokens)
	}
}

func TestTokenizeSGR(t *testing.T) {
	input := []byte("\x1b[31mred\x1b[0m")
	tokens := Tokenize(input)
	if len(tokens) != 3 {
		t.Fatalf("got %d tokens, want 3: %+v", len(tokens), tokens)
	}
	if tokens[0].Kind != TokenCSI || tokens[0].Final != 'm' {
		t.Fatalf("first token = %+v, want CSI 'm'", tokens[0])
	}
	if !reflect.DeepEqual(tokens[0].Params, []int{31}) {
		t.Fatalf("params = %v, want [31]", tokens[0].Params)
	}
	if tokens[1].Kind != TokenText {
		t.Fatalf("second token = %+v, want text", tokens[1])
	}
}

func TestTokenizeDCS(t *testing.T) {
	input := []byte("\x1bPq#0;2;0;0;0#1~~@@vv@@~~@@~~$#1~~@@\x1b\\")
	tokens := Tokenize(input)
	if len(tokens) != 1 || tokens[0].Kind != TokenDCS {
		t.Fatalf("got %+v, want single DCS token", tokens)
	}
	if string(tokens[0].Text) != string(input) {
		t.Fatalf("text = %q, want %q", tokens[0].Text, input)
	}
}

func TestTokenizeDCSDoesNotStopAtBEL(t *testing.T) {
	// DCS payloads may contain a raw BEL byte; only ESC \ ends the sequence.
	input := []byte("\x1bPabc\x07def\x1b\\")
	tokens := Tokenize(input)
	if len(tokens) != 1 || tokens[0].Kind != TokenDCS {
		t.Fatalf("got %+v, want single DCS token", tokens)
	}
	if len(tokens[0].Text) != len(input) {
		t.Fatalf("text = %q, want it to span the whole input", tokens[0].Text)
	}
}

func TestTokenizeDCSUnterminated(t *testing.T) {
	input := []byte("\x1bPq#0;2;0;0")
	tokens := Tokenize(input)
	if len(tokens) != 1 || tokens[0].Kind != TokenUnknown {
		t.Fatalf("got %+v, want single unknown token", tokens)
	}
}

func TestTokenizeDECPrivateMode(t *testing.T) {
	input := []byte("\x1b[?1049h")
	tokens := Tokenize(input)
	if len(tokens) != 1 || tokens[0].Kind != TokenCSI {
		t.Fatalf("got %+v, want single CSI token", tokens)
	}
	tok := tokens[0]
	if !tok.Private {
		t.Fatalf("token = %+v, want Private set", tok)
	}
	if !reflect.DeepEqual(tok.Params, []int{1049}) {
		t.Fatalf("params = %v, want [1049]", tok.Params)
	}
	if tok.Final != 'h' {
		t.Fatalf("final = %q, want 'h'", tok.Final)
	}
}

func TestTokenizeCSIWithoutPrivateMarkerIsNotPrivate(t *testing.T) {
	tokens := Tokenize([]byte("\x1b[4h"))
	if len(tokens) != 1 || tokens[0].Private {
		t.Fatalf("got %+v, want non-private CSI token", tokens)
	}
}

func TestParseParamsEmptyField(t *testing.T) {
	got := parseParams([]byte("1;;3"))
	want := []int{1, 0, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestExplainSGRReset(t *testing.T) {
	tok := Token{Kind: TokenCSI, Final: 'm', Params: []int{0}}
	exp := Explain(tok)
	if exp.Detail != "reset" {
		t.Fatalf("detail = %q, want %q", exp.Detail, "reset")
	}
}

func TestExplainIsDeterministic(t *testing.T) {
	tok := Token{Kind: TokenSimple, Final: '7'}
	if Explain(tok) != Explain(tok) {
		t.Fatalf("Explain is not pure: same input gave different results")
	}
}
