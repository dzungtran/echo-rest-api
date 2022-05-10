package datastore

import (
	"time"

	"github.com/dzungtran/echo-rest-api/pkg/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type MasterDbInstance struct {
	dbx *sqlx.DB
}

func (m *MasterDbInstance) DBX() *sqlx.DB {
	return m.dbx
}

type SlaveDbInstance struct {
	dbx *sqlx.DB
}

func (m *SlaveDbInstance) DBX() *sqlx.DB {
	return m.dbx
}

func NewMasterDbInstance(databaseURL string) *MasterDbInstance {
	return &MasterDbInstance{
		dbx: NewSqlXInstance(databaseURL),
	}
}

func NewSlaveDbInstance(databaseURL string) *SlaveDbInstance {
	return &SlaveDbInstance{
		dbx: NewSqlXInstance(databaseURL),
	}
}

// NewDatabase will create new database instance
func NewSqlXInstance(databaseURL string) *sqlx.DB {
	dbx, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		logger.Log().Panic(err)
	}

	if err := dbx.Ping(); err != nil {
		logger.Log().Fatal(err)
	}

	dbx.SetConnMaxLifetime(time.Minute * 5)
	dbx.SetMaxIdleConns(0)
	dbx.SetMaxOpenConns(5)
	return dbx
}
