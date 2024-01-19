// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// OrganisationStateSa is an object representing the database table.
type OrganisationStateSa struct {
	ID                int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	OrganisationID    int         `boil:"organisation_id" json:"organisation_id" toml:"organisation_id" yaml:"organisation_id"`
	Name              string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Abn               null.String `boil:"abn" json:"abn,omitempty" toml:"abn" yaml:"abn,omitempty"`
	Address           null.String `boil:"address" json:"address,omitempty" toml:"address" yaml:"address,omitempty"`
	RecordDefunctRisk null.String `boil:"record_defunct_risk" json:"record_defunct_risk,omitempty" toml:"record_defunct_risk" yaml:"record_defunct_risk,omitempty"`
	Region            null.String `boil:"region" json:"region,omitempty" toml:"region" yaml:"region,omitempty"`
	Phone             null.String `boil:"phone" json:"phone,omitempty" toml:"phone" yaml:"phone,omitempty"`
	Mobile            null.String `boil:"mobile" json:"mobile,omitempty" toml:"mobile" yaml:"mobile,omitempty"`
	Freecall          null.String `boil:"freecall" json:"freecall,omitempty" toml:"freecall" yaml:"freecall,omitempty"`
	Fax               null.String `boil:"fax" json:"fax,omitempty" toml:"fax" yaml:"fax,omitempty"`

	R *organisationStateSaR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L organisationStateSaL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var OrganisationStateSaColumns = struct {
	ID                string
	OrganisationID    string
	Name              string
	Abn               string
	Address           string
	RecordDefunctRisk string
	Region            string
	Phone             string
	Mobile            string
	Freecall          string
	Fax               string
}{
	ID:                "id",
	OrganisationID:    "organisation_id",
	Name:              "name",
	Abn:               "abn",
	Address:           "address",
	RecordDefunctRisk: "record_defunct_risk",
	Region:            "region",
	Phone:             "phone",
	Mobile:            "mobile",
	Freecall:          "freecall",
	Fax:               "fax",
}

var OrganisationStateSaTableColumns = struct {
	ID                string
	OrganisationID    string
	Name              string
	Abn               string
	Address           string
	RecordDefunctRisk string
	Region            string
	Phone             string
	Mobile            string
	Freecall          string
	Fax               string
}{
	ID:                "organisation_state_sa.id",
	OrganisationID:    "organisation_state_sa.organisation_id",
	Name:              "organisation_state_sa.name",
	Abn:               "organisation_state_sa.abn",
	Address:           "organisation_state_sa.address",
	RecordDefunctRisk: "organisation_state_sa.record_defunct_risk",
	Region:            "organisation_state_sa.region",
	Phone:             "organisation_state_sa.phone",
	Mobile:            "organisation_state_sa.mobile",
	Freecall:          "organisation_state_sa.freecall",
	Fax:               "organisation_state_sa.fax",
}

// Generated where

var OrganisationStateSaWhere = struct {
	ID                whereHelperint
	OrganisationID    whereHelperint
	Name              whereHelperstring
	Abn               whereHelpernull_String
	Address           whereHelpernull_String
	RecordDefunctRisk whereHelpernull_String
	Region            whereHelpernull_String
	Phone             whereHelpernull_String
	Mobile            whereHelpernull_String
	Freecall          whereHelpernull_String
	Fax               whereHelpernull_String
}{
	ID:                whereHelperint{field: "\"organisation_state_sa\".\"id\""},
	OrganisationID:    whereHelperint{field: "\"organisation_state_sa\".\"organisation_id\""},
	Name:              whereHelperstring{field: "\"organisation_state_sa\".\"name\""},
	Abn:               whereHelpernull_String{field: "\"organisation_state_sa\".\"abn\""},
	Address:           whereHelpernull_String{field: "\"organisation_state_sa\".\"address\""},
	RecordDefunctRisk: whereHelpernull_String{field: "\"organisation_state_sa\".\"record_defunct_risk\""},
	Region:            whereHelpernull_String{field: "\"organisation_state_sa\".\"region\""},
	Phone:             whereHelpernull_String{field: "\"organisation_state_sa\".\"phone\""},
	Mobile:            whereHelpernull_String{field: "\"organisation_state_sa\".\"mobile\""},
	Freecall:          whereHelpernull_String{field: "\"organisation_state_sa\".\"freecall\""},
	Fax:               whereHelpernull_String{field: "\"organisation_state_sa\".\"fax\""},
}

// OrganisationStateSaRels is where relationship names are stored.
var OrganisationStateSaRels = struct {
	Organisation string
}{
	Organisation: "Organisation",
}

// organisationStateSaR is where relationships are stored.
type organisationStateSaR struct {
	Organisation *Organisation `boil:"Organisation" json:"Organisation" toml:"Organisation" yaml:"Organisation"`
}

// NewStruct creates a new relationship struct
func (*organisationStateSaR) NewStruct() *organisationStateSaR {
	return &organisationStateSaR{}
}

func (r *organisationStateSaR) GetOrganisation() *Organisation {
	if r == nil {
		return nil
	}
	return r.Organisation
}

// organisationStateSaL is where Load methods for each relationship are stored.
type organisationStateSaL struct{}

var (
	organisationStateSaAllColumns            = []string{"id", "organisation_id", "name", "abn", "address", "record_defunct_risk", "region", "phone", "mobile", "freecall", "fax"}
	organisationStateSaColumnsWithoutDefault = []string{"organisation_id", "name"}
	organisationStateSaColumnsWithDefault    = []string{"id", "abn", "address", "record_defunct_risk", "region", "phone", "mobile", "freecall", "fax"}
	organisationStateSaPrimaryKeyColumns     = []string{"id"}
	organisationStateSaGeneratedColumns      = []string{}
)

type (
	// OrganisationStateSaSlice is an alias for a slice of pointers to OrganisationStateSa.
	// This should almost always be used instead of []OrganisationStateSa.
	OrganisationStateSaSlice []*OrganisationStateSa
	// OrganisationStateSaHook is the signature for custom OrganisationStateSa hook methods
	OrganisationStateSaHook func(context.Context, boil.ContextExecutor, *OrganisationStateSa) error

	organisationStateSaQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	organisationStateSaType                 = reflect.TypeOf(&OrganisationStateSa{})
	organisationStateSaMapping              = queries.MakeStructMapping(organisationStateSaType)
	organisationStateSaPrimaryKeyMapping, _ = queries.BindMapping(organisationStateSaType, organisationStateSaMapping, organisationStateSaPrimaryKeyColumns)
	organisationStateSaInsertCacheMut       sync.RWMutex
	organisationStateSaInsertCache          = make(map[string]insertCache)
	organisationStateSaUpdateCacheMut       sync.RWMutex
	organisationStateSaUpdateCache          = make(map[string]updateCache)
	organisationStateSaUpsertCacheMut       sync.RWMutex
	organisationStateSaUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var organisationStateSaAfterSelectHooks []OrganisationStateSaHook

var organisationStateSaBeforeInsertHooks []OrganisationStateSaHook
var organisationStateSaAfterInsertHooks []OrganisationStateSaHook

var organisationStateSaBeforeUpdateHooks []OrganisationStateSaHook
var organisationStateSaAfterUpdateHooks []OrganisationStateSaHook

var organisationStateSaBeforeDeleteHooks []OrganisationStateSaHook
var organisationStateSaAfterDeleteHooks []OrganisationStateSaHook

var organisationStateSaBeforeUpsertHooks []OrganisationStateSaHook
var organisationStateSaAfterUpsertHooks []OrganisationStateSaHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *OrganisationStateSa) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *OrganisationStateSa) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *OrganisationStateSa) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *OrganisationStateSa) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *OrganisationStateSa) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *OrganisationStateSa) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *OrganisationStateSa) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *OrganisationStateSa) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *OrganisationStateSa) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateSaAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddOrganisationStateSaHook registers your hook function for all future operations.
func AddOrganisationStateSaHook(hookPoint boil.HookPoint, organisationStateSaHook OrganisationStateSaHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		organisationStateSaAfterSelectHooks = append(organisationStateSaAfterSelectHooks, organisationStateSaHook)
	case boil.BeforeInsertHook:
		organisationStateSaBeforeInsertHooks = append(organisationStateSaBeforeInsertHooks, organisationStateSaHook)
	case boil.AfterInsertHook:
		organisationStateSaAfterInsertHooks = append(organisationStateSaAfterInsertHooks, organisationStateSaHook)
	case boil.BeforeUpdateHook:
		organisationStateSaBeforeUpdateHooks = append(organisationStateSaBeforeUpdateHooks, organisationStateSaHook)
	case boil.AfterUpdateHook:
		organisationStateSaAfterUpdateHooks = append(organisationStateSaAfterUpdateHooks, organisationStateSaHook)
	case boil.BeforeDeleteHook:
		organisationStateSaBeforeDeleteHooks = append(organisationStateSaBeforeDeleteHooks, organisationStateSaHook)
	case boil.AfterDeleteHook:
		organisationStateSaAfterDeleteHooks = append(organisationStateSaAfterDeleteHooks, organisationStateSaHook)
	case boil.BeforeUpsertHook:
		organisationStateSaBeforeUpsertHooks = append(organisationStateSaBeforeUpsertHooks, organisationStateSaHook)
	case boil.AfterUpsertHook:
		organisationStateSaAfterUpsertHooks = append(organisationStateSaAfterUpsertHooks, organisationStateSaHook)
	}
}

// One returns a single organisationStateSa record from the query.
func (q organisationStateSaQuery) One(ctx context.Context, exec boil.ContextExecutor) (*OrganisationStateSa, error) {
	o := &OrganisationStateSa{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for organisation_state_sa")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all OrganisationStateSa records from the query.
func (q organisationStateSaQuery) All(ctx context.Context, exec boil.ContextExecutor) (OrganisationStateSaSlice, error) {
	var o []*OrganisationStateSa

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to OrganisationStateSa slice")
	}

	if len(organisationStateSaAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all OrganisationStateSa records in the query.
func (q organisationStateSaQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count organisation_state_sa rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q organisationStateSaQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if organisation_state_sa exists")
	}

	return count > 0, nil
}

// Organisation pointed to by the foreign key.
func (o *OrganisationStateSa) Organisation(mods ...qm.QueryMod) organisationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.OrganisationID),
	}

	queryMods = append(queryMods, mods...)

	return Organisations(queryMods...)
}

// LoadOrganisation allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (organisationStateSaL) LoadOrganisation(ctx context.Context, e boil.ContextExecutor, singular bool, maybeOrganisationStateSa interface{}, mods queries.Applicator) error {
	var slice []*OrganisationStateSa
	var object *OrganisationStateSa

	if singular {
		var ok bool
		object, ok = maybeOrganisationStateSa.(*OrganisationStateSa)
		if !ok {
			object = new(OrganisationStateSa)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeOrganisationStateSa)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeOrganisationStateSa))
			}
		}
	} else {
		s, ok := maybeOrganisationStateSa.(*[]*OrganisationStateSa)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeOrganisationStateSa)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeOrganisationStateSa))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &organisationStateSaR{}
		}
		args = append(args, object.OrganisationID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &organisationStateSaR{}
			}

			for _, a := range args {
				if a == obj.OrganisationID {
					continue Outer
				}
			}

			args = append(args, obj.OrganisationID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`organisation`),
		qm.WhereIn(`organisation.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Organisation")
	}

	var resultSlice []*Organisation
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Organisation")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for organisation")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for organisation")
	}

	if len(organisationAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Organisation = foreign
		if foreign.R == nil {
			foreign.R = &organisationR{}
		}
		foreign.R.OrganisationStateSa = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.OrganisationID == foreign.ID {
				local.R.Organisation = foreign
				if foreign.R == nil {
					foreign.R = &organisationR{}
				}
				foreign.R.OrganisationStateSa = local
				break
			}
		}
	}

	return nil
}

// SetOrganisation of the organisationStateSa to the related item.
// Sets o.R.Organisation to related.
// Adds o to related.R.OrganisationStateSa.
func (o *OrganisationStateSa) SetOrganisation(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Organisation) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"organisation_state_sa\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"organisation_id"}),
		strmangle.WhereClause("\"", "\"", 2, organisationStateSaPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.OrganisationID = related.ID
	if o.R == nil {
		o.R = &organisationStateSaR{
			Organisation: related,
		}
	} else {
		o.R.Organisation = related
	}

	if related.R == nil {
		related.R = &organisationR{
			OrganisationStateSa: o,
		}
	} else {
		related.R.OrganisationStateSa = o
	}

	return nil
}

// OrganisationStateSas retrieves all the records using an executor.
func OrganisationStateSas(mods ...qm.QueryMod) organisationStateSaQuery {
	mods = append(mods, qm.From("\"organisation_state_sa\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"organisation_state_sa\".*"})
	}

	return organisationStateSaQuery{q}
}

// FindOrganisationStateSa retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindOrganisationStateSa(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*OrganisationStateSa, error) {
	organisationStateSaObj := &OrganisationStateSa{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"organisation_state_sa\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, organisationStateSaObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from organisation_state_sa")
	}

	if err = organisationStateSaObj.doAfterSelectHooks(ctx, exec); err != nil {
		return organisationStateSaObj, err
	}

	return organisationStateSaObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *OrganisationStateSa) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no organisation_state_sa provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(organisationStateSaColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	organisationStateSaInsertCacheMut.RLock()
	cache, cached := organisationStateSaInsertCache[key]
	organisationStateSaInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			organisationStateSaAllColumns,
			organisationStateSaColumnsWithDefault,
			organisationStateSaColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(organisationStateSaType, organisationStateSaMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(organisationStateSaType, organisationStateSaMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"organisation_state_sa\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"organisation_state_sa\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into organisation_state_sa")
	}

	if !cached {
		organisationStateSaInsertCacheMut.Lock()
		organisationStateSaInsertCache[key] = cache
		organisationStateSaInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the OrganisationStateSa.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *OrganisationStateSa) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	organisationStateSaUpdateCacheMut.RLock()
	cache, cached := organisationStateSaUpdateCache[key]
	organisationStateSaUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			organisationStateSaAllColumns,
			organisationStateSaPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update organisation_state_sa, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"organisation_state_sa\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, organisationStateSaPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(organisationStateSaType, organisationStateSaMapping, append(wl, organisationStateSaPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update organisation_state_sa row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for organisation_state_sa")
	}

	if !cached {
		organisationStateSaUpdateCacheMut.Lock()
		organisationStateSaUpdateCache[key] = cache
		organisationStateSaUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q organisationStateSaQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for organisation_state_sa")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for organisation_state_sa")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o OrganisationStateSaSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organisationStateSaPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"organisation_state_sa\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, organisationStateSaPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in organisationStateSa slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all organisationStateSa")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *OrganisationStateSa) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no organisation_state_sa provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(organisationStateSaColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	organisationStateSaUpsertCacheMut.RLock()
	cache, cached := organisationStateSaUpsertCache[key]
	organisationStateSaUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			organisationStateSaAllColumns,
			organisationStateSaColumnsWithDefault,
			organisationStateSaColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			organisationStateSaAllColumns,
			organisationStateSaPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert organisation_state_sa, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(organisationStateSaPrimaryKeyColumns))
			copy(conflict, organisationStateSaPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"organisation_state_sa\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(organisationStateSaType, organisationStateSaMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(organisationStateSaType, organisationStateSaMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert organisation_state_sa")
	}

	if !cached {
		organisationStateSaUpsertCacheMut.Lock()
		organisationStateSaUpsertCache[key] = cache
		organisationStateSaUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single OrganisationStateSa record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *OrganisationStateSa) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no OrganisationStateSa provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), organisationStateSaPrimaryKeyMapping)
	sql := "DELETE FROM \"organisation_state_sa\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from organisation_state_sa")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for organisation_state_sa")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q organisationStateSaQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no organisationStateSaQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from organisation_state_sa")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for organisation_state_sa")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o OrganisationStateSaSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(organisationStateSaBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organisationStateSaPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"organisation_state_sa\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, organisationStateSaPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from organisationStateSa slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for organisation_state_sa")
	}

	if len(organisationStateSaAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *OrganisationStateSa) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindOrganisationStateSa(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OrganisationStateSaSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := OrganisationStateSaSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organisationStateSaPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"organisation_state_sa\".* FROM \"organisation_state_sa\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, organisationStateSaPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in OrganisationStateSaSlice")
	}

	*o = slice

	return nil
}

// OrganisationStateSaExists checks if the OrganisationStateSa row exists.
func OrganisationStateSaExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"organisation_state_sa\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if organisation_state_sa exists")
	}

	return exists, nil
}

// Exists checks if the OrganisationStateSa row exists.
func (o *OrganisationStateSa) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return OrganisationStateSaExists(ctx, exec, o.ID)
}
