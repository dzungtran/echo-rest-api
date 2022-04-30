package migrations

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
)

//go:embed sql/*
var fs embed.FS

func RunAutoMigrate(db *sql.DB) {
	d, err := iofs.New(fs, "sql")
	if err != nil {
		logger.Log().Fatalw("auto migration - init iofs", "err", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		logger.Log().Fatalw("auto migration - new instance", "err", err.Error())
	}

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
