package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/modules/projects/domains"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	sqlTools "github.com/dzungtran/echo-rest-api/pkg/sql-tools"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/jmoiron/sqlx"
)

const (
	projectsTableName = "projects"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *domains.Project) (int64, error)
	GetByID(ctx context.Context, id int64) (*domains.Project, error)
	Fetch(ctx context.Context, p ParamsForFetchProjects) ([]*domains.Project, int64, error)
	Update(ctx context.Context, project *domains.Project, fieldsToUpdate []string) error
	DeleteById(ctx context.Context, id int64) error
}

type (
	pgsqlProjectRepository struct {
		db  *sqlx.DB
		sdb *sqlx.DB
	}
	ParamsForFetchProjects struct {
		contexts.CommonParamsForFetch
	}
)

// NewPgsqlProjectRepository will create new an projectRepository object representation of ProjectRepository interface
func NewPgsqlProjectRepository(mdbi *datastore.MasterDbInstance, sdbi *datastore.SlaveDbInstance) ProjectRepository {
	return &pgsqlProjectRepository{
		db:  mdbi.DBX(),
		sdb: sdbi.DBX(),
	}
}

func (r *pgsqlProjectRepository) Create(ctx context.Context, project *domains.Project) (newId int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	cols, vals := sqlTools.GetColumnsAndValuesFromStruct(
		ctx,
		project,
		sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
		sqlTools.WithMapValuesAutoDateTimeFields([]string{"created_at", "updated_at"}),
	)

	query := psql.Insert(projectsTableName).
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

func (r *pgsqlProjectRepository) GetByID(ctx context.Context, id int64) (project *domains.Project, err error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &domains.Project{})
	query, args, err := psql.Select(cols...).From(projectsTableName).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return
	}

	project = &domains.Project{}
	err = r.sdb.Get(project, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return
}

func (r *pgsqlProjectRepository) Fetch(ctx context.Context, params ParamsForFetchProjects) (projects []*domains.Project, count int64, err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.sdb)
	type projectWithCount struct {
		domains.Project
		Count int64 `db:"_count"` // special field for count
	}

	cols, _ := sqlTools.GetColumnsAndValuesFromStruct(ctx, &projectWithCount{})
	query := psql.Select(sqlTools.ParseColumnsForSelect(cols)...).From(projectsTableName)
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

	projects = make([]*domains.Project, 0)
	for rows.Next() {
		var uwc projectWithCount
		err = rows.StructScan(&uwc)
		if err != nil {
			return nil, count, err
		}

		count = uwc.Count
		u := uwc.Project
		projects = append(projects, &u)
	}

	return
}

func (r *pgsqlProjectRepository) Update(ctx context.Context, project *domains.Project, fieldsToUpdate []string) (err error) {
	if len(fieldsToUpdate) == 0 {
		fieldsToUpdate = make([]string, 0)
	}

	if project.Id <= 0 {
		return errors.New("missing project id")
	}

	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	query := psql.Update(projectsTableName).
		SetMap(sqlTools.GetMapValuesFromStruct(
			ctx, project,
			sqlTools.WithMapValuesSelectFields(fieldsToUpdate),
			sqlTools.WithMapValuesIgnoreFields([]string{"id"}),
			sqlTools.WithMapValuesAutoDateTimeFields([]string{"updated_at"}),
		)).
		Where(squirrel.Eq{
			"id": project.Id,
		})

	affect, err := query.ExecContext(ctx)
	if err != nil {
		return
	}

	_, err = affect.RowsAffected()
	return
}

func (r *pgsqlProjectRepository) DeleteById(ctx context.Context, id int64) (err error) {
	psql := sqlTools.NewPSQLStatementBuilder(r.db)
	query := psql.Delete(projectsTableName).Where(squirrel.Eq{
		"id": id,
	})

	_, err = query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return
}
