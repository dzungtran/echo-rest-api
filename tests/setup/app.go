package setup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dzungtran/echo-rest-api/cmd/api/di"
	"github.com/dzungtran/echo-rest-api/config"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/migrations"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
)

var (
	diContainer *dig.Container
	appConf     *config.AppConfig
)

func LoadTestEnv() {
	logger.Reload(
		// logger.WithConfigLevel("info"),
		logger.WithConfigEncoding("console"),
	)
	envFiles := []string{
		"../../.env.test",
		".env.test",
	}
	fPath := ""
	for _, fp := range envFiles {
		if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
			continue
		}
		fPath = filepath.FromSlash(fp)
	}

	logger.Log().Infof("Load env file: %s", fPath)
	err := godotenv.Load(fPath)
	if err != nil {
		logger.Log().Fatalf("Error loading .env file, details: %v", err)
	}
}

func TruncateTables() {
	tables := []string{
		"orgs",
		"schema_migrations",
		"users",
		"users_orgs",
	}

	var db *sqlx.DB

	di := GetTestDIContainer()
	if di == nil {
		logger.Log().Fatal("needs init test env first")
	}

	di.Invoke(func(mdbi *datastore.MasterDbInstance) {
		db = mdbi.DBX()
	})

	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS=0;")
	query := ""
	for _, v := range tables {
		query = query + fmt.Sprintf("TRUNCATE TABLE %s;", v)
	}
	_, _ = db.Exec(query)
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS=1;")

	logger.Log().Info("truncating table is finished")
}

func InitTest() {
	LoadTestEnv()
	appConf, _ = config.InitAppConfig()
	logger.Log().Debugw("Config", "data", appConf)
	mDBInstance := datastore.NewMasterDbInstance(appConf.DatabaseURL)
	sDBInstance := datastore.NewSlaveDbInstance(appConf.DatabaseURL)
	migrations.RunAutoMigrate(mDBInstance.DBX().DB)
	diContainer = di.BuildDIContainer(mDBInstance, sDBInstance, appConf)
}

func GetTestDIContainer() *dig.Container {
	return diContainer
}

func GetTestAppConfig() *config.AppConfig {
	return appConf
}
