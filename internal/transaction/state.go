package transaction

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"transactionsearch/internal/handlers"
	"transactionsearch/internal/tokenize"
	"transactionsearch/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TransactionState struct{}

func NewTransactionState() TransactionHandler {
	return TransactionState{}
}

func (ts TransactionState) Handle(ctx context.Context, h handlers.Handlers, transaction *Transaction) error {
	states, err := models.States().All(ctx, h.DB)
	if err != nil {
		return err
	}

	// Get state from last token. If state exists within the memo and is not
	// the last token, it is most likely part of the organisation name.
	lastToken := transaction.tokenize.Last()
	for _, state := range states {
		shortState := state.Name[:2]
		re := regexp.MustCompile(fmt.Sprintf(`(?i)(?P<prefix>\w+)?(?P<state>%s|%s)$`, state.Name, shortState))
		match := re.FindStringSubmatch(lastToken.ValueString())
		if len(match) == 3 {
			lastToken.SetState(true)

			// if last token not exactly same as state name then mark as locality
			if strings.ToLower(lastToken.ValueString()) != strings.ToLower(state.Name) {
				lastToken.SetLocality(true)
			}

			transaction.state = state
			return nil
		}
	}

	// Lookup state using postcode. Get the last 2 tokens and lookup in
	// postcode. If postcode found, use the state from the postcode.
	combined := []*tokenize.Token{}
	tokens := transaction.tokenize.TokensReversed()
	for idx, token := range tokens {
		if idx == 2 {
			break
		}

		if token.IsCountry() && !token.IsLocality() {
			continue
		}

		re := regexp.MustCompile(`(?i)AUS$`)
		locality := re.ReplaceAllString(token.ValueString(), "")

		combined = append(combined, token)

		q := []qm.QueryMod{
			qm.Where("locality=?", locality),
			qm.InnerJoin("state s on postcode.state_id = s.id"),
		}
		postcodes, err := models.Postcodes(q...).All(ctx, h.DB)
		if err != nil {
			return nil
		}

		postcodeStates := map[*models.State][]*models.Postcode{}
		for _, postcode := range postcodes {
			state, err := postcode.State().One(ctx, h.DB)
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

				// if token contains .*AUS$ set token to Country but it will also be locality.
				if strings.ToLower(locality) != strings.ToLower(token.ValueString()) {
					token.SetCountry(true)
				}
			}

			q := []qm.QueryMod{
				qm.Where("locality ilike ?", fmt.Sprintf("%s", strings.Join(s, " "))+"%"),
				qm.InnerJoin("state s on postcode.state_id = s.id"),
			}
			postcodes, err := models.Postcodes(q...).All(ctx, h.DB)
			if err != nil {
				return err
			}

			for _, postcode := range postcodes {
				state, err := postcode.State().One(ctx, h.DB)
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
