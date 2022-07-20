package sqlTools

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type (
	CommonRepository interface {
		Close()
	}
)

func NewPSQLStatementBuilder(db *sqlx.DB) squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db.DB)
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Log().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
