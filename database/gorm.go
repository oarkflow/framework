package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/oarkflow/framework/contracts/database/orm"
	"github.com/oarkflow/framework/database/support"
	"github.com/oarkflow/framework/facades"
)

type GormDB struct {
	orm.Query
	instance *gorm.DB
}

func NewGormDB(ctx context.Context, connection string, config *gorm.Config, disableLog bool) (orm.DB, error) {
	db, err := NewGormInstance(connection, config, disableLog)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, nil
	}

	if ctx != nil {
		db = db.WithContext(ctx)
	}

	if facades.Config.GetBool("database.enable_connection_pool") {
		maxIdleTime := facades.Config.GetInt("database.max_idle_time")
		maxLifeTime := facades.Config.GetInt("database.max_life_time")
		maxIdleConnections := facades.Config.GetInt("database.max_idle_connections")
		maxOpenConnections := facades.Config.GetInt("database.max_open_connections")
		if maxOpenConnections == 0 || maxOpenConnections > 100 {
			maxOpenConnections = 100
		}
		if maxIdleConnections == 0 {
			maxIdleConnections = int(float64(maxOpenConnections) - (20 / float64(maxOpenConnections) * 100))
		}
		resolver := dbresolver.Register(dbresolver.Config{})
		if maxIdleTime > 0 {
			resolver = resolver.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)
		}
		if maxLifeTime > 0 {
			resolver = resolver.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Hour)
		}
		err = db.Use(resolver.SetMaxIdleConns(maxIdleConnections).SetMaxOpenConns(maxOpenConnections))
		if err != nil {
			return nil, err
		}
	}
	return &GormDB{
		Query:    NewGormQuery(db),
		instance: db,
	}, nil
}

func NewGormInstance(connection string, config *gorm.Config, disableLog bool) (*gorm.DB, error) {
	var cfg *gorm.Config

	gormConfig, err := getGormConfig(connection)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("init gorm config error: %v", err))
	}
	if gormConfig == nil {
		return nil, nil
	}
	if config == nil {
		cfg = &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
		}
	} else {
		cfg = config
	}
	if disableLog {
		cfg.Logger = nil
	} else {
		if cfg.Logger == nil {
			logger := getLogger()
			cfg.Logger = logger
		}
	}
	return gorm.Open(gormConfig, cfg)
}

func getLogger() glog.Interface {
	var logLevel glog.LogLevel
	if facades.Config.GetBool("app.debug") {
		logLevel = glog.Info
	} else {
		logLevel = glog.Error
	}

	logger := New(log.New(os.Stdout, "\r\n", log.LstdFlags), glog.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  glog.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})
	logger.LogMode(logLevel)
	return logger
}

func (r *GormDB) Begin() (orm.Transaction, error) {
	tx := r.instance.Begin()

	return NewGormTransaction(tx), tx.Error
}

func (r *GormDB) Instance() *gorm.DB {
	return r.instance
}

type GormTransaction struct {
	orm.Query
	instance *gorm.DB
}

func NewGormTransaction(instance *gorm.DB) orm.Transaction {
	return &GormTransaction{Query: NewGormQuery(instance), instance: instance}
}

func (r *GormTransaction) Commit() *gorm.DB {
	return r.instance.Commit()
}

func (r *GormTransaction) Rollback() *gorm.DB {
	return r.instance.Rollback()
}

type GormQuery struct {
	instance *gorm.DB
}

func NewGormQuery(instance *gorm.DB) orm.Query {
	return &GormQuery{instance: instance}
}

func (r *GormQuery) Driver() orm.Driver {
	return orm.Driver(r.instance.Dialector.Name())
}

func (r *GormQuery) Count(count *int64) *gorm.DB {
	return r.instance.Count(count)
}

func (r *GormQuery) Create(value interface{}) *gorm.DB {
	return r.instance.Create(value)
}

func (r *GormQuery) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return r.instance.Delete(value, conds...)
}

func (r *GormQuery) Distinct(args ...interface{}) orm.Query {
	tx := r.instance.Distinct(args...)

	return NewGormQuery(tx)
}

func (r *GormQuery) Exec(sql string, values ...interface{}) *gorm.DB {
	return r.instance.Exec(sql, values...)
}

func (r *GormQuery) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return r.instance.Find(dest, conds...)
}

func (r *GormQuery) First(dest interface{}) *gorm.DB {
	return r.instance.First(dest)
}

func (r *GormQuery) FirstOrCreate(dest interface{}, conds ...interface{}) *gorm.DB {
	if len(conds) > 1 {
		return r.instance.Attrs([]interface{}{conds[1]}...).FirstOrCreate(dest, []interface{}{conds[0]}...)
	}
	return r.instance.FirstOrCreate(dest, conds...)
}

func (r *GormQuery) ForceDelete(value interface{}, conds ...interface{}) *gorm.DB {
	return r.instance.Unscoped().Delete(value, conds...)
}

func (r *GormQuery) Get(dest interface{}) *gorm.DB {
	return r.instance.Find(dest)
}

func (r *GormQuery) Group(name string) orm.Query {
	tx := r.instance.Group(name)

	return NewGormQuery(tx)
}

func (r *GormQuery) Having(query interface{}, args ...interface{}) orm.Query {
	tx := r.instance.Having(query, args...)

	return NewGormQuery(tx)
}

func (r *GormQuery) Join(query string, args ...interface{}) orm.Query {
	tx := r.instance.Joins(query, args...)

	return NewGormQuery(tx)
}

func (r *GormQuery) Limit(limit int) orm.Query {
	tx := r.instance.Limit(limit)

	return NewGormQuery(tx)
}

func (r *GormQuery) Model(value interface{}) orm.Query {
	tx := r.instance.Model(value)

	return NewGormQuery(tx)
}

func (r *GormQuery) Offset(offset int) orm.Query {
	tx := r.instance.Offset(offset)

	return NewGormQuery(tx)
}

func (r *GormQuery) Order(value interface{}) orm.Query {
	tx := r.instance.Order(value)

	return NewGormQuery(tx)
}

func (r *GormQuery) OrWhere(query interface{}, args ...interface{}) orm.Query {
	tx := r.instance.Or(query, args...)

	return NewGormQuery(tx)
}

func (r *GormQuery) Pluck(column string, dest interface{}) *gorm.DB {
	return r.instance.Pluck(column, dest)
}

func (r *GormQuery) Raw(sql string, values ...interface{}) orm.Query {
	tx := r.instance.Raw(sql, values...)

	return NewGormQuery(tx)
}

func (r *GormQuery) Save(value interface{}) *gorm.DB {
	return r.instance.Save(value)
}

func (r *GormQuery) Scan(dest interface{}) *gorm.DB {
	return r.instance.Scan(dest)
}

func (r *GormQuery) Select(query interface{}, args ...interface{}) orm.Query {
	tx := r.instance.Select(query, args...)

	return NewGormQuery(tx)
}

func (r *GormQuery) Table(name string, args ...interface{}) orm.Query {
	tx := r.instance.Table(name, args...)

	return NewGormQuery(tx)
}

func (r *GormQuery) Update(column string, value interface{}) *gorm.DB {
	return r.instance.Update(column, value)
}

func (r *GormQuery) Updates(values interface{}) *gorm.DB {
	return r.instance.Updates(values)
}

func (r *GormQuery) Where(query interface{}, args ...interface{}) orm.Query {
	tx := r.instance.Where(query, args...)

	return NewGormQuery(tx)
}

func (r *GormQuery) With(query string, args ...interface{}) orm.Query {
	tx := r.instance.Preload(query, args...)
	return NewGormQuery(tx)
}

func (r *GormQuery) WithTrashed() orm.Query {
	tx := r.instance.Unscoped()
	return NewGormQuery(tx)
}

func (r *GormQuery) Fields(schema, name string) (fields []string, err error) {
	var f []map[string]string
	switch r.Driver().String() {
	case "mysql", "mariadb":
		err = r.instance.Raw("SELECT column_name as name FROm information_schema.columns WHERE table_schema = ? AND table_name = ?", schema, name).Scan(&f).Error
		if err != nil {
			return
		}
		for _, field := range f {
			fields = append(fields, field["name"])
		}
		return
	case "psql", "postgres", "postgresql":
		err = r.instance.Raw("SELECT column_name as name FROm information_schema.columns WHERE table_catalog = ? AND table_name = ?", schema, name).Scan(&f).Error
		if err != nil {
			return
		}
		for _, field := range f {
			fields = append(fields, field["name"])
		}
		return
	}
	return nil, nil
}

func (r *GormQuery) Scopes(funcs ...func(orm.Query) orm.Query) orm.Query {
	var gormFuncs []func(*gorm.DB) *gorm.DB
	for _, item := range funcs {
		gormFuncs = append(gormFuncs, func(db *gorm.DB) *gorm.DB {
			item(&GormQuery{instance: db})

			return db
		})
	}

	tx := r.instance.Scopes(gormFuncs...)

	return NewGormQuery(tx)
}

func getGormConfig(connection string) (gorm.Dialector, error) {
	driver := facades.Config.GetString("database.connections." + connection + ".driver")

	switch driver {
	case support.Mysql:
		return getMysqlGormConfig(connection), nil
	case support.Postgresql:
		return getPostgresqlGormConfig(connection), nil
	case support.Sqlite:
		return getSqliteGormConfig(connection), nil
	case support.Sqlserver:
		return getSqlserverGormConfig(connection), nil
	default:
		return nil, errors.New(fmt.Sprintf("err database driver: %s, only support mysql, postgresql, sqlite and sqlserver", driver))
	}
}

func getMysqlGormConfig(connection string) gorm.Dialector {
	dsn := support.GetMysqlDsn(connection)
	if dsn == "" {
		return nil
	}

	return mysql.New(mysql.Config{
		DSN: dsn,
	})
}

func getPostgresqlGormConfig(connection string) gorm.Dialector {
	dsn := support.GetPostgresqlDsn(connection)
	if dsn == "" {
		return nil
	}

	return postgres.New(postgres.Config{
		DSN: dsn,
	})
}

func getSqliteGormConfig(connection string) gorm.Dialector {
	dsn := support.GetSqliteDsn(connection)
	if dsn == "" {
		return nil
	}

	return sqlite.Open(dsn)
}

func getSqlserverGormConfig(connection string) gorm.Dialector {
	dsn := support.GetSqlserverDsn(connection)
	if dsn == "" {
		return nil
	}

	return sqlserver.New(sqlserver.Config{
		DSN: dsn,
	})
}
