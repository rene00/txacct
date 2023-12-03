package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"transactionsearch/internal/tokenize"
	"transactionsearch/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TransactionState struct{}

func NewTransactionState() TransactionHandler {
	return TransactionState{}
}

func (ts TransactionState) Handle(ctx context.Context, db *sql.DB, transaction *Transaction) error {
	states, err := models.States().All(ctx, db)
	if err != nil {
		return err
	}

	tokens := transaction.tokenize.TokensReversed()
	for _, token := range tokens {
		for _, state := range states {
			shortState := state.Name[:2]
			re := regexp.MustCompile(fmt.Sprintf(`(?i)(?P<prefix>\w+)?(?P<state>%s|%s)$`, state.Name, shortState))
			match := re.FindStringSubmatch(token.ValueString())
			if len(match) == 3 {
				token.SetLocality(true)
				transaction.state = state
				return nil
			}
		}
	}

	// Lookup state using postcode. Get the last 2 tokens and lookup in
	// postcode. If postcode found, use the state from the postcode.
	combined := []*tokenize.Token{}
	for idx, token := range tokens {
		if idx == 2 {
			break
		}
		locality := token.ValueString()

		if strings.ToLower(locality) == "aus" {
			token.SetLocality(true)
			continue
		}

		re := regexp.MustCompile(`(?i)AUS$`)
		locality = re.ReplaceAllString(locality, "")

		combined = append(combined, token)

		q := []qm.QueryMod{
			qm.Where("locality=?", locality),
			qm.InnerJoin("state s on postcode.state_id = s.id"),
		}
		postcodes, err := models.Postcodes(q...).All(ctx, db)
		if err != nil {
			return nil
		}

		postcodeStates := map[*models.State][]*models.Postcode{}
		for _, postcode := range postcodes {
			state, err := postcode.State().One(ctx, db)
			if err != nil {
				return err
			}

			if postcodes, ok := postcodeStates[state]; ok {
				postcodeStates[state] = append(postcodes, postcode)
			} else {
				postcodeStates[state] = []*models.Postcode{postcode}
			}
		}

		if len(combined) == 2 {
			slices.Reverse(combined)
			s := []string{}
			re := regexp.MustCompile(`(?i)AUS$`)
			for _, token := range combined {
				locality = re.ReplaceAllString(token.ValueString(), "")
				s = append(s, locality)
			}

			q := []qm.QueryMod{
				qm.Where("locality ilike ?", fmt.Sprintf("%s", strings.Join(s, " "))+"%"),
				qm.InnerJoin("state s on postcode.state_id = s.id"),
			}
			postcodes, err := models.Postcodes(q...).All(ctx, db)
			if err != nil {
				return err
			}

			for _, postcode := range postcodes {
				state, err := postcode.State().One(ctx, db)
				if err != nil {
					return err

				}
				if postcodes, ok := postcodeStates[state]; ok {
					postcodeStates[state] = append(postcodes, postcode)
				} else {
					postcodeStates[state] = []*models.Postcode{postcode}
				}

				for _, token := range combined {
					token.SetLocality(true)
				}
			}
		}

		for _, name := range allStatesPreferenced() {
			for k, _ := range postcodeStates {
				if k.Name == name {
					token.SetLocality(true)
					transaction.state = k
					return nil
				}
			}
		}
	}

	return nil
}
