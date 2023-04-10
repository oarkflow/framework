package orm

import (
	"context"

	"gorm.io/gorm"
)

type Orm interface {
	Connection(name string, config *gorm.Config, disableLog bool) Orm
	Query(database ...string) DB
	Transaction(txFunc func(tx Transaction) error) error
	WithContext(ctx context.Context) Orm
}

type DB interface {
	Query
	Begin() (Transaction, error)
	Instance() *gorm.DB
}

type Transaction interface {
	Query
	Commit() *gorm.DB
	Rollback() *gorm.DB
}

type Query interface {
	Driver() Driver
	Count(count *int64) *gorm.DB
	Create(value any) *gorm.DB
	Delete(value any, conds ...any) *gorm.DB
	Distinct(args ...any) Query
	Exec(sql string, values ...any) *gorm.DB
	Find(dest any, conds ...any) *gorm.DB
	First(dest any) *gorm.DB
	FirstOrCreate(dest any, conds ...any) *gorm.DB
	ForceDelete(value any, conds ...any) *gorm.DB
	Get(dest any) *gorm.DB
	Group(name string) Query
	Having(query any, args ...any) Query
	Join(query string, args ...any) Query
	Limit(limit int) Query
	Model(value any) Query
	Offset(offset int) Query
	Order(value any) Query
	OrWhere(query any, args ...any) Query
	Pluck(column string, dest any) *gorm.DB
	Raw(sql string, values ...any) Query
	Save(value any) *gorm.DB
	Scan(dest any) *gorm.DB
	Scopes(funcs ...func(Query) Query) Query
	Select(query any, args ...any) Query
	Table(name string, args ...any) Query
	Update(column string, value any) *gorm.DB
	Updates(values any) *gorm.DB
	Where(query any, args ...any) Query
	With(query string, args ...any) Query
	WithTrashed() Query
	Fields(schema, name string) (fields []string, err error)
}
