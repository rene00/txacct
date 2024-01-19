package dataimporter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"transactionsearch/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type OrganisationStateStorer interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

type OrganisationState interface {
	models.OrganisationStateVic | models.OrganisationStateNSW | models.OrganisationStateTasmanium | models.OrganisationStateAct | models.OrganisationStateQLD | models.OrganisationStateSa | models.OrganisationStateNT | models.OrganisationStateWa
}

type DataImporter[T OrganisationState] struct {
	organisationState T
}

func NewDataImporter[T OrganisationState](o T) DataImporter[T] {
	return DataImporter[T]{organisationState: o}
}

func (d DataImporter[T]) getState(ctx context.Context, db *sql.DB, r Row) (*models.State, error) {
	var state *models.State
	var err error

	stateName := strings.ToUpper(r.GetCellValueWithCol("STATE"))
	if stateName == "" {
		return state, fmt.Errorf("no state found")
	}

	state, err = models.States(qm.Where("name=?", stateName)).One(ctx, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			state := &models.State{Name: stateName}
			if err := state.Insert(ctx, db, boil.Infer()); err != nil {
				switch {
				case IsUniqueConsstraint(err):
					state, err = models.States(qm.Where("name=?", stateName)).One(ctx, db)
					if err != nil {
						return state, fmt.Errorf("failed to find state after insert: %w", err)
					}
				default:
					return state, fmt.Errorf("failed to insert state: %w", err)
				}
			}
			return state, nil
		}
		return state, fmt.Errorf("failed to find state: %w", err)
	}
	return state, nil
}

func (d DataImporter[T]) getPostcode(ctx context.Context, db *sql.DB, r Row) (*models.Postcode, error) {
	var postcode *models.Postcode
	var err error

	rowPostcode := r.GetCellValueWithCol("POSTCODE")
	if rowPostcode == "" {
		return postcode, fmt.Errorf("no postcode found in row")
	}

	rowLocation := r.GetCellValueWithCol("LOCATION")
	if rowLocation == "" {
		return postcode, fmt.Errorf("no location found in row")
	}

	state, err := d.getState(ctx, db, r)
	if err != nil {
		return postcode, err
	}

	postcode, err = models.Postcodes(qm.Where("postcode=? AND locality=? AND state_id=?", rowPostcode, rowLocation, state.ID)).One(ctx, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			postcode := &models.Postcode{
				Postcode: rowPostcode,
				Locality: rowLocation,
				StateID:  state.ID,
			}
			if err := postcode.Insert(ctx, db, boil.Infer()); err != nil {
				switch {
				case IsUniqueConsstraint(err):
					postcode, err = models.Postcodes(qm.Where("postcode=? AND locality=? AND state_id=?", rowPostcode, rowLocation, state.ID)).One(ctx, db)
					if err != nil {
						return postcode, fmt.Errorf("failed to find postcode after insert: %s, %s, %d, %w", rowPostcode, rowLocation, state.ID, err)
					}
				default:
					return postcode, fmt.Errorf("failed to insert postcode: %w", err)
				}
			}
			return postcode, nil
		}
		return postcode, fmt.Errorf("failed to find postcodes: %w", err)
	}
	return postcode, nil
}

func (d DataImporter[T]) getBusinessCode(ctx context.Context, db *sql.DB, r Row) (*models.BusinessCode, error) {
	var businessCode *models.BusinessCode
	var err error

	rowBuscode := r.GetCellValueWithCol("BUSCODE")
	if rowBuscode == "" {
		return businessCode, fmt.Errorf("no business code found in row")
	}

	rowBusinessDescription := r.GetCellValueWithCol("BUSINESS_DESCRIPTION")

	businessCode, err = models.BusinessCodes(qm.Where("code=? AND description=?", rowBuscode, rowBusinessDescription)).One(ctx, db)
	if err != nil && err != sql.ErrNoRows {
		return businessCode, fmt.Errorf("failed finding business code: %w", err)
	}

	if businessCode == nil {
		businessCode = &models.BusinessCode{
			Code:        null.StringFrom(rowBuscode),
			Description: null.StringFrom(rowBusinessDescription),
		}
		if err := businessCode.Insert(ctx, db, boil.Infer()); err != nil {
			switch {
			case IsUniqueConsstraint(err):
				businessCode, err = models.BusinessCodes(qm.Where("code=? AND description=?", rowBuscode, rowBusinessDescription)).One(ctx, db)
				if err != nil && err != sql.ErrNoRows {
					return businessCode, fmt.Errorf("failed finding business code after insert: %w", err)
				}
			default:
				return businessCode, fmt.Errorf("failed to insert business code: %w", err)
			}
		}
	}

	return businessCode, nil
}

func (d DataImporter[T]) getEmail(ctx context.Context, db *sql.DB, r Row) ([]*models.Email, error) {
	var email *models.Email
	var emails []*models.Email
	var err error

	for _, i := range []string{r.GetCellValueWithCol("EMAIL"), r.GetCellValueWithCol("EMAIL-2")} {
		if i == "" {
			continue
		}
		email, err = models.Emails(qm.Where("email=?", i)).One(ctx, db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				email = &models.Email{Email: i}
				if err := email.Insert(ctx, db, boil.Infer()); err != nil {
					switch {
					case IsUniqueConsstraint(err):
						email, err = models.Emails(qm.Where("email=?", i)).One(ctx, db)
						if err != nil {
							return emails, fmt.Errorf("failed finding email after insert: %w", err)
						}
					default:
						return emails, fmt.Errorf("failed finding email 1: %w", err)
					}
				}
				emails = append(emails, email)
				continue
			}
			return emails, fmt.Errorf("failed finding email: %w", err)
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func (d DataImporter[T]) Do(ctx context.Context, db *sql.DB, r Row) error {
	idOrganisation := r.GetCellValueWithCol("ID_ORGANISATION")
	organisationExists, err := models.Organisations(qm.Where("organisation_source_id=?", idOrganisation)).Exists(ctx, db)
	if err != nil {
		return err
	}

	if organisationExists {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	businessCode, err := d.getBusinessCode(ctx, db, r)
	if err != nil {
		return err
	}

	postcode, err := d.getPostcode(ctx, db, r)
	if err != nil {
		return err
	}

	organisationSourceID, err := strconv.Atoi(r.GetCellValueWithCol("ID_ORGANISATION"))
	if err != nil {
		return fmt.Errorf("failed to convert organisation source id")
	}

	if organisationSourceID == 0 {
		return fmt.Errorf("organisation source id is 0")
	}

	organisation := models.Organisation{
		BusinessCodeID:       businessCode.ID,
		PostcodeID:           postcode.ID,
		OrganisationSourceID: organisationSourceID,
	}
	if err := organisation.Insert(ctx, tx, boil.Infer()); err != nil {
		switch {
		case IsUniqueConsstraint(err):
			fmt.Printf("orgisation already imported: %d\n", organisationSourceID)
			return nil
		default:
			return fmt.Errorf("failed to insert organisation: %w", err)
		}
	}

	var organisationStateStore OrganisationStateStorer

	switch any(d.organisationState).(type) {
	case models.OrganisationStateVic:
		organisationStateStore = &models.OrganisationStateVic{
			Name:              r.GetCellValueWithCol("ORGANISATION"),
			OrganisationID:    organisation.ID,
			Abn:               null.StringFrom(r.GetCellValueWithCol("ABN")),
			Address:           null.StringFrom(r.GetCellValueWithCol("ADDRESS")),
			RecordDefunctRisk: null.StringFrom(r.GetCellValueWithCol("RECORD_DEFUNCT_RISK")),
			Region:            null.StringFrom(r.GetCellValueWithCol("REGION")),
			Phone:             null.StringFrom(r.GetCellValueWithCol("PHONE")),
			Mobile:            null.StringFrom(r.GetCellValueWithCol("MOBILE")),
			Freecall:          null.StringFrom(r.GetCellValueWithCol("FREECALL")),
			Fax:               null.StringFrom(r.GetCellValueWithCol("FAX")),
		}
	case models.OrganisationStateNSW:
		organisationStateStore = &models.OrganisationStateNSW{
			Name:              r.GetCellValueWithCol("ORGANISATION"),
			OrganisationID:    organisation.ID,
			Abn:               null.StringFrom(r.GetCellValueWithCol("ABN")),
			Address:           null.StringFrom(r.GetCellValueWithCol("ADDRESS")),
			RecordDefunctRisk: null.StringFrom(r.GetCellValueWithCol("RECORD_DEFUNCT_RISK")),
			Region:            null.StringFrom(r.GetCellValueWithCol("REGION")),
			Phone:             null.StringFrom(r.GetCellValueWithCol("PHONE")),
			Mobile:            null.StringFrom(r.GetCellValueWithCol("MOBILE")),
			Freecall:          null.StringFrom(r.GetCellValueWithCol("FREECALL")),
			Fax:               null.StringFrom(r.GetCellValueWithCol("FAX")),
		}
	case models.OrganisationStateTasmanium:
		organisationStateStore = &models.OrganisationStateTasmanium{
			Name:              r.GetCellValueWithCol("ORGANISATION"),
			OrganisationID:    organisation.ID,
			Abn:               null.StringFrom(r.GetCellValueWithCol("ABN")),
			Address:           null.StringFrom(r.GetCellValueWithCol("ADDRESS")),
			RecordDefunctRisk: null.StringFrom(r.GetCellValueWithCol("RECORD_DEFUNCT_RISK")),
			Region:            null.StringFrom(r.GetCellValueWithCol("REGION")),
			Phone:             null.StringFrom(r.GetCellValueWithCol("PHONE")),
			Mobile:            null.StringFrom(r.GetCellValueWithCol("MOBILE")),
			Freecall:          null.StringFrom(r.GetCellValueWithCol("FREECALL")),
			Fax:               null.StringFrom(r.GetCellValueWithCol("FAX")),
		}
	case models.OrganisationStateAct:
		organisationStateStore = &models.OrganisationStateAct{
			Name:              r.GetCellValueWithCol("ORGANISATION"),
			OrganisationID:    organisation.ID,
			Abn:               null.StringFrom(r.GetCellValueWithCol("ABN")),
			Address:           null.StringFrom(r.GetCellValueWithCol("ADDRESS")),
			RecordDefunctRisk: null.StringFrom(r.GetCellValueWithCol("RECORD_DEFUNCT_RISK")),
			Region:            null.StringFrom(r.GetCellValueWithCol("REGION")),
			Phone:             null.StringFrom(r.GetCellValueWithCol("PHONE")),
			Mobile:            null.StringFrom(r.GetCellValueWithCol("MOBILE")),
			Freecall:          null.StringFrom(r.GetCellValueWithCol("FREECALL")),
			Fax:               null.StringFrom(r.GetCellValueWithCol("FAX")),
		}
	case models.OrganisationStateQLD:
		organisationStateStore = &models.OrganisationStateQLD{
			Name:              r.GetCellValueWithCol("ORGANISATION"),
			OrganisationID:    organisation.ID,
			Abn:               null.StringFrom(r.GetCellValueWithCol("ABN")),
			Address:           null.StringFrom(r.GetCellValueWithCol("ADDRESS")),
			RecordDefunctRisk: null.StringFrom(r.GetCellValueWithCol("RECORD_DEFUNCT_RISK")),
			Region:            null.StringFrom(r.GetCellValueWithCol("REGION")),
			Phone:             null.StringFrom(r.GetCellValueWithCol("PHONE")),
			Mobile:            null.StringFrom(r.GetCellValueWithCol("MOBILE")),
			Freecall:          null.StringFrom(r.GetCellValueWithCol("FREECALL")),
			Fax:               null.StringFrom(r.GetCellValueWithCol("FAX")),
		}
	case models.OrganisationStateNT:
		organisationStateStore = &models.OrganisationStateNT{
			Name:              r.GetCellValueWithCol("ORGANISATION"),
			OrganisationID:    organisation.ID,
			Abn:               null.StringFrom(r.GetCellValueWithCol("ABN")),
			Address:           null.StringFrom(r.GetCellValueWithCol("ADDRESS")),
			RecordDefunctRisk: null.StringFrom(r.GetCellValueWithCol("RECORD_DEFUNCT_RISK")),
			Region:            null.StringFrom(r.GetCellValueWithCol("REGION")),
			Phone:             null.StringFrom(r.GetCellValueWithCol("PHONE")),
			Mobile:            null.StringFrom(r.GetCellValueWithCol("MOBILE")),
			Freecall:          null.StringFrom(r.GetCellValueWithCol("FREECALL")),
			Fax:               null.StringFrom(r.GetCellValueWithCol("FAX")),
		}
	case models.OrganisationStateSa:
		organisationStateStore = &models.OrganisationStateSa{
			Name:              r.GetCellValueWithCol("ORGANISATION"),
			OrganisationID:    organisation.ID,
			Abn:               null.StringFrom(r.GetCellValueWithCol("ABN")),
			Address:           null.StringFrom(r.GetCellValueWithCol("ADDRESS")),
			RecordDefunctRisk: null.StringFrom(r.GetCellValueWithCol("RECORD_DEFUNCT_RISK")),
			Region:            null.StringFrom(r.GetCellValueWithCol("REGION")),
			Phone:             null.StringFrom(r.GetCellValueWithCol("PHONE")),
			Mobile:            null.StringFrom(r.GetCellValueWithCol("MOBILE")),
			Freecall:          null.StringFrom(r.GetCellValueWithCol("FREECALL")),
			Fax:               null.StringFrom(r.GetCellValueWithCol("FAX")),
		}
	case models.OrganisationStateWa:
		organisationStateStore = &models.OrganisationStateWa{
			Name:              r.GetCellValueWithCol("ORGANISATION"),
			OrganisationID:    organisation.ID,
			Abn:               null.StringFrom(r.GetCellValueWithCol("ABN")),
			Address:           null.StringFrom(r.GetCellValueWithCol("ADDRESS")),
			RecordDefunctRisk: null.StringFrom(r.GetCellValueWithCol("RECORD_DEFUNCT_RISK")),
			Region:            null.StringFrom(r.GetCellValueWithCol("REGION")),
			Phone:             null.StringFrom(r.GetCellValueWithCol("PHONE")),
			Mobile:            null.StringFrom(r.GetCellValueWithCol("MOBILE")),
			Freecall:          null.StringFrom(r.GetCellValueWithCol("FREECALL")),
			Fax:               null.StringFrom(r.GetCellValueWithCol("FAX")),
		}
	default:
		return fmt.Errorf("not supported")
	}

	if err := organisationStateStore.Insert(ctx, tx, boil.Infer()); err != nil {
		switch {
		case IsUniqueConsstraint(err):
			tx.Rollback()
			return nil
		default:
			return fmt.Errorf("failed to insert organisation state: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
	}

	emails, err := d.getEmail(ctx, db, r)
	if err != nil {
		return err
	}

	for _, email := range emails {
		emailOrganisation := models.EmailOrganisation{
			EmailID:        email.ID,
			OrganisationID: organisation.ID,
		}
		if err := emailOrganisation.Insert(ctx, db, boil.Infer()); err != nil {
			switch {
			case IsUniqueConsstraint(err):
				return nil
			default:
				return fmt.Errorf("failed inserting email: %w", err)
			}
		}
	}

	return nil
}

func IsUniqueConsstraint(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}
