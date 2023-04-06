package console

import (
	"database/sql"
	"errors"
	"os"

	"github.com/oarkflow/migration"

	"github.com/oarkflow/framework/database/support"
	"github.com/oarkflow/framework/facades"
)

func getMigrate(con ...string) (*migration.Migrate, error) {
	connection := facades.Config.GetString("database.default")
	if len(con) > 0 {
		connection = con[0]
	}
	driver := facades.Config.GetString("database.connections." + connection + ".driver")
	dir := "./database/migrations"
	os.MkdirAll(dir, os.ModePerm)
	files, _ := os.ReadDir(dir)
	if len(files) == 0 {
		os.Create(dir + "/.gitkeep")
	}
	switch driver {
	case support.Mysql:
		dsn := support.GetMysqlDsn(connection)
		if dsn == "" {
			return nil, nil
		}

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		return migration.New(migration.Config{
			DB:         db,
			IsEmbedded: false,
			Dir:        dir,
			TableName:  facades.Config.GetString("database.migrations"),
			Dialect:    "mysql",
		}), nil
	case support.Postgresql:
		dsn := support.GetPostgresqlDsn(connection)
		if dsn == "" {
			return nil, nil
		}

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}

		return migration.New(migration.Config{
			DB:         db,
			IsEmbedded: false,
			Dir:        dir,
			TableName:  facades.Config.GetString("database.migrations"),
			Dialect:    "postgresql",
		}), nil
	case support.Sqlite:
		dsn := support.GetSqliteDsn(connection)
		if dsn == "" {
			return nil, nil
		}

		db, err := sql.Open("sqlite3", dsn)
		if err != nil {
			return nil, err
		}

		return migration.New(migration.Config{
			DB:         db,
			IsEmbedded: false,
			Dir:        dir,
			TableName:  facades.Config.GetString("database.migrations"),
			Dialect:    "sqlite3",
		}), nil
	default:
		return nil, errors.New("database driver only support mysql, postgresql and sqlite")
	}
}
