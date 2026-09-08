package main

// Token is a single unit found in a byte stream: either a run of plain
// bytes or a parsed escape sequence.
type Token struct {
	Kind    TokenKind
	Text    []byte // raw bytes that make up this token
	Final   byte   // final byte of a CSI or simple sequence; 0 if not applicable
	Params  []int  // numeric parameters parsed from a CSI sequence
	Private bool   // CSI carried a leading '?', marking a DEC private sequence
}

type TokenKind int

const (
	TokenText    TokenKind = iota // plain, non-escape bytes
	TokenCSI                      // ESC [ params... final
	TokenOSC                      // ESC ] ... BEL or ESC \
	TokenDCS                      // ESC P ... ST (ESC \)
	TokenSimple                   // ESC followed by exactly one byte
	TokenUnknown                  // ESC introduced something we could not parse
)

const (
	esc = 0x1b
	bel = 0x07
)

// Tokenize splits raw terminal output into a sequence of Tokens. It does
// not mutate input, and the Text of every returned token, concatenated in
// order, reproduces input exactly.
func Tokenize(input []byte) []Token {
	var tokens []Token
	i := 0
	for i < len(input) {
		if input[i] != esc {
			start := i
			for i < len(input) && input[i] != esc {
				i++
			}
			tokens = append(tokens, Token{Kind: TokenText, Text: input[start:i]})
			continue
		}
		tok, n := parseEscape(input[i:])
		tokens = append(tokens, tok)
		i += n
	}
	return tokens
}

// parseEscape parses a single escape sequence starting at input[0] == ESC.
// It returns the token and the number of bytes consumed, always at least 1.
func parseEscape(input []byte) (Token, int) {
	if len(input) < 2 {
		return Token{Kind: TokenUnknown, Text: input}, len(input)
	}
	switch input[1] {
	case '[':
		return parseCSI(input)
	case ']':
		return parseOSC(input)
	case 'P':
		return parseDCS(input)
	default:
		return Token{Kind: TokenSimple, Text: input[:2], Final: input[1]}, 2
	}
}

// parseCSI parses "ESC [ params intermediates final" per ECMA-48: params
// occupy 0x30-0x3F, intermediates 0x20-0x2F, and the final byte 0x40-0x7E. A
// leading '?' among the params marks a DEC private sequence (e.g. the
// cursor-visibility and alternate-screen toggles terminals send constantly)
// rather than a standard ECMA-48 one; it is stripped before the numeric
// params are parsed and recorded on the token instead.
func parseCSI(input []byte) (Token, int) {
	i := 2
	for i < len(input) && input[i] >= 0x30 && input[i] <= 0x3F {
		i++
	}
	paramBytes := input[2:i]
	private := len(paramBytes) > 0 && paramBytes[0] == '?'
	if private {
		paramBytes = paramBytes[1:]
	}
	for i < len(input) && input[i] >= 0x20 && input[i] <= 0x2F {
		i++
	}
	if i >= len(input) || input[i] < 0x40 || input[i] > 0x7E {
		return Token{Kind: TokenUnknown, Text: input[:i]}, i
	}
	final := input[i]
	i++
	return Token{
		Kind:    TokenCSI,
		Text:    input[:i],
		Final:   final,
		Params:  parseParams(paramBytes),
		Private: private,
	}, i
}

// parseOSC parses "ESC ] ... terminator", where the terminator is BEL or
// the two-byte String Terminator ESC \.
func parseOSC(input []byte) (Token, int) {
	i := 2
	for i < len(input) {
		if input[i] == bel {
			return Token{Kind: TokenOSC, Text: input[:i+1]}, i + 1
		}
		if input[i] == esc && i+1 < len(input) && input[i+1] == '\\' {
			return Token{Kind: TokenOSC, Text: input[:i+2]}, i + 2
		}
		i++
	}
	return Token{Kind: TokenUnknown, Text: input}, len(input)
}

// parseDCS parses "ESC P ... ST", where ST is the two-byte string
// terminator ESC \. Unlike OSC, DCS is not terminated by a bare BEL: DCS
// payloads (e.g. Sixel graphics, terminfo queries) can legitimately contain
// 0x07 as data.
func parseDCS(input []byte) (Token, int) {
	i := 2
	for i < len(input) {
		if input[i] == esc && i+1 < len(input) && input[i+1] == '\\' {
			return Token{Kind: TokenDCS, Text: input[:i+2]}, i + 2
		}
		i++
	}
	return Token{Kind: TokenUnknown, Text: input}, len(input)
}

// parseParams turns the semicolon-separated parameter bytes of a CSI
// sequence into ints. A missing field (as in the middle of "1;;3") becomes
// 0, matching the ECMA-48 default-value convention.
func parseParams(b []byte) []int {
	if len(b) == 0 {
		return nil
	}
	var params []int
	n := 0
	for _, c := range b {
		if c == ';' {
			params = append(params, n)
			n = 0
			continue
		}
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	params = append(params, n)
	return params
}
