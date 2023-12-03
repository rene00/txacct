package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"transactionsearch/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TransactionPostcode struct{}

func NewTransactionPostcode() TransactionHandler {
	return TransactionPostcode{}
}

func (tp TransactionPostcode) Handle(ctx context.Context, db *sql.DB, transaction *Transaction) error {
	for _, token := range transaction.tokenize.TokensReversed() {
		if strings.ToLower(token.ValueString()) == "aus" {
			continue
		}

		q := []qm.QueryMod{
			qm.Where("lower(locality) = lower(?)", token.ValueString()),
			qm.InnerJoin("state s on postcode.state_id = s.id"),
		}

		s, err := stripStateFromString(ctx, db, token.ValueString())
		if err != nil {
			return err
		}

		s = stripCountryFromString(s)

		if token.Previous() != nil {
			previousTokenString, err := stripStateFromString(ctx, db, token.Previous().ValueString())
			if err != nil {
				return err
			}
			previousTokenString = stripCountryFromString(previousTokenString)
			query := fmt.Sprintf("%s %s", previousTokenString, s)
			q = append(q, qm.Or("locality ilike ?", query+"%"))
		}

		postcodes, err := models.Postcodes(q...).All(ctx, db)
		if err != nil {
			return err
		}

		for _, stateName := range allStatesPreferenced() {
			for _, postcode := range postcodes {
				state, err := postcode.State().One(ctx, db)
				if err != nil {
					return err
				}
				if stateName == state.Name {
					transaction.postcode = postcode
					transaction.state = state
					token.SetLocality(true)
					return nil
				}
			}
		}
	}

	return nil
}
