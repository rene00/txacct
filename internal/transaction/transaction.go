package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"transactionsearch/internal/handlers"
	"transactionsearch/internal/tokenize"
	"transactionsearch/models"
)

type TransactionHandler interface {
	Handle(ctx context.Context, h handlers.Handlers, transaction *Transaction) error
}

type Transaction struct {
	db                   *sql.DB
	input                string
	tokenize             tokenize.Tokenize
	state                *models.State
	postcode             *models.Postcode
	postcodes            []models.Postcode
	organisation         *models.Organisation
	organisationStateVic *models.OrganisationStateVic
	organisationStateNSW *models.OrganisationStateNSW
	organisationStateQLD *models.OrganisationStateQLD
	organisationStateAct *models.OrganisationStateAct
	organisationStateTas *models.OrganisationStateTasmanium
	organisationStateSa  *models.OrganisationStateSa
	organisationStateNT  *models.OrganisationStateNT
	organisationStateWa  *models.OrganisationStateWa
}

func NewTransaction(s string, db *sql.DB) *Transaction {
	t := tokenize.NewTokenize()
	t.Parse(s)
	return &Transaction{
		input:    s,
		tokenize: t,
		db:       db,
	}
}

// TransactionJSONResponse is the response given back by the router for the
// transaction endpoint.
type TransactionJSONResponse struct {
	Memo         string  `json:"memo"`
	State        *string `json:"state"`
	Postcode     *int    `json:"postcode"`
	Organisation *string `json:"organisation"`
	// Address is from Organisation.Address
	Address *string `json:"address"`
	// Description is from Organisation.BusinessCode.Description
	Description *string `json:"description"`
}

func (t TransactionJSONResponse) GetState() string {
	if t.State == nil {
		return ""
	}
	return *t.State
}

func (t TransactionJSONResponse) GetPostcode() int {
	if t.Postcode == nil {
		return 0
	}
	return *t.Postcode
}

func (t TransactionJSONResponse) GetOrganisation() string {
	if t.Organisation == nil {
		return ""
	}
	return *t.Organisation
}

func (t TransactionJSONResponse) GetAddress() string {
	if t.Address == nil {
		return ""
	}
	return *t.Address
}

func (t TransactionJSONResponse) GetDescription() string {
	if t.Description == nil {
		return ""
	}
	return *t.Description
}

type TransactionJSONRequest struct {
	Memo string `json:"memo" binding:"required"`
}

func NewTransactionJSONResponse(t Transaction) (TransactionJSONResponse, error) {
	r := TransactionJSONResponse{Memo: t.input}

	if t.state != nil {
		state := t.state.Name
		r.State = &state
	}

	if t.postcode != nil {
		postcode, err := strconv.Atoi(t.postcode.Postcode)
		if err != nil {
			return r, err
		}
		r.Postcode = &postcode
	}

	address := ""
	organisation := ""
	switch {
	case t.organisationStateVic != nil:
		address = fmt.Sprintf("%s, %s, %s", t.organisationStateVic.Address.String, t.organisation.R.Postcode.Locality, t.organisation.R.Postcode.R.State.Name)
		organisation = t.organisationStateVic.Name
	case t.organisationStateNSW != nil:
		address = fmt.Sprintf("%s, %s, %s", t.organisationStateNSW.Address.String, t.organisation.R.Postcode.Locality, t.organisation.R.Postcode.R.State.Name)
		organisation = t.organisationStateNSW.Name
	case t.organisationStateQLD != nil:
		address = fmt.Sprintf("%s, %s, %s", t.organisationStateQLD.Address.String, t.organisation.R.Postcode.Locality, t.organisation.R.Postcode.R.State.Name)
		organisation = t.organisationStateQLD.Name
	case t.organisationStateAct != nil:
		address = fmt.Sprintf("%s, %s, %s", t.organisationStateAct.Address.String, t.organisation.R.Postcode.Locality, t.organisation.R.Postcode.R.State.Name)
		organisation = t.organisationStateAct.Name
	case t.organisationStateTas != nil:
		address = fmt.Sprintf("%s, %s, %s", t.organisationStateTas.Address.String, t.organisation.R.Postcode.Locality, t.organisation.R.Postcode.R.State.Name)
		organisation = t.organisationStateTas.Name
	case t.organisationStateNT != nil:
		address = fmt.Sprintf("%s, %s, %s", t.organisationStateNT.Address.String, t.organisation.R.Postcode.Locality, t.organisation.R.Postcode.R.State.Name)
		organisation = t.organisationStateNT.Name
	case t.organisationStateSa != nil:
		address = fmt.Sprintf("%s, %s, %s", t.organisationStateSa.Address.String, t.organisation.R.Postcode.Locality, t.organisation.R.Postcode.R.State.Name)
		organisation = t.organisationStateSa.Name
	case t.organisationStateWa != nil:
		address = fmt.Sprintf("%s, %s, %s", t.organisationStateWa.Address.String, t.organisation.R.Postcode.Locality, t.organisation.R.Postcode.R.State.Name)
		organisation = t.organisationStateWa.Name
	default:
		return r, nil
	}

	r.Address = &address
	r.Organisation = &organisation

	if t.organisation != nil && t.organisation.R != nil && t.organisation.R.BusinessCode != nil {
		description := t.organisation.R.BusinessCode.Description.String
		r.Description = &description
	}

	return r, nil
}

// stripCountryFromString removes leading "AUS" from string.
func stripCountryFromString(s string) string {
	re := regexp.MustCompile(`(?i)AUS$`)
	return re.ReplaceAllString(s, "")
}

// stripStateFromString removes leading state name from string.
func stripStateFromString(ctx context.Context, db *sql.DB, s string) (string, error) {
	newString := s

	states, err := models.States().All(ctx, db)
	if err != nil {
		return newString, err
	}

	for _, state := range states {
		re := regexp.MustCompile(fmt.Sprintf(`(i?)%s$`, state.Name))
		newString = re.ReplaceAllString(s, "")
		if newString != s {
			break
		}
	}

	return newString, nil
}

func allStatesPreferenced() []string {
	return []string{"VIC", "NSW", "QLD", "SA", "ACT", "WA", "TAS", "NT"}
}
