package transaction

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
	"transactionsearch/internal/handlers"
	"transactionsearch/internal/tokenize"
	"transactionsearch/models"
	"unicode"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TransactionPostcode struct{}

func NewTransactionPostcode() TransactionHandler {
	return TransactionPostcode{}
}

type PostcodeResult struct {
	Similarity      float64 `boil:"similarity"`
	models.Postcode `boil:",bind"`
}

func (tp TransactionPostcode) Handle(ctx context.Context, h handlers.Handlers, transaction *Transaction) error {
	var postcodes []models.Postcode
	var tokens []*tokenize.Token
	var err error

	for idx, token := range transaction.tokenize.TokensReversed() {

		h.Logger.Debug(fmt.Sprintf("process token (%s) in (%s)(%d)", token.ValueString(), transaction.input, idx))

		if idx == 2 {
			break
		}

		if strings.ToLower(token.ValueString()) == "aus" {
			h.Logger.Debug(fmt.Sprintf("skipping token (%s) in (%s) because country match", token.ValueString(), transaction.input))
			continue
		}

		re := regexp.MustCompile("[0-9]")
		if re.FindString(token.ValueString()) != "" {
			h.Logger.Debug(fmt.Sprintf("skipping token (%s) in (%s) because digit match", token.ValueString(), transaction.input))
			continue
		}

		// token contains only state which will be skipped.
		if token.IsState() && !token.IsLocality() {
			h.Logger.Debug(fmt.Sprintf("skipping token (%s) in (%s) because state match and is not locality", token.ValueString(), transaction.input))
			continue
		}

		h.Logger.Debug(fmt.Sprintf("including token (%s) in postcode slice", token.ValueString()))
		tokens = append(tokens, token)
	}

	// If length of tokens is 1 and state is set
	// if len(tokens) == 1 && transaction.state != nil {

	allResults := []PostcodeResult{}

	firstQuery := ""

	for idx, token := range tokens {
		results := []PostcodeResult{}

		q := []qm.QueryMod{
			qm.From("postcode"),
			qm.OrderBy("similarity DESC, postcode.postcode ASC"),
			qm.InnerJoin("state on postcode.state_id = state.id"),
		}

		qSelect := []qm.QueryMod{}
		qWhere := []qm.QueryMod{}

		tokensToMarkAsLocality := []*tokenize.Token{}

		// If this is the first token in the tokens slice, use where clause on
		// that. If it's the second then concat the first and second token.
		if idx == 0 {
			s, err := stripStateFromString(ctx, h.DB, token.ValueString())
			if err != nil {
				return err
			}
			s = stripCountryFromString(s)
			h.Logger.Debug(fmt.Sprintf("single token (%s) query is (%s) for (%s)", token.ValueString(), s, transaction.input))
			qSelect = append(qSelect, qm.Select(
				"postcode.id",
				"postcode.locality",
				"postcode.postcode",
				"postcode.state_id",
				fmt.Sprintf("similarity(postcode.locality, '%s') as similarity", s),
				"state.id",
				"state.name",
			))
			qWhere = append(qWhere, qm.Where("postcode.locality % ?", s))
			firstQuery = s
			tokensToMarkAsLocality = append(tokensToMarkAsLocality, token)
		} else {
			q := ""

			previousTokenInSlice := tokens[idx-1]
			s1, err := stripStateFromString(ctx, h.DB, previousTokenInSlice.ValueString())
			if err != nil {
				return err
			}
			s1 = stripCountryFromString(s1)

			s2, err := stripStateFromString(ctx, h.DB, token.ValueString())
			if err != nil {
				return err
			}
			s2 = stripCountryFromString(s2)

			// both tokens contain not strings which can be localities so skip this token.
			if len(s1) == 0 && len(s2) == 0 {
				continue
			}

			if len(s1) != 0 {
				q = s1
			}

			if len(s2) != 0 {
				q = fmt.Sprintf("%s %s", s2, s1)
			}

			q = strings.TrimRight(q, " ")
			q = strings.TrimLeft(q, " ")

			// if first query is the same as this and has returned results, dont perform a combined query.
			if firstQuery == q && len(allResults) >= 1 {
				h.Logger.Debug(fmt.Sprintf("skipping query (%s) (first:%s) (allResults len:%d) for (%s)", q, firstQuery, len(allResults), transaction.input))
				continue
			}

			h.Logger.Debug(fmt.Sprintf("combined query is (%s) for (%s)", q, transaction.input))

			qSelect = append(qSelect, qm.Select(
				"postcode.id",
				"postcode.locality",
				"postcode.postcode",
				"postcode.state_id",
				fmt.Sprintf("similarity(postcode.locality, '%s') as similarity", q),
				"state.id",
				"state.name",
			))
			qWhere = append(qWhere, qm.Where("postcode.locality % ?", q))
			tokensToMarkAsLocality = append(tokensToMarkAsLocality, previousTokenInSlice)
			tokensToMarkAsLocality = append(tokensToMarkAsLocality, token)
		}

		if transaction.state != nil {
			qWhere = append(qWhere, qm.And("postcode.state_id=?", transaction.state.ID))
		}

		q = append(q, qSelect...)
		q = append(q, qWhere...)

		if err = models.NewQuery(q...).Bind(ctx, h.DB, &results); err != nil {
			return fmt.Errorf("failed to query postcode table: %w", err)
		}

		if len(results) > 0 {
			for _, token := range tokensToMarkAsLocality {
				h.Logger.Debug(fmt.Sprintf("setting locality for token (%s) in (%s)", token.ValueString(), transaction.input))
				token.SetLocality(true)
			}
		}

		allResults = append(allResults, results...)
	}

	if len(allResults) == 0 {
		return nil
	}

	resultsOrderBySimilarity := map[float64][]PostcodeResult{}

	for _, res := range allResults {
		// filter postcode to remove bad postcodes
		if !filterPostcodeResult(res) {
			continue
		}

		foundResult, ok := resultsOrderBySimilarity[res.Similarity]
		if !ok {
			resultsOrderBySimilarity[res.Similarity] = []PostcodeResult{res}
			continue
		}
		foundResult = append(foundResult, res)
		resultsOrderBySimilarity[res.Similarity] = foundResult
		h.Logger.Debug(fmt.Sprintf("adding postcode (%s) to resultsOrderBySimilarity[%f]", res.Postcode.Postcode, res.Similarity))
	}

	resultsOrderBySimilaritySortedKeys := make([]float64, 0, len(resultsOrderBySimilarity))
	for k := range resultsOrderBySimilarity {
		resultsOrderBySimilaritySortedKeys = append(resultsOrderBySimilaritySortedKeys, k)
	}
	sort.Float64s(resultsOrderBySimilaritySortedKeys)
	slices.Reverse(resultsOrderBySimilaritySortedKeys)

	// Take highest ordered results
	highest := resultsOrderBySimilaritySortedKeys[0]
	h2 := resultsOrderBySimilarity[highest]

	// Take first result in highest result and use that as postcode
	h1 := h2[0]
	transaction.postcode = &h1.Postcode

	for _, postcodeResult := range h2 {
		h.Logger.Debug(fmt.Sprintf("appending found postcode (%s, %s) to postcodes for (%s)", postcodeResult.Postcode.Locality, postcodeResult.Postcode.Postcode, transaction.input))
		postcodes = append(postcodes, postcodeResult.Postcode)
	}
	transaction.postcodes = postcodes

	for _, postcode := range postcodes {
		h.Logger.Debug(fmt.Sprintf("postcode2 found (%s, %s) for (%s)", postcode.Locality, postcode.Postcode, transaction.input))
	}

	for _, postcode := range transaction.postcodes {
		h.Logger.Debug(fmt.Sprintf("postcode found (%s, %s) for (%s)", postcode.Locality, postcode.Postcode, transaction.input))
	}

	/*
		for idx, token := range transaction.tokenize.TokensReversed() {

			// Only iterate through last 2 tokens as they will have postcode locality.
			if idx == 2 {
				break
			}

			if strings.ToLower(token.ValueString()) == "aus" {
				continue
			}

			re := regexp.MustCompile("[0-9]")
			if re.FindString(token.ValueString()) != "" {
				continue
			}

			q := []qm.QueryMod{
				qm.Where("lower(locality) = lower(?)", token.ValueString()),
				qm.InnerJoin("state s on postcode.state_id = s.id"),
			}

			s, err := stripStateFromString(ctx, h.DB, token.ValueString())
			if err != nil {
				return err
			}

			s = stripCountryFromString(s)
			if s == "" {
				continue
			}

			previousToken := token.Previous()

			if previousToken != nil && idx == 0 {
				previousTokenString, err := stripStateFromString(ctx, h.DB, token.Previous().ValueString())
				if err != nil {
					return err
				}
				previousTokenString = stripCountryFromString(previousTokenString)
				query := previousTokenString
				if len(s) != 0 {
					query = fmt.Sprintf("%s %s", previousTokenString, s)
				}
				q = append(q, qm.Or("locality ilike ?", query+"%"))
			}

			postcodeSlice, err := models.Postcodes(q...).All(ctx, h.DB)
			if err != nil {
				return err
			}

			for _, postcode := range postcodeSlice {
				found := false
				for _, foundPostcode := range transaction.postcodes {
					if postcode.ID == foundPostcode.ID {
						found = true
						break
					}
				}
				if !found {
					transaction.postcodes = append(transaction.postcodes, postcode)
				}
			}

			token.SetLocality(true)

			if previousToken != nil {
				for _, i := range postcodeSlice {
					if strings.ToLower(previousToken.ValueString()) == strings.ToLower(i.Locality) {
						previousToken.SetLocality(true)
						break
					}
				}
			}

			// If 1 or more postcodes found, set the current token locality to
			// true. This ensures that it is excluded when searching for
			// organisation.
			if len(postcodeSlice) == 0 {
				token.SetLocality(false)
				if previousToken != nil {
					previousToken.SetLocality(false)
				}
			}

			// If a single postcode was found, assign the state of this postcode to
			// the transaction.
			if len(postcodes) == 1 {
				postcode := postcodes[0]
				state, err := postcode.State().One(ctx, h.DB)
				if err != nil {
					return err
				}
				transaction.postcode = postcode
				transaction.state = state
			}

		}

	*/

	/*
		for _, i := range transaction.postcodes {
			fmt.Printf("DEBUG1a:%s, %d\n", i.Locality, i.ID)
		}
	*/

	return nil
}

func filterPostcodeResult(res PostcodeResult) bool {
	postcode := res.Postcode
	firstChar := rune(postcode.Locality[0])
	if !unicode.IsLetter(firstChar) {
		return false
	}

	lastChar := rune(postcode.Locality[len(postcode.Locality)-1])
	if !unicode.IsLetter(lastChar) {
		return false
	}
	return true
}
