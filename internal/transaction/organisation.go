package transaction

import (
	"context"
	"fmt"
	"slices"
	"sort"
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

	type result struct {
		Similarity          float64 `boil:"similarity"`
		models.Organisation `boil:",bind"`
		models.Postcode     `boil:",bind"`
		models.BusinessCode `boil:",bind"`
	}

	results := []result{}

	foundName := ""

	names := to.buildNameQueryContents(*transaction)
	for _, name := range names {
		cacheOrganisation, err := store.Cache.Get(name)
		if err != nil && err != cachier.ErrNotFound {
			return err
		}

		if cacheOrganisation != nil {
			transaction.organisation = cacheOrganisation
			return nil
		}
		if err = models.NewQuery(
			qm.Select("organisation.id", "organisation.name", "organisation.address", "organisation.postcode_id", fmt.Sprintf("similarity(name, '%s') as similarity", name), "postcode.id", "postcode.locality", "business_code.id"),
			qm.From("organisation"),
			qm.Where("organisation.name % ?", name),
			qm.OrderBy("similarity DESC, organisation.name"),
			qm.InnerJoin("postcode on postcode.id = organisation.postcode_id"),
			qm.InnerJoin("business_code on business_code.id = organisation.business_code_id"),
		).Bind(ctx, store.DB, &results); err != nil {
			return err
		}
		if len(results) >= 1 {
			foundName = name
			break
		}
	}

	resultsOrderBySimilarity := map[float64][]result{}

	for _, res := range results {
		foundResult, ok := resultsOrderBySimilarity[res.Similarity]
		if !ok {
			resultsOrderBySimilarity[res.Similarity] = []result{res}
			continue
		}
		foundResult = append(foundResult, res)
		resultsOrderBySimilarity[res.Similarity] = foundResult
	}

	resultsOrderBySimilaritySortedKeys := make([]float64, 0, len(resultsOrderBySimilarity))

	for k := range resultsOrderBySimilarity {
		resultsOrderBySimilaritySortedKeys = append(resultsOrderBySimilaritySortedKeys, k)
	}
	sort.Float64s(resultsOrderBySimilaritySortedKeys)
	slices.Reverse(resultsOrderBySimilaritySortedKeys)

	// Iterate through all organisations order by similarity desc and set the
	// transaction organisation of the first organisation that matches the
	// transacode postcode.
	for _, similarity := range resultsOrderBySimilaritySortedKeys {
		for _, result := range resultsOrderBySimilarity[similarity] {
			for _, postcode := range transaction.postcodes {
				if postcode.ID == result.Postcode.ID {
					// Perform another select to get organisation with eager
					// loading BusinessCode. The eager loading of BusinessCode
					// can't be done with the previous query bind.
					organisation, err := models.Organisations(qm.Load("BusinessCode"), qm.Load("Postcode"), qm.Where("id = ?", result.Organisation.ID)).One(ctx, store.DB)
					if err != nil {
						return err
					}
					transaction.organisation = organisation

					if err = store.Cache.Set(foundName, organisation); err != nil {
						return err
					}

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

// buildNameQueryContents accepts a transaction and returns a slice of strings
// which contain Organisation.name queries for the transaction.
func (to TransactionOrganisation) buildNameQueryContents(transaction Transaction) []string {
	var s []string
	for _, token := range transaction.tokenize.Tokens() {
		if to.querySkipToken(*token) {
			continue
		}

		if token.Previous() == nil {
			s = append(s, token.ValueString())
			continue
		}

		if len(s) == 0 {
			s = append(s, token.ValueString())
		} else {
			s = append(s, fmt.Sprintf("%s %s", s[len(s)-1], token.ValueString()))
		}
	}
	slices.Reverse(s)
	return s
}
