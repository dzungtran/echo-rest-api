package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	sqlTools "github.com/dzungtran/echo-rest-api/pkg/sql-tools"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/jmoiron/sqlx"
)

const (
	orgsTableName = "orgs"
)

type OrgRepository interface {
	Create(ctx context.Context, org *domains.Org) (int64, error)
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, org *domains.Org) (int64, error)
	GetByID(ctx context.Context, id int64) (*domains.Org, error)
	Fetch(ctx context.Context, p ParamsForFetchOrgs) ([]*domains.Org, int64, error)
	Update(ctx context.Context, org *domains.Org, fieldsToUpdate []string) error
	DeleteById(ctx context.Context, id int64) error
}

type (
	pgsqlOrgRepository struct {
		db  *sqlx.DB
		sdb *sqlx.DB
	}
	ParamsForFetchOrgs struct {
		contexts.CommonParamsForFetch
		Ids      []int64
		Statuses []string
	}
)

// NewPgsqlOrgRepository will create new an orgRepository object representation of OrgRepository interface
func NewPgsqlOrgRepository(mdbi *datastore.MasterDbInstance, sdbi *datastore.SlaveDbInstance) OrgRepository {
	return &pgsqlOrgRepository{
		db:  mdbi.DBX(),
		sdb: sdbi.DBX(),
	}
}

func (r *pgsqlOrgRepository) Create(ctx context.Context, org *domains.Org) (newId int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	cols, vals := sqlTools.GetColumnsAndValuesFromStruct(
		ctx,
		org,
		sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
		sqlTools.WithMapValuesAutoDateTimeFields([]string{"created_at", "updated_at"}),
	)

	query := psql.Insert(orgsTableName).
		Columns(cols...).
		Values(vals...).
		Suffix(`RETURNING id`)

	err = query.QueryRowContext(ctx).Scan(&newId)
	if err != nil {
		if utils.IsDuplicatedError(err) {
			err = constants.ErrDuplicated
		}
		return
	}
	return
}

func (r *pgsqlOrgRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, org *domains.Org) (newId int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	cols, vals := sqlTools.GetColumnsAndValuesFromStruct(
		ctx,
		org,
		sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
		sqlTools.WithMapValuesAutoDateTimeFields([]string{"created_at", "updated_at"}),
	)

	query := psql.Insert(orgsTableName).
		Columns(cols...).
		Values(vals...).
		Suffix(`RETURNING id`)

	sqlQuery, agrs, err := query.ToSql()
	if err != nil {
		return
	}

	err = tx.QueryRowxContext(ctx, sqlQuery, agrs...).Scan(&newId)
	if err != nil {
		if utils.IsDuplicatedError(err) {
			err = constants.ErrDuplicated
		}
		return
	}
	return
}

func (r *pgsqlOrgRepository) GetByID(ctx context.Context, id int64) (org *domains.Org, err error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &domains.Org{})
	query, args, err := psql.Select(cols...).From(orgsTableName).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return
	}

	org = &domains.Org{}
	err = r.sdb.Get(org, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return
}

func (r *pgsqlOrgRepository) Fetch(ctx context.Context, params ParamsForFetchOrgs) (orgs []*domains.Org, count int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	type orgWithCount struct {
		domains.Org
		Count int64 `db:"_count"` // special field for count
	}

	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &orgWithCount{})
	query := psql.Select(sqlTools.ParseColumnsForSelect(cols)...).From(orgsTableName)
	query = r.buildQueryFilters(query, params)
	sqlQuery, args, err := sqlTools.BindCommonParamsToSelectBuilder(query, params.CommonParamsForFetch).
		OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, count, err
	}

	rows, err := r.sdb.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, count, constants.ErrNotFound
		}
		return nil, count, err
	}

	orgs = make([]*domains.Org, 0)
	for rows.Next() {
		var uwc orgWithCount
		err = rows.StructScan(&uwc)
		if err != nil {
			return nil, count, err
		}

		count = uwc.Count
		u := uwc.Org
		orgs = append(orgs, &u)
	}

	return
}

func (r *pgsqlOrgRepository) Update(ctx context.Context, org *domains.Org, fieldsToUpdate []string) (err error) {
	if len(fieldsToUpdate) == 0 {
		fieldsToUpdate = make([]string, 0)
	}

	if org.Id <= 0 {
		return errors.New("missing org id")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	query := psql.Update(orgsTableName).
		SetMap(sqlTools.GetMapValuesFromStruct(
			ctx, org,
			sqlTools.WithMapValuesSelectFields(fieldsToUpdate),
			sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
			sqlTools.WithMapValuesAutoDateTimeFields([]string{"updated_at"}),
		)).
		Where(squirrel.Eq{
			"id": org.Id,
		})

	affect, err := query.ExecContext(ctx)
	if err != nil {
		return
	}

	_, err = affect.RowsAffected()
	return
}

func (r *pgsqlOrgRepository) DeleteById(ctx context.Context, id int64) (err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	query := psql.Delete(orgsTableName).Where(squirrel.Eq{
		"id": id,
	})

	_, err = query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return
}

func (r *pgsqlOrgRepository) buildQueryFilters(builder squirrel.SelectBuilder, params ParamsForFetchOrgs) squirrel.SelectBuilder {
	if len(params.Ids) > 0 {
		builder = builder.Where(squirrel.Eq{
			"id": params.Ids,
		})
	}

	if len(params.Statuses) > 0 {
		builder = builder.Where(squirrel.Eq{
			"status": params.Statuses,
		})
	}
	return builder
}
