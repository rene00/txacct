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

// OrganisationStateTasmanium is an object representing the database table.
type OrganisationStateTasmanium struct {
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

	R *organisationStateTasmaniumR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L organisationStateTasmaniumL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var OrganisationStateTasmaniumColumns = struct {
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

var OrganisationStateTasmaniumTableColumns = struct {
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
	ID:                "organisation_state_tasmania.id",
	OrganisationID:    "organisation_state_tasmania.organisation_id",
	Name:              "organisation_state_tasmania.name",
	Abn:               "organisation_state_tasmania.abn",
	Address:           "organisation_state_tasmania.address",
	RecordDefunctRisk: "organisation_state_tasmania.record_defunct_risk",
	Region:            "organisation_state_tasmania.region",
	Phone:             "organisation_state_tasmania.phone",
	Mobile:            "organisation_state_tasmania.mobile",
	Freecall:          "organisation_state_tasmania.freecall",
	Fax:               "organisation_state_tasmania.fax",
}

// Generated where

var OrganisationStateTasmaniumWhere = struct {
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
	ID:                whereHelperint{field: "\"organisation_state_tasmania\".\"id\""},
	OrganisationID:    whereHelperint{field: "\"organisation_state_tasmania\".\"organisation_id\""},
	Name:              whereHelperstring{field: "\"organisation_state_tasmania\".\"name\""},
	Abn:               whereHelpernull_String{field: "\"organisation_state_tasmania\".\"abn\""},
	Address:           whereHelpernull_String{field: "\"organisation_state_tasmania\".\"address\""},
	RecordDefunctRisk: whereHelpernull_String{field: "\"organisation_state_tasmania\".\"record_defunct_risk\""},
	Region:            whereHelpernull_String{field: "\"organisation_state_tasmania\".\"region\""},
	Phone:             whereHelpernull_String{field: "\"organisation_state_tasmania\".\"phone\""},
	Mobile:            whereHelpernull_String{field: "\"organisation_state_tasmania\".\"mobile\""},
	Freecall:          whereHelpernull_String{field: "\"organisation_state_tasmania\".\"freecall\""},
	Fax:               whereHelpernull_String{field: "\"organisation_state_tasmania\".\"fax\""},
}

// OrganisationStateTasmaniumRels is where relationship names are stored.
var OrganisationStateTasmaniumRels = struct {
	Organisation string
}{
	Organisation: "Organisation",
}

// organisationStateTasmaniumR is where relationships are stored.
type organisationStateTasmaniumR struct {
	Organisation *Organisation `boil:"Organisation" json:"Organisation" toml:"Organisation" yaml:"Organisation"`
}

// NewStruct creates a new relationship struct
func (*organisationStateTasmaniumR) NewStruct() *organisationStateTasmaniumR {
	return &organisationStateTasmaniumR{}
}

func (r *organisationStateTasmaniumR) GetOrganisation() *Organisation {
	if r == nil {
		return nil
	}
	return r.Organisation
}

// organisationStateTasmaniumL is where Load methods for each relationship are stored.
type organisationStateTasmaniumL struct{}

var (
	organisationStateTasmaniumAllColumns            = []string{"id", "organisation_id", "name", "abn", "address", "record_defunct_risk", "region", "phone", "mobile", "freecall", "fax"}
	organisationStateTasmaniumColumnsWithoutDefault = []string{"organisation_id", "name"}
	organisationStateTasmaniumColumnsWithDefault    = []string{"id", "abn", "address", "record_defunct_risk", "region", "phone", "mobile", "freecall", "fax"}
	organisationStateTasmaniumPrimaryKeyColumns     = []string{"id"}
	organisationStateTasmaniumGeneratedColumns      = []string{}
)

type (
	// OrganisationStateTasmaniumSlice is an alias for a slice of pointers to OrganisationStateTasmanium.
	// This should almost always be used instead of []OrganisationStateTasmanium.
	OrganisationStateTasmaniumSlice []*OrganisationStateTasmanium
	// OrganisationStateTasmaniumHook is the signature for custom OrganisationStateTasmanium hook methods
	OrganisationStateTasmaniumHook func(context.Context, boil.ContextExecutor, *OrganisationStateTasmanium) error

	organisationStateTasmaniumQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	organisationStateTasmaniumType                 = reflect.TypeOf(&OrganisationStateTasmanium{})
	organisationStateTasmaniumMapping              = queries.MakeStructMapping(organisationStateTasmaniumType)
	organisationStateTasmaniumPrimaryKeyMapping, _ = queries.BindMapping(organisationStateTasmaniumType, organisationStateTasmaniumMapping, organisationStateTasmaniumPrimaryKeyColumns)
	organisationStateTasmaniumInsertCacheMut       sync.RWMutex
	organisationStateTasmaniumInsertCache          = make(map[string]insertCache)
	organisationStateTasmaniumUpdateCacheMut       sync.RWMutex
	organisationStateTasmaniumUpdateCache          = make(map[string]updateCache)
	organisationStateTasmaniumUpsertCacheMut       sync.RWMutex
	organisationStateTasmaniumUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var organisationStateTasmaniumAfterSelectHooks []OrganisationStateTasmaniumHook

var organisationStateTasmaniumBeforeInsertHooks []OrganisationStateTasmaniumHook
var organisationStateTasmaniumAfterInsertHooks []OrganisationStateTasmaniumHook

var organisationStateTasmaniumBeforeUpdateHooks []OrganisationStateTasmaniumHook
var organisationStateTasmaniumAfterUpdateHooks []OrganisationStateTasmaniumHook

var organisationStateTasmaniumBeforeDeleteHooks []OrganisationStateTasmaniumHook
var organisationStateTasmaniumAfterDeleteHooks []OrganisationStateTasmaniumHook

var organisationStateTasmaniumBeforeUpsertHooks []OrganisationStateTasmaniumHook
var organisationStateTasmaniumAfterUpsertHooks []OrganisationStateTasmaniumHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *OrganisationStateTasmanium) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *OrganisationStateTasmanium) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *OrganisationStateTasmanium) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *OrganisationStateTasmanium) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *OrganisationStateTasmanium) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *OrganisationStateTasmanium) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *OrganisationStateTasmanium) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *OrganisationStateTasmanium) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *OrganisationStateTasmanium) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range organisationStateTasmaniumAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddOrganisationStateTasmaniumHook registers your hook function for all future operations.
func AddOrganisationStateTasmaniumHook(hookPoint boil.HookPoint, organisationStateTasmaniumHook OrganisationStateTasmaniumHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		organisationStateTasmaniumAfterSelectHooks = append(organisationStateTasmaniumAfterSelectHooks, organisationStateTasmaniumHook)
	case boil.BeforeInsertHook:
		organisationStateTasmaniumBeforeInsertHooks = append(organisationStateTasmaniumBeforeInsertHooks, organisationStateTasmaniumHook)
	case boil.AfterInsertHook:
		organisationStateTasmaniumAfterInsertHooks = append(organisationStateTasmaniumAfterInsertHooks, organisationStateTasmaniumHook)
	case boil.BeforeUpdateHook:
		organisationStateTasmaniumBeforeUpdateHooks = append(organisationStateTasmaniumBeforeUpdateHooks, organisationStateTasmaniumHook)
	case boil.AfterUpdateHook:
		organisationStateTasmaniumAfterUpdateHooks = append(organisationStateTasmaniumAfterUpdateHooks, organisationStateTasmaniumHook)
	case boil.BeforeDeleteHook:
		organisationStateTasmaniumBeforeDeleteHooks = append(organisationStateTasmaniumBeforeDeleteHooks, organisationStateTasmaniumHook)
	case boil.AfterDeleteHook:
		organisationStateTasmaniumAfterDeleteHooks = append(organisationStateTasmaniumAfterDeleteHooks, organisationStateTasmaniumHook)
	case boil.BeforeUpsertHook:
		organisationStateTasmaniumBeforeUpsertHooks = append(organisationStateTasmaniumBeforeUpsertHooks, organisationStateTasmaniumHook)
	case boil.AfterUpsertHook:
		organisationStateTasmaniumAfterUpsertHooks = append(organisationStateTasmaniumAfterUpsertHooks, organisationStateTasmaniumHook)
	}
}

// One returns a single organisationStateTasmanium record from the query.
func (q organisationStateTasmaniumQuery) One(ctx context.Context, exec boil.ContextExecutor) (*OrganisationStateTasmanium, error) {
	o := &OrganisationStateTasmanium{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for organisation_state_tasmania")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all OrganisationStateTasmanium records from the query.
func (q organisationStateTasmaniumQuery) All(ctx context.Context, exec boil.ContextExecutor) (OrganisationStateTasmaniumSlice, error) {
	var o []*OrganisationStateTasmanium

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to OrganisationStateTasmanium slice")
	}

	if len(organisationStateTasmaniumAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all OrganisationStateTasmanium records in the query.
func (q organisationStateTasmaniumQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count organisation_state_tasmania rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q organisationStateTasmaniumQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if organisation_state_tasmania exists")
	}

	return count > 0, nil
}

// Organisation pointed to by the foreign key.
func (o *OrganisationStateTasmanium) Organisation(mods ...qm.QueryMod) organisationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.OrganisationID),
	}

	queryMods = append(queryMods, mods...)

	return Organisations(queryMods...)
}

// LoadOrganisation allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (organisationStateTasmaniumL) LoadOrganisation(ctx context.Context, e boil.ContextExecutor, singular bool, maybeOrganisationStateTasmanium interface{}, mods queries.Applicator) error {
	var slice []*OrganisationStateTasmanium
	var object *OrganisationStateTasmanium

	if singular {
		var ok bool
		object, ok = maybeOrganisationStateTasmanium.(*OrganisationStateTasmanium)
		if !ok {
			object = new(OrganisationStateTasmanium)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeOrganisationStateTasmanium)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeOrganisationStateTasmanium))
			}
		}
	} else {
		s, ok := maybeOrganisationStateTasmanium.(*[]*OrganisationStateTasmanium)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeOrganisationStateTasmanium)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeOrganisationStateTasmanium))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &organisationStateTasmaniumR{}
		}
		args = append(args, object.OrganisationID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &organisationStateTasmaniumR{}
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
		foreign.R.OrganisationStateTasmanium = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.OrganisationID == foreign.ID {
				local.R.Organisation = foreign
				if foreign.R == nil {
					foreign.R = &organisationR{}
				}
				foreign.R.OrganisationStateTasmanium = local
				break
			}
		}
	}

	return nil
}

// SetOrganisation of the organisationStateTasmanium to the related item.
// Sets o.R.Organisation to related.
// Adds o to related.R.OrganisationStateTasmanium.
func (o *OrganisationStateTasmanium) SetOrganisation(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Organisation) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"organisation_state_tasmania\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"organisation_id"}),
		strmangle.WhereClause("\"", "\"", 2, organisationStateTasmaniumPrimaryKeyColumns),
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
		o.R = &organisationStateTasmaniumR{
			Organisation: related,
		}
	} else {
		o.R.Organisation = related
	}

	if related.R == nil {
		related.R = &organisationR{
			OrganisationStateTasmanium: o,
		}
	} else {
		related.R.OrganisationStateTasmanium = o
	}

	return nil
}

// OrganisationStateTasmania retrieves all the records using an executor.
func OrganisationStateTasmania(mods ...qm.QueryMod) organisationStateTasmaniumQuery {
	mods = append(mods, qm.From("\"organisation_state_tasmania\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"organisation_state_tasmania\".*"})
	}

	return organisationStateTasmaniumQuery{q}
}

// FindOrganisationStateTasmanium retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindOrganisationStateTasmanium(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*OrganisationStateTasmanium, error) {
	organisationStateTasmaniumObj := &OrganisationStateTasmanium{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"organisation_state_tasmania\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, organisationStateTasmaniumObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from organisation_state_tasmania")
	}

	if err = organisationStateTasmaniumObj.doAfterSelectHooks(ctx, exec); err != nil {
		return organisationStateTasmaniumObj, err
	}

	return organisationStateTasmaniumObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *OrganisationStateTasmanium) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no organisation_state_tasmania provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(organisationStateTasmaniumColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	organisationStateTasmaniumInsertCacheMut.RLock()
	cache, cached := organisationStateTasmaniumInsertCache[key]
	organisationStateTasmaniumInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			organisationStateTasmaniumAllColumns,
			organisationStateTasmaniumColumnsWithDefault,
			organisationStateTasmaniumColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(organisationStateTasmaniumType, organisationStateTasmaniumMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(organisationStateTasmaniumType, organisationStateTasmaniumMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"organisation_state_tasmania\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"organisation_state_tasmania\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into organisation_state_tasmania")
	}

	if !cached {
		organisationStateTasmaniumInsertCacheMut.Lock()
		organisationStateTasmaniumInsertCache[key] = cache
		organisationStateTasmaniumInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the OrganisationStateTasmanium.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *OrganisationStateTasmanium) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	organisationStateTasmaniumUpdateCacheMut.RLock()
	cache, cached := organisationStateTasmaniumUpdateCache[key]
	organisationStateTasmaniumUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			organisationStateTasmaniumAllColumns,
			organisationStateTasmaniumPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update organisation_state_tasmania, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"organisation_state_tasmania\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, organisationStateTasmaniumPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(organisationStateTasmaniumType, organisationStateTasmaniumMapping, append(wl, organisationStateTasmaniumPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update organisation_state_tasmania row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for organisation_state_tasmania")
	}

	if !cached {
		organisationStateTasmaniumUpdateCacheMut.Lock()
		organisationStateTasmaniumUpdateCache[key] = cache
		organisationStateTasmaniumUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q organisationStateTasmaniumQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for organisation_state_tasmania")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for organisation_state_tasmania")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o OrganisationStateTasmaniumSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organisationStateTasmaniumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"organisation_state_tasmania\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, organisationStateTasmaniumPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in organisationStateTasmanium slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all organisationStateTasmanium")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *OrganisationStateTasmanium) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no organisation_state_tasmania provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(organisationStateTasmaniumColumnsWithDefault, o)

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

	organisationStateTasmaniumUpsertCacheMut.RLock()
	cache, cached := organisationStateTasmaniumUpsertCache[key]
	organisationStateTasmaniumUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			organisationStateTasmaniumAllColumns,
			organisationStateTasmaniumColumnsWithDefault,
			organisationStateTasmaniumColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			organisationStateTasmaniumAllColumns,
			organisationStateTasmaniumPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert organisation_state_tasmania, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(organisationStateTasmaniumPrimaryKeyColumns))
			copy(conflict, organisationStateTasmaniumPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"organisation_state_tasmania\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(organisationStateTasmaniumType, organisationStateTasmaniumMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(organisationStateTasmaniumType, organisationStateTasmaniumMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert organisation_state_tasmania")
	}

	if !cached {
		organisationStateTasmaniumUpsertCacheMut.Lock()
		organisationStateTasmaniumUpsertCache[key] = cache
		organisationStateTasmaniumUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single OrganisationStateTasmanium record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *OrganisationStateTasmanium) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no OrganisationStateTasmanium provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), organisationStateTasmaniumPrimaryKeyMapping)
	sql := "DELETE FROM \"organisation_state_tasmania\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from organisation_state_tasmania")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for organisation_state_tasmania")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q organisationStateTasmaniumQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no organisationStateTasmaniumQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from organisation_state_tasmania")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for organisation_state_tasmania")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o OrganisationStateTasmaniumSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(organisationStateTasmaniumBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organisationStateTasmaniumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"organisation_state_tasmania\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, organisationStateTasmaniumPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from organisationStateTasmanium slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for organisation_state_tasmania")
	}

	if len(organisationStateTasmaniumAfterDeleteHooks) != 0 {
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
func (o *OrganisationStateTasmanium) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindOrganisationStateTasmanium(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OrganisationStateTasmaniumSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := OrganisationStateTasmaniumSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), organisationStateTasmaniumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"organisation_state_tasmania\".* FROM \"organisation_state_tasmania\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, organisationStateTasmaniumPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in OrganisationStateTasmaniumSlice")
	}

	*o = slice

	return nil
}

// OrganisationStateTasmaniumExists checks if the OrganisationStateTasmanium row exists.
func OrganisationStateTasmaniumExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"organisation_state_tasmania\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if organisation_state_tasmania exists")
	}

	return exists, nil
}

// Exists checks if the OrganisationStateTasmanium row exists.
func (o *OrganisationStateTasmanium) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return OrganisationStateTasmaniumExists(ctx, exec, o.ID)
}