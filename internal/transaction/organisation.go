package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"transactionsearch/internal/tokenize"
	"transactionsearch/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TransactionOrganisation struct{}

func NewTransactionOrganisation() TransactionHandler {
	return TransactionOrganisation{}
}

func (to TransactionOrganisation) Handle(ctx context.Context, db *sql.DB, transaction *Transaction) error {
	var q []qm.QueryMod

	if transaction.postcode != nil {
		q = []qm.QueryMod{
			qm.InnerJoin("business_code bc on organisation.business_code_id = bc.id"),
			qm.InnerJoin(fmt.Sprintf("postcode p on organisation.postcode_id = %d", transaction.postcode.ID)),
		}

		var likeQueryContents []string

		likeQueryContents = to.buildLikeQueryContents(*transaction)
		for i := len(likeQueryContents) - 1; i >= 0; i-- {
			v := likeQueryContents[i]
			q = append(q, qm.Or("name ILIKE ?", v+"%"))
		}

		organisations, err := models.Organisations(q...).All(ctx, db)
		if err != nil {
			return err
		}

		for _, organisation := range organisations {
			for i := len(likeQueryContents) - 1; i >= 0; i-- {
				q := likeQueryContents[i]
				if strings.ToLower(q) == strings.ToLower(organisation.Name) {
					transaction.organisation = organisation
					return nil
				}
			}
		}
	}
	return nil
}

// querySkipToken accepts a token and attempts to determine if this token
// should be skipped when building the query used to find the organisation.
func (to TransactionOrganisation) querySkipToken(token tokenize.Token) bool {
	if token.IsLocality() {
		return true
	}

	// Skip tokens based on their position and value. In the future, tokens
	// could be classified as payment processors such as Stripe and Square.
	skip := map[int][]string{
		0: []string{
			"SP",
			"SQ",
		},
	}

	skipSlice, exists := skip[token.Position()]
	if exists {
		for _, skipString := range skipSlice {
			if strings.ToLower(skipString) == strings.ToLower(token.ValueString()) {
				return true
			}
		}
	}

	return false
}

// buildLikeQueryContents accepts a transaction and returns a slice of strings which
// are the ILIKE queries that will be used when searching for an organisation.
func (to TransactionOrganisation) buildLikeQueryContents(transaction Transaction) []string {
	var s []string
	for _, token := range transaction.tokenize.Tokens() {
		if to.querySkipToken(*token) {
			continue
		}
		v := token.ValueString()
		previousToken := token.Previous()
		if previousToken != nil {
			v = fmt.Sprintf("%s %s", token.Previous().ValueString(), token.ValueString())
		}
		s = append(s, strings.TrimPrefix(v, " "))
	}
	return s
}
