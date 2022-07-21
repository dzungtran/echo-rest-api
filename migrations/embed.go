package migrations

import (
	"database/sql"
	"embed"

	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*
var fs embed.FS

type MLog struct {
	log logger.Logger
}

func (MLog) Verbose() bool { return false }
func (l MLog) Printf(format string, v ...interface{}) {
	l.log.Infof(format, v...)
}

func RunAutoMigrate(db *sql.DB) {
	d, err := iofs.New(fs, "sql")
	if err != nil {
		logger.Log().Fatalw("auto migration - init iofs", "err", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log().Fatalw("auto migration - init driver", "err", err.Error())
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		logger.Log().Fatalw("auto migration - new instance", "err", err.Error())
	}

	defer m.Close()
	m.Log = &MLog{log: logger.Log()}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Log().Fatalw("auto migration - run up", "err", err.Error())
	}
	dbversion, dirty, err := m.Version()
	if err != nil {
		logger.Log().Errorw("auto migration - error get db version", "err", err.Error())
	}

	logger.Log().Infof("Current DB version: %v, dirty %v", dbversion, dirty)
}
