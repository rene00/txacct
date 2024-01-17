package transaction

import (
	"context"
	"fmt"
	"strings"
	"transactionsearch/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TransactionPostcode struct{}

func NewTransactionPostcode() TransactionHandler {
	return TransactionPostcode{}
}

func (tp TransactionPostcode) Handle(ctx context.Context, store Store, transaction *Transaction) error {
	var postcodes []*models.Postcode
	for _, token := range transaction.tokenize.TokensReversed() {
		if strings.ToLower(token.ValueString()) == "aus" {
			continue
		}

		q := []qm.QueryMod{
			qm.Where("lower(locality) = lower(?)", token.ValueString()),
			qm.InnerJoin("state s on postcode.state_id = s.id"),
		}

		s, err := stripStateFromString(ctx, store.DB, token.ValueString())
		if err != nil {
			return err
		}

		s = stripCountryFromString(s)

		if token.Previous() != nil {
			previousTokenString, err := stripStateFromString(ctx, store.DB, token.Previous().ValueString())
			if err != nil {
				return err
			}
			previousTokenString = stripCountryFromString(previousTokenString)
			query := fmt.Sprintf("%s", previousTokenString)
			if len(s) != 0 {
				query = fmt.Sprintf("%s %s", previousTokenString, s)
			}
			q = append(q, qm.Or("locality ilike ?", query+"%"))
		}

		postcodeSlice, err := models.Postcodes(q...).All(ctx, store.DB)
		if err != nil {
			return err
		}

		for _, postcode := range postcodeSlice {
			postcodes = append(postcodes, postcode)
		}
		transaction.postcodes = postcodes

		if len(postcodes) == 1 {
			postcode := postcodes[0]
			state, err := postcode.State().One(ctx, store.DB)
			if err != nil {
				return err
			}
			transaction.postcode = postcode
			transaction.state = state
		}
	}

	return nil
}
