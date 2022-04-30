package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
)

func NewSqlxTransaction(mdbi *datastore.MasterDbInstance) *SqlxTransaction {
	return &SqlxTransaction{
		db: mdbi.DBX(),
	}
}

type SqlxTransaction struct {
	db *sqlx.DB
}

func (s *SqlxTransaction) Init() *sqlx.Tx {
	return s.db.MustBegin()
}
