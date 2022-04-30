package postgres

// Target: repositories/postgres/{{ .ModuleName | ToSnake }}.go

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/dzungtran/echo-rest-api/delivery/defines"
	"github.com/dzungtran/echo-rest-api/domains"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	sql_tools "github.com/dzungtran/echo-rest-api/pkg/sql-tools"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
)

const (
	{{ .ModuleName | ToLowerCamel }}sTableName = "{{ .ModuleName | ToSnake }}s"
)

type {{ .ModuleName }}Repository interface {
	Create(ctx context.Context, {{ .ModuleName | ToLowerCamel }} *domains.{{ .ModuleName }}) (int64, error)
	GetByID(ctx context.Context, id int64) (*domains.{{ .ModuleName }}, error)
	Fetch(ctx context.Context, p ParamsForFetch{{ .ModuleName }}s) ([]*domains.{{ .ModuleName }}, int64, error)
	Update(ctx context.Context, {{ .ModuleName | ToLowerCamel }} *domains.{{ .ModuleName }}, fieldsToUpdate []string) error
	DeleteById(ctx context.Context, id int64) error
}

type (
	pgsql{{ .ModuleName }}Repository struct {
		db  *sqlx.DB
		sdb *sqlx.DB
	}
	ParamsForFetch{{ .ModuleName }}s struct {
		defines.CommonParamsForFetch
	}
)

// NewPgsql{{ .ModuleName }}Repository will create new an {{ .ModuleName | ToLowerCamel }}Repository object representation of {{ .ModuleName }}Repository interface
func NewPgsql{{ .ModuleName }}Repository(mdbi *datastore.MasterDbInstance, sdbi *datastore.SlaveDbInstance) {{ .ModuleName }}Repository {
	return &pgsql{{ .ModuleName }}Repository{
		db:  mdbi.DBX(),
		sdb: sdbi.DBX(),
	}
}

func (r *pgsql{{ .ModuleName }}Repository) Create(ctx context.Context, {{ .ModuleName | ToLowerCamel }} *domains.{{ .ModuleName }}) (newId int64, err error) {
	psql := NewPSQLStatementBuilder(r.db)
	cols, vals := sql_tools.GetColumnsAndValuesFromStruct(
		ctx,
		{{ .ModuleName | ToLowerCamel }},
		sql_tools.WithMapValuesIgnoreFields([]string{"id"}),
		sql_tools.WithMapValuesAutoDateTimeFields([]string{"created_at", "updated_at"}),
	)

	query := psql.Insert({{ .ModuleName | ToLowerCamel }}sTableName).
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

func (r *pgsql{{ .ModuleName }}Repository) GetByID(ctx context.Context, id int64) ({{ .ModuleName | ToLowerCamel }} *domains.{{ .ModuleName }}, err error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	psql := NewPSQLStatementBuilder(r.sdb)
	cols, _ := sql_tools.GetColumnsAndValuesFromStruct(ctx, &domains.{{ .ModuleName }}{})
	query, args, err := psql.Select(cols...).From({{ .ModuleName | ToLowerCamel }}sTableName).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return
	}

	{{ .ModuleName | ToLowerCamel }} = &domains.{{ .ModuleName }}{}
	err = r.sdb.Get({{ .ModuleName | ToLowerCamel }}, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return
}

func (r *pgsql{{ .ModuleName }}Repository) Fetch(ctx context.Context, params ParamsForFetch{{ .ModuleName }}s) ({{ .ModuleName | ToLowerCamel }}s []*domains.{{ .ModuleName }}, count int64, err error) {
	psql := NewPSQLStatementBuilder(r.sdb)
	type {{ .ModuleName | ToLowerCamel }}WithCount struct {
		domains.{{ .ModuleName }}
		Count int64 `db:"_count"` // special field for count
	}

	cols, _ := sql_tools.GetColumnsAndValuesFromStruct(ctx, &{{ .ModuleName | ToLowerCamel }}WithCount{})
	query := psql.Select(sql_tools.ParseColumnsForSelect(cols)...).From({{ .ModuleName | ToLowerCamel }}sTableName)
	sqlQuery, args, err := sql_tools.BindCommonParamsToSelectBuilder(query, params.CommonParamsForFetch).
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

	{{ .ModuleName | ToLowerCamel }}s = make([]*domains.{{ .ModuleName }}, 0)
	for rows.Next() {
		var uwc {{ .ModuleName | ToLowerCamel }}WithCount
		err = rows.StructScan(&uwc)
		if err != nil {
			return nil, count, err
		}

		count = uwc.Count
		u := uwc.{{ .ModuleName }}
		{{ .ModuleName | ToLowerCamel }}s = append({{ .ModuleName | ToLowerCamel }}s, &u)
	}

	return
}

func (r *pgsql{{ .ModuleName }}Repository) Update(ctx context.Context, {{ .ModuleName | ToLowerCamel }} *domains.{{ .ModuleName }}, fieldsToUpdate []string) (err error) {
	if len(fieldsToUpdate) == 0 {
		fieldsToUpdate = make([]string, 0)
	}

	if {{ .ModuleName | ToLowerCamel }}.Id <= 0 {
		return errors.New("missing {{ .ModuleName | ToLowerCamel }} id")
	}

	psql := NewPSQLStatementBuilder(r.db)
	query := psql.Update({{ .ModuleName | ToLowerCamel }}sTableName).
		SetMap(sql_tools.GetMapValuesFromStruct(
			ctx, {{ .ModuleName | ToLowerCamel }},
			sql_tools.WithMapValuesSelectFields(fieldsToUpdate),
			sql_tools.WithMapValuesIgnoreFields([]string{"id"}),
			sql_tools.WithMapValuesAutoDateTimeFields([]string{"updated_at"}),
		)).
		Where(squirrel.Eq{
			"id": {{ .ModuleName | ToLowerCamel }}.Id,
		})

	affect, err := query.ExecContext(ctx)
	if err != nil {
		return
	}

	_, err = affect.RowsAffected()
	return
}

func (r *pgsql{{ .ModuleName }}Repository) DeleteById(ctx context.Context, id int64) (err error) {
	psql := NewPSQLStatementBuilder(r.db)
	query := psql.Delete({{ .ModuleName | ToLowerCamel }}sTableName).Where(squirrel.Eq{
		"id": id,
	})

	_, err = query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return
}
