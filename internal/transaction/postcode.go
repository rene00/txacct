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
			query := fmt.Sprintf("%s %s", previousTokenString, s)
			q = append(q, qm.Or("locality ilike ?", query+"%"))
		}

		postcodes, err := models.Postcodes(q...).All(ctx, store.DB)
		if err != nil {
			return err
		}

		for _, stateName := range allStatesPreferenced() {
			for _, postcode := range postcodes {
				state, err := postcode.State().One(ctx, store.DB)
				if err != nil {
					return err
				}
				if stateName == state.Name {
					transaction.postcode = postcode
					transaction.state = state
					token.SetLocality(true)
					for _, s := range strings.Split(postcode.Locality, " ") {
						if strings.ToLower(token.Previous().ValueString()) == strings.ToLower(s) {
							token.Previous().SetLocality(true)
						}
					}
					return nil
				}
			}
		}
	}

	return nil
}
