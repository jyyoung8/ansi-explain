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
