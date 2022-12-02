package orm

import (
	"context"
	"gorm.io/gorm"
)

type Orm interface {
	Connection(name string, config *gorm.Config, disableLog bool) Orm
	Query() DB
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
	Commit() error
	Rollback() error
}

type Query interface {
	Driver() Driver
	Count(count *int64) error
	Create(value any) error
	Delete(value any, conds ...any) error
	Distinct(args ...any) Query
	Exec(sql string, values ...any) error
	Find(dest any, conds ...any) error
	First(dest any) error
	FirstOrCreate(dest any, conds ...any) error
	ForceDelete(value any, conds ...any) error
	Get(dest any) error
	Group(name string) Query
	Having(query any, args ...any) Query
	Join(query string, args ...any) Query
	Limit(limit int) Query
	Model(value any) Query
	Offset(offset int) Query
	Order(value any) Query
	OrWhere(query any, args ...any) Query
	Pluck(column string, dest any) error
	Raw(sql string, values ...any) Query
	Save(value any) error
	Scan(dest any) error
	Scopes(funcs ...func(Query) Query) Query
	Select(query any, args ...any) Query
	Table(name string, args ...any) Query
	Update(column string, value any) error
	Updates(values any) error
	Where(query any, args ...any) Query
	WithTrashed() Query
}
