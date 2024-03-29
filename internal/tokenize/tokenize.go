package tokenize

import (
	"slices"
	"strings"
)

type Tokenize struct {
	raw    string
	tokens []*Token
}

// Parse accepts a string which is parsed and tokenized into t.tokens.
func (t *Tokenize) Parse(s string) {
	t.raw = s
	previousToken := &Token{}
	for idx, field := range strings.Fields(s) {
		token := &Token{value: field, position: idx, previous: previousToken}
		if idx == 0 {
			token.previous = nil
		}
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
	r := []*Token{}
	for _, token := range t.Tokens() {
		r = append(r, token)
	}
	slices.Reverse(r)
	return r
}

// Last returns to the last token.
func (t Tokenize) Last() *Token {
	token := t.tokens[len(t.tokens)-1]
	return token
}

// NewTokenize returns a new Tokenize.
func NewTokenize() Tokenize {
	return Tokenize{}
}

type Token struct {
	value    string
	position int
	// geo is true if token contains geo data which is postcode, state or country.
	geo bool
	// locality is true if the token contains locality data (i.e. postcode).
	locality bool
	// state is true if token contains state data (i.e. VIC).
	state bool
	// country is true if token contains country data (i.e. AUS).
	country  bool
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
	t.geo = true
	return
}

func (t *Token) SetState(b bool) {
	t.state = b
	t.geo = true
	return
}

func (t Token) IsState() bool {
	return t.state
}

func (t *Token) SetCountry(b bool) {
	t.country = b
	t.geo = true
	return
}

func (t Token) IsCountry() bool {
	return t.country
}

func (t Token) IsGeo() bool {
	return t.geo
}
