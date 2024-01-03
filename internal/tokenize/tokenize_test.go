package tokenize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		// string to tokenize
		s string
		// length of tokens
		length int
	}{
		{
			"foo",
			1,
		},
		{
			"a b c",
			3,
		},
		{
			"SQ *FOO BAR",
			3,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			tk := NewTokenize()
			tk.Parse(test.s)
			require.Equal(t, test.length, len(tk.Tokens()))
		})
	}
}

func TestPrevious(t *testing.T) {
	tests := []struct {
		// string to tokenize
		s string
		// length of tokens
		length int
	}{
		{
			"foo",
			0,
		},
		{
			"a b c",
			2,
		},
		{
			"SQ *FOO BAR",
			2,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			tk := NewTokenize()
			tk.Parse(test.s)
			count := 0
			for _, token := range tk.Tokens() {
				if token.Previous() == nil {
					continue
				}
				count++
			}
			require.Equal(t, test.length, count, fmt.Sprintf("string is '%s'", test.s))
		})
	}
}
