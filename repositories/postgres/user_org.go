package postgres

// Target: repositories/postgres/user_org.go

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dzungtran/echo-rest-api/delivery/defines"
	"github.com/dzungtran/echo-rest-api/domains"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	sql_tools "github.com/dzungtran/echo-rest-api/pkg/sql-tools"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/jmoiron/sqlx"
)

const (
	userOrgsTableName = "users_orgs"
)

type UserOrgRepository interface {
	Create(ctx context.Context, userOrg *domains.UserOrg) (int64, error)
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, userOrg *domains.UserOrg) (int64, error)
	UpdateByUserIdAndOrgId(ctx context.Context, userOrg *domains.UserOrg, fieldsToUpdate []string) error
	DeleteByUserIdAndOrgId(ctx context.Context, userId, orgId int64) error
	Fetch(ctx context.Context, params ParamsForFetchUserOrgs) (rs []*domains.UserOrg, count int64, err error)
}

type (
	pgsqlUserOrgRepository struct {
		db  *sqlx.DB
		sdb *sqlx.DB
	}
	ParamsForFetchUserOrgs struct {
		UserIds []int64
		OrgId   int64
		Emails  []string
		defines.CommonParamsForFetch
	}
)

// NewPgsqlUserOrgRepository will create new an userOrgRepository object representation of UserOrgRepository interface
func NewPgsqlUserOrgRepository(mdbi *datastore.MasterDbInstance, sdbi *datastore.SlaveDbInstance) UserOrgRepository {
	return &pgsqlUserOrgRepository{
		db:  mdbi.DBX(),
		sdb: sdbi.DBX(),
	}
}

func (r *pgsqlUserOrgRepository) Create(ctx context.Context, userOrg *domains.UserOrg) (newId int64, err error) {
	psql := NewPSQLStatementBuilder(r.db)
	cols, vals := sql_tools.GetColumnsAndValuesFromStruct(
		ctx,
		userOrg,
		sql_tools.WithMapValuesIgnoreFields([]string{"id"}),
		sql_tools.WithMapValuesAutoDateTimeFields([]string{"created_at", "updated_at"}),
	)

	query := psql.Insert(userOrgsTableName).
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

func (r *pgsqlUserOrgRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, userOrg *domains.UserOrg) (newId int64, err error) {
	psql := NewPSQLStatementBuilder(r.db)
	cols, vals := sql_tools.GetColumnsAndValuesFromStruct(
		ctx,
		userOrg,
		sql_tools.WithMapValuesIgnoreFields([]string{"id"}),
		sql_tools.WithMapValuesAutoDateTimeFields([]string{"created_at", "updated_at"}),
	)

	query := psql.Insert(userOrgsTableName).
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

func (r *pgsqlUserOrgRepository) GetByID(ctx context.Context, id int64) (userOrg *domains.UserOrg, err error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	psql := NewPSQLStatementBuilder(r.sdb)
	cols, _ := sql_tools.GetColumnsAndValuesFromStruct(ctx, &domains.UserOrg{})
	query, args, err := psql.Select(cols...).From(userOrgsTableName).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return
	}

	userOrg = &domains.UserOrg{}
	err = r.sdb.Get(userOrg, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return
}

func (r *pgsqlUserOrgRepository) UpdateByUserIdAndOrgId(ctx context.Context, userOrg *domains.UserOrg, fieldsToUpdate []string) (err error) {
	if len(fieldsToUpdate) == 0 {
		fieldsToUpdate = make([]string, 0)
	}

	if userOrg.UserId <= 0 {
		return errors.New("missing user id")
	}

	if userOrg.OrgId <= 0 {
		return errors.New("missing org id")
	}

	psql := NewPSQLStatementBuilder(r.db)
	query := psql.Update(userOrgsTableName).
		SetMap(sql_tools.GetMapValuesFromStruct(
			ctx, userOrg,
			sql_tools.WithMapValuesSelectFields(fieldsToUpdate),
			sql_tools.WithMapValuesIgnoreFields([]string{"id"}),
			sql_tools.WithMapValuesAutoDateTimeFields([]string{"updated_at"}),
		)).
		Where(squirrel.Eq{
			"user_id": userOrg.UserId,
			"org_id":  userOrg.OrgId,
		})

	affect, err := query.ExecContext(ctx)
	if err != nil {
		return
	}

	_, err = affect.RowsAffected()
	return
}

func (r *pgsqlUserOrgRepository) DeleteByUserIdAndOrgId(ctx context.Context, uid, oid int64) (err error) {
	psql := NewPSQLStatementBuilder(r.db)
	query := psql.Delete(userOrgsTableName).Where(squirrel.Eq{
		"user_id": uid,
		"org_id":  oid,
	})

	_, err = query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return
}

func (r *pgsqlUserOrgRepository) Fetch(ctx context.Context, params ParamsForFetchUserOrgs) (rs []*domains.UserOrg, count int64, err error) {
	psql := NewPSQLStatementBuilder(r.sdb)
	type orgWithCount struct {
		domains.UserOrg
		Count int64 `db:"_count"` // special field for count
	}

	tAlias := "ou"
	cols, _ := sql_tools.GetColumnsAndValuesFromStruct(ctx, &orgWithCount{})
	query := psql.Select(sql_tools.ParseColumnsForSelectWithAlias(cols, tAlias)...).From(userOrgsTableName + " AS " + tAlias)
	query = r.buildQueryFilters(query, params, tAlias)
	sqlQuery, args, err := sql_tools.
		BindCommonParamsToSelectBuilder(query, params.CommonParamsForFetch).
		OrderBy(tAlias + ".created_at DESC").ToSql()
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

	rs = make([]*domains.UserOrg, 0)
	for rows.Next() {
		var uwc orgWithCount
		err = rows.StructScan(&uwc)
		if err != nil {
			return nil, count, err
		}

		count = uwc.Count
		u := uwc.UserOrg
		rs = append(rs, &u)
	}

	return
}

func (r *pgsqlUserOrgRepository) buildQueryFilters(builder squirrel.SelectBuilder, params ParamsForFetchUserOrgs, tAlias string) squirrel.SelectBuilder {
	if len(params.UserIds) > 0 {
		builder = builder.Where(squirrel.Eq{
			tAlias + ".user_id": params.UserIds,
		})
	}

	if params.OrgId > 0 {
		builder = builder.Where(squirrel.Eq{
			tAlias + ".org_id": params.OrgId,
		})
	}

	if len(params.Emails) > 0 {
		builder = builder.Join(fmt.Sprintf("%s AS u ON u.id = %s.user_id", usersTableName, tAlias))
		builder = builder.Where(squirrel.Eq{
			"u.email": params.Emails,
		})
	}

	return builder
}
