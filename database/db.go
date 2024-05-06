package database

import (
	"context"

	"github.com/oarkflow/log"

	sqlx "github.com/oarkflow/squealx"
	"github.com/pkg/errors"

	"github.com/oarkflow/framework/contracts/database"
	"github.com/oarkflow/framework/database/support"
	"github.com/oarkflow/framework/facades"
)

type DB struct {
	connection      string
	defaultInstance database.Sqlx
	instances       map[string]database.Sqlx
}

func (r *DB) Connection(name string) database.DB {
	defaultConnection := facades.Config.GetString("database.default")
	if name == "" {
		name = defaultConnection
	}

	r.connection = name
	if r.instances == nil {
		r.instances = make(map[string]database.Sqlx)
	}

	if _, exist := r.instances[name]; exist {
		return r
	}

	dsn, err := GetDsn(name)
	if err != nil {
		log.Error().Err(err).Msg("database dsn error: %v")
		return r
	}

	db, err := sqlx.Open("mysql", dsn, name)
	if err != nil {
		log.Error().Err(err).Msg("database open error: %v")
		return r
	}

	r.instances[name] = db

	if name == defaultConnection {
		r.defaultInstance = db
	}

	return r
}

func (r *DB) Query() database.Sqlx {
	if r.connection == "" {
		if r.defaultInstance == nil {
			r.Connection("")
		}

		return r.defaultInstance
	}

	instance, exist := r.instances[r.connection]
	if !exist {
		return nil
	}

	r.connection = ""

	return instance
}

func (r *DB) Transaction(ctx context.Context, txFunc func(tx *sqlx.Tx) error) error {
	tx, err := r.Query().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := txFunc(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrapf(err, "rollback error: %v", err)
		}

		return err
	} else {
		return tx.Commit()
	}
}

func GetDsn(connection string) (string, error) {
	driver := facades.Config.GetString("database.connections." + connection + ".driver")
	switch driver {
	case support.Mysql:
		return support.GetMysqlDsn(connection), nil
	case support.Postgresql:
		return support.GetPostgresqlDsn(connection), nil
	case support.Sqlite:
		return support.GetSqliteDsn(connection), nil
	case support.Sqlserver:
		return support.GetSqlserverDsn(connection), nil
	default:
		return "", errors.New("database driver only support mysql, postgresql, sqlite and sqlserver")
	}
}
