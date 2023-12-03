package tokenize

import (
	"strings"
)

type Tokenize struct {
	raw    string
	tokens []*Token
}

func (t *Tokenize) Parse(s string) {
	t.raw = s
	previousToken := &Token{}
	for idx, field := range strings.Fields(s) {
		token := &Token{value: field, position: idx, previous: previousToken}
		t.tokens = append(t.tokens, token)
		previousToken = token
	}
}

// Tokens returns a slice of tokens.
func (t *Tokenize) Tokens() []*Token {
	return t.tokens
}

// Tokens returns a slice of tokens in reverse.
func (t Tokenize) TokensReversed() []*Token {
	tokens := make([]*Token, len(t.tokens))
	for i, j := len(t.tokens)-1, 0; i >= 0; i, j = i-1, j+1 {
		tokens[j] = t.tokens[i]
	}
	return tokens
}

// Last returns to the last token.
func (t Tokenize) Last() Token {
	token := t.tokens[len(t.tokens)-1]
	return *token
}

// NewTokenize returns a new Tokenize.
func NewTokenize() Tokenize {
	return Tokenize{}
}

type Token struct {
	value    string
	position int
	// locality is true if the token contains locality data (i.e. postcode).
	locality bool
	previous *Token
}

func (t Token) Position() int {
	return t.position
}

func (t Token) Previous() *Token {
	return t.previous
}

func (t Token) ValueString() string {
	return t.value
}

func (t Token) IsLocality() bool {
	return t.locality
}

func (t *Token) SetLocality(b bool) {
	t.locality = b
	return
}
