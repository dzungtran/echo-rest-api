package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"{{ .RootPackage }}/infrastructure/datastore"
	"{{ .RootPackage }}/modules/{{ .PluralName | ToKebab }}/domains"
	"{{ .RootPackage }}/pkg/constants"
	"{{ .RootPackage }}/pkg/contexts"
	sqlTools "{{ .RootPackage }}/pkg/sql-tools"
	"{{ .RootPackage }}/pkg/utils"
	"github.com/jmoiron/sqlx"
)

const (
	{{ .PluralName | ToLowerCamel }}TableName = "{{ .PluralName | ToSnake }}"
)

type {{ .SingularName }}Repository interface {
	Create(ctx context.Context, {{ .SingularName | ToLowerCamel }} *domains.{{ .SingularName }}) (int64, error)
	GetByID(ctx context.Context, id int64) (*domains.{{ .SingularName }}, error)
	Fetch(ctx context.Context, p ParamsForFetch{{ .PluralName }}) ([]*domains.{{ .SingularName }}, int64, error)
	Update(ctx context.Context, {{ .SingularName | ToLowerCamel }} *domains.{{ .SingularName }}, fieldsToUpdate []string) error
	DeleteById(ctx context.Context, id int64) error
}

type (
	pgsql{{ .SingularName }}Repository struct {
		db  *sqlx.DB
		sdb *sqlx.DB
	}
	ParamsForFetch{{ .PluralName }} struct {
		contexts.CommonParamsForFetch
	}
)

// NewPgsql{{ .SingularName }}Repository will create new an {{ .SingularName | ToLowerCamel }}Repository object representation of {{ .SingularName }}Repository interface
func NewPgsql{{ .SingularName }}Repository(mdbi *datastore.MasterDbInstance, sdbi *datastore.SlaveDbInstance) {{ .SingularName }}Repository {
	return &pgsql{{ .SingularName }}Repository{
		db:  mdbi.DBX(),
		sdb: sdbi.DBX(),
	}
}

func (r *pgsql{{ .SingularName }}Repository) Create(ctx context.Context, {{ .SingularName | ToLowerCamel }} *domains.{{ .SingularName }}) (newId int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	cols, vals := sqlTools.GetColumnsAndValuesFromStruct(
		ctx,
		{{ .SingularName | ToLowerCamel }},
		sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
		sqlTools.WithMapValuesAutoDateTimeFields([]string{"created_at", "updated_at"}),
	)

	query := psql.Insert({{ .PluralName | ToLowerCamel }}TableName).
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

func (r *pgsql{{ .SingularName }}Repository) GetByID(ctx context.Context, id int64) ({{ .SingularName | ToLowerCamel }} *domains.{{ .SingularName }}, err error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &domains.{{ .SingularName }}{})
	query, args, err := psql.Select(cols...).From({{ .PluralName | ToLowerCamel }}TableName).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return
	}

	{{ .SingularName | ToLowerCamel }} = &domains.{{ .SingularName }}{}
	err = r.sdb.Get({{ .SingularName | ToLowerCamel }}, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return
}

func (r *pgsql{{ .SingularName }}Repository) Fetch(ctx context.Context, params ParamsForFetch{{ .PluralName }}) ({{ .PluralName | ToLowerCamel }} []*domains.{{ .SingularName }}, count int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	type {{ .SingularName | ToLowerCamel }}WithCount struct {
		domains.{{ .SingularName }}
		Count int64 `db:"_count"` // special field for count
	}

	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &{{ .SingularName | ToLowerCamel }}WithCount{})
	query := psql.Select(sqlTools.ParseColumnsForSelect(cols)...).From({{ .PluralName | ToLowerCamel }}TableName)
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

	{{ .PluralName | ToLowerCamel }} = make([]*domains.{{ .SingularName }}, 0)
	for rows.Next() {
		var uwc {{ .SingularName | ToLowerCamel }}WithCount
		err = rows.StructScan(&uwc)
		if err != nil {
			return nil, count, err
		}

		count = uwc.Count
		u := uwc.{{ .SingularName }}
		{{ .PluralName | ToLowerCamel }} = append({{ .PluralName | ToLowerCamel }}, &u)
	}

	return
}

func (r *pgsql{{ .SingularName }}Repository) Update(ctx context.Context, {{ .SingularName | ToLowerCamel }} *domains.{{ .SingularName }}, fieldsToUpdate []string) (err error) {
	if len(fieldsToUpdate) == 0 {
		fieldsToUpdate = make([]string, 0)
	}

	if {{ .SingularName | ToLowerCamel }}.Id <= 0 {
		return errors.New("missing {{ .SingularName | ToLowerCamel }} id")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	query := psql.Update({{ .PluralName | ToLowerCamel }}TableName).
		SetMap(sqlTools.GetMapValuesFromStruct(
			ctx, {{ .SingularName | ToLowerCamel }},
			sqlTools.WithMapValuesSelectFields(fieldsToUpdate),
			sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
			sqlTools.WithMapValuesAutoDateTimeFields([]string{"updated_at"}),
		)).
		Where(squirrel.Eq{
			"id": {{ .SingularName | ToLowerCamel }}.Id,
		})

	affect, err := query.ExecContext(ctx)
	if err != nil {
		return
	}

	_, err = affect.RowsAffected()
	return
}

func (r *pgsql{{ .SingularName }}Repository) DeleteById(ctx context.Context, id int64) (err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	query := psql.Delete({{ .PluralName | ToLowerCamel }}TableName).Where(squirrel.Eq{
		"id": id,
	})

	_, err = query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return
}
