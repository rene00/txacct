package transaction

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"
	"transactionsearch/internal/handlers"
	"transactionsearch/internal/tokenize"
	"transactionsearch/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TransactionOrganisation struct{}

func NewTransactionOrganisation() TransactionHandler {
	return TransactionOrganisation{}
}

type Result struct {
	Similarity                        float64 `boil:"similarity"`
	OrganisationName                  string  `boil:"organisation_name"`
	models.OrganisationStateVic       `boil:",bind"`
	models.OrganisationStateNSW       `boil:",bind"`
	models.OrganisationStateAct       `boil:",bind"`
	models.OrganisationStateQLD       `boil:",bind"`
	models.OrganisationStateNT        `boil:",bind"`
	models.OrganisationStateSa        `boil:",bind"`
	models.OrganisationStateTasmanium `boil:",bind"`
	models.OrganisationStateWa        `boil:",bind"`
	models.Organisation               `boil:",bind"`
	models.Postcode                   `boil:",bind"`
	models.BusinessCode               `boil:",bind"`
	models.State                      `boil:",bind"`
}

func (to TransactionOrganisation) Handle(ctx context.Context, h handlers.Handlers, transaction *Transaction) error {
	var err error

	results := []Result{}

	stateNames := []string{}
	if transaction.state != nil {
		stateNames = append(stateNames, strings.ToLower(transaction.state.Name))
	} else {
		stateNames = append(stateNames, []string{"vic", "nsw"}...)
	}

	names := to.buildNameQueryContents(*transaction)
	for _, name := range names {

		if len(results) >= 1 {
			break
		}

		for _, stateName := range stateNames {
			stateNameTable := fmt.Sprintf("organisation_state_%s", stateName)
			q := []qm.QueryMod{
				qm.Select(
					fmt.Sprintf("%s.id", stateNameTable),
					fmt.Sprintf("%s.name", stateNameTable),
					fmt.Sprintf("%s.name AS organisation_name", stateNameTable),
					fmt.Sprintf("%s.address", stateNameTable),
					fmt.Sprintf("similarity(%s.name, '%s') as similarity", stateNameTable, name),
					"organisation.id",
					"organisation.postcode_id",
					"organisation.business_code_id",
					"postcode.id",
					"postcode.postcode",
					"postcode.locality",
					"postcode.state_id",
					"business_code.id",
					"business_code.code",
					"business_code.description",
					"state.id",
					"state.name",
				),
				qm.From(stateNameTable),
				qm.Where(stateNameTable+".name % ?", name),
				qm.OrderBy(fmt.Sprintf("similarity DESC, %s.name", stateNameTable)),
				qm.InnerJoin(fmt.Sprintf("organisation on organisation.id = %s.organisation_id", stateNameTable)),
				qm.InnerJoin("business_code on business_code.id = organisation.business_code_id"),
				qm.InnerJoin("postcode on postcode.id = organisation.postcode_id"),
				qm.InnerJoin("state on state.id = postcode.state_id"),
			}

			if err = models.NewQuery(q...).Bind(ctx, h.DB, &results); err != nil {
				return fmt.Errorf("failed to query organisation state %s table: %w", stateNameTable, err)
			}

			if len(results) >= 1 {
				break
			}
		}
	}

	resultsOrderBySimilarity := map[float64][]Result{}

	for _, res := range results {
		foundResult, ok := resultsOrderBySimilarity[res.Similarity]
		if !ok {
			resultsOrderBySimilarity[res.Similarity] = []Result{res}
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
	for _, postcode := range transaction.postcodes {
		h.Logger.Debug(fmt.Sprintf("transaction.postcode is (%s) (len:%d) for (%s)", postcode.Postcode, len(transaction.postcodes), transaction.input))
	}

	for _, similarity := range resultsOrderBySimilaritySortedKeys {
		for _, result := range resultsOrderBySimilarity[similarity] {
			for _, postcode := range transaction.postcodes {
				h.Logger.Debug(fmt.Sprintf("comparing postcode (%s) (len:%d) with (%s) for (%s) with similarity (%f)", postcode.Postcode, len(transaction.postcodes), result.Postcode.Postcode, transaction.input, similarity))
				if postcode.Postcode == result.Postcode.Postcode {
					// Perform another select to get organisation with eager
					// loading BusinessCode. The eager loading of BusinessCode
					// can't be done with the previous query bind.
					q := []qm.QueryMod{
						qm.Load("BusinessCode"),
						qm.Load("Postcode"),
						qm.Load("Postcode.State"),
						qm.Where("organisation.id = ?", result.Organisation.ID),
						qm.InnerJoin("postcode on postcode.id = organisation.postcode_id"),
						qm.InnerJoin("state on state.id = postcode.state_id"),
					}

					organisation, err := models.Organisations(q...).One(ctx, h.DB)
					if err != nil {
						return fmt.Errorf("failed to query organisation: %w", err)
					}
					transaction.organisation = organisation
					transaction.state = organisation.R.Postcode.R.State

					// if there are multiple transaction.postcodes that share
					// the same locality, iterate through list and choose
					// lowest postcode to assign to this organisation.
					transaction.postcode = organisation.R.Postcode
					if len(transaction.postcodes) > 0 {
						lowestPostcode := *organisation.R.Postcode
						for _, postcode := range transaction.postcodes {
							if strings.ToLower(postcode.Locality) == strings.ToLower(organisation.R.Postcode.Locality) {
								h.Logger.Debug(fmt.Sprintf("comparing postcode locality (%s,%s) with (%s,%s) for (%s) with similarity (%f)", postcode.Locality, postcode.Postcode, lowestPostcode.Locality, lowestPostcode.Postcode, transaction.input, similarity))
								if postcode.Postcode < lowestPostcode.Postcode {
									h.Logger.Debug("true")
									lowestPostcode = postcode
								}
							}
						}
						transaction.postcode = &lowestPostcode
					}

					if result.OrganisationStateVic != (models.OrganisationStateVic{}) {
						transaction.organisationStateVic = &result.OrganisationStateVic
					} else if result.OrganisationStateNSW != (models.OrganisationStateNSW{}) {
						transaction.organisationStateNSW = &result.OrganisationStateNSW
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
	if token.IsGeo() {
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
