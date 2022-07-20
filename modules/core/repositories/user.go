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
	usersTableName = "users"
)

type UserRepository interface {
	Create(ctx context.Context, user *domains.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*domains.User, error)
	Fetch(ctx context.Context, p ParamsForFetchUsers) ([]*domains.User, int64, error)
	Update(ctx context.Context, user *domains.User, fieldsToUpdate []string) error
	DeleteById(ctx context.Context, id int64) error
	GetByCode(ctx context.Context, code string) (*domains.User, error)
	GetByEmail(ctx context.Context, email string) (*domains.User, error)
}

type (
	pgsqlUserRepository struct {
		db  *sqlx.DB
		sdb *sqlx.DB
	}
	ParamsForFetchUsers struct {
		contexts.CommonParamsForFetch
	}
)

// NewPgsqlUserRepository will create new an userRepository object representation of UserRepository interface
func NewPgsqlUserRepository(mdbi *datastore.MasterDbInstance, sdbi *datastore.SlaveDbInstance) UserRepository {
	return &pgsqlUserRepository{
		db:  mdbi.DBX(),
		sdb: sdbi.DBX(),
	}
}

func (r *pgsqlUserRepository) Create(ctx context.Context, user *domains.User) (newId int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	cols, vals := sqlTools.GetColumnsAndValuesFromStruct(
		ctx,
		user,
		sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
		sqlTools.WithMapValuesAutoDateTimeFields([]string{"created_at", "updated_at"}),
	)

	query := psql.Insert(usersTableName).
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

func (r *pgsqlUserRepository) GetByID(ctx context.Context, id int64) (user *domains.User, err error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &domains.User{})
	query, args, err := psql.Select(cols...).From(usersTableName).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return
	}

	user = &domains.User{}
	err = r.sdb.Get(user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return
}

func (r *pgsqlUserRepository) Fetch(ctx context.Context, params ParamsForFetchUsers) (users []*domains.User, count int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	type userWithCount struct {
		domains.User
		Count int64 `db:"_count"` // special field for count
	}

	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &userWithCount{})
	query := psql.Select(sqlTools.ParseColumnsForSelect(cols)...).From(usersTableName)
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

	users = make([]*domains.User, 0)
	for rows.Next() {
		var uwc userWithCount
		err = rows.StructScan(&uwc)
		if err != nil {
			return nil, count, err
		}

		count = uwc.Count
		u := uwc.User
		users = append(users, &u)
	}

	return
}

func (r *pgsqlUserRepository) Update(ctx context.Context, user *domains.User, fieldsToUpdate []string) (err error) {
	if len(fieldsToUpdate) == 0 {
		fieldsToUpdate = make([]string, 0)
	}

	if user.Id <= 0 {
		return errors.New("missing user id")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	query := psql.Update(usersTableName).
		SetMap(sqlTools.GetMapValuesFromStruct(
			ctx, user,
			sqlTools.WithMapValuesSelectFields(fieldsToUpdate),
			sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
			sqlTools.WithMapValuesAutoDateTimeFields([]string{"updated_at"}),
		)).
		Where(squirrel.Eq{
			"id": user.Id,
		})

	affect, err := query.ExecContext(ctx)
	if err != nil {
		return
	}

	_, err = affect.RowsAffected()
	return
}

func (r *pgsqlUserRepository) DeleteById(ctx context.Context, id int64) (err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	query := psql.Delete(usersTableName).Where(squirrel.Eq{
		"id": id,
	})

	_, err = query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return
}

func (r *pgsqlUserRepository) GetByCode(ctx context.Context, code string) (user *domains.User, err error) {
	if len(code) == 0 {
		return nil, errors.New("invalid code")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &domains.User{})
	query, args, err := psql.Select(cols...).From(usersTableName).
		Where(squirrel.Eq{
			"code": code,
		}).ToSql()
	if err != nil {
		return nil, err
	}

	user = &domains.User{}
	err = r.sdb.Get(user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return
}

func (r *pgsqlUserRepository) GetByEmail(ctx context.Context, email string) (user *domains.User, err error) {
	if len(email) == 0 {
		return nil, errors.New("invalid email")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &domains.User{})
	query, args, err := psql.Select(cols...).From(usersTableName).
		Where(squirrel.Eq{
			"email": email,
		}).ToSql()
	if err != nil {
		return nil, err
	}

	user = &domains.User{}
	err = r.sdb.Get(user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}
	return
}
