package transaction

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"transactionsearch/internal/tokenize"
	"transactionsearch/models"

	"github.com/datasapiens/cachier"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TransactionOrganisation struct{}

func NewTransactionOrganisation() TransactionHandler {
	return TransactionOrganisation{}
}

func (to TransactionOrganisation) Handle(ctx context.Context, store Store, transaction *Transaction) error {

	if transaction.postcode == nil {
		return nil
	}

	var likeQueryContents []string

	likeQueryContents = to.buildLikeQueryContents(*transaction)
	for i := len(likeQueryContents) - 1; i >= 0; i-- {
		var organisations models.OrganisationSlice
		var q []qm.QueryMod

		v := likeQueryContents[i]
		q = []qm.QueryMod{
			qm.InnerJoin(fmt.Sprintf("postcode p on organisation.postcode_id = %d", transaction.postcode.ID)),
			qm.Where("name ILIKE ?", v+"%"),
			qm.Load("BusinessCode"),
		}

		cachedOrganisations, err := store.Cache.Get(v)
		if err != nil && err != cachier.ErrNotFound {
			return err
		}

		if cachedOrganisations != nil {
			o := *cachedOrganisations
			organisations = o[0]
		} else {
			organisations, err = models.Organisations(q...).All(ctx, store.DB)
			if err != nil {
				return err
			}
			data := []models.OrganisationSlice{organisations}
			err = store.Cache.Set(v, &data)
			if err != nil {
				return err
			}
		}

		for _, organisation := range organisations {
			for i := len(likeQueryContents) - 1; i >= 0; i-- {
				q := likeQueryContents[i]

				re, err := regexp.Compile("(?i)" + q)
				if err != nil {
					return err
				}
				if re.MatchString(organisation.Name) {
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

		if token.Previous() == nil {
			s = append(s, token.ValueString())
			continue
		}

		s = append(s, fmt.Sprintf("%s %s", s[len(s)-1], token.ValueString()))
	}
	return s
}
