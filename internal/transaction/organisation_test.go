package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildLikeQueryContents(t *testing.T) {
	tests := []struct {
		input  string
		expect []string
	}{
		{
			"foo",
			[]string{"foo"},
		},
		{
			"foo bar",
			[]string{"foo", "foo bar"},
		},
		{
			"SQ *foo",
			[]string{"*foo"},
		},
		{
			"SQ *pants bar",
			[]string{"*pants", "*pants bar"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			tx := NewTransaction(test.input, nil)
			to := TransactionOrganisation{}
			require.Equal(t, test.expect, to.buildLikeQueryContents(*tx))
		})
	}
}
