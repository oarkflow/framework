package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	contractsorm "github.com/oarkflow/framework/contracts/database/orm"
	"github.com/oarkflow/framework/facades"
)

type Orm struct {
	Name            string
	ctx             context.Context
	connection      string
	defaultInstance contractsorm.DB
	instances       map[string]contractsorm.DB
	config          *gorm.Config
	disableLog      bool
	mu              *sync.RWMutex
}

func NewOrm(ctx context.Context, name string, config *gorm.Config, disableLog bool) contractsorm.Orm {
	orm := &Orm{ctx: ctx, Name: name, config: config, disableLog: disableLog, mu: &sync.RWMutex{}}

	return orm.Connection(name, config, disableLog)
}

func (r *Orm) Connection(name string, config *gorm.Config, disableLog bool) contractsorm.Orm {
	r.mu.Lock()
	defer r.mu.Unlock()
	defaultConnection := facades.Config.GetString("database.default")
	if name == "" {
		name = defaultConnection
	}

	r.connection = name
	if r.instances == nil {
		r.instances = make(map[string]contractsorm.DB)
	}

	if _, exist := r.instances[name]; exist {
		return r
	}

	g, err := NewGormDB(r.ctx, name, config, disableLog)
	if err != nil {
		color.Redln(fmt.Sprintf("[Orm] Init connection error, %v", err))

		return nil
	}
	if g == nil {
		return nil
	}

	r.instances[name] = g

	if name == defaultConnection {
		r.defaultInstance = g
	}

	return r
}

func (r *Orm) Query(database ...string) contractsorm.DB {
	r.mu.Lock()
	defer r.mu.Unlock()
	// the rationale behind this is that if the user passes in a database name, we will use that
	// to get the instance, otherwise we will use the default database
	if len(database) > 0 {
		// if the user passes in a database name, we will use that to get the instance
		// no matter how many database names are passed in, we will only use the first one
		instance, exist := r.instances[database[0]]
		if !exist {
			return nil
		}
		return instance
	}

	// if the user does not pass in a database name, we will use check if the connection is set
	// if it is set, we will use that to get the instance
	// if it is not set, we will use the default database
	if r.connection == "" {
		if r.defaultInstance == nil {
			r.Connection("", r.config, r.disableLog)
		}

		return r.defaultInstance
	}

	// get the instance from the connection if it is set
	instance, exist := r.instances[r.connection]
	if !exist {
		return nil
	}

	r.connection = ""

	return instance
}

func (r *Orm) Transaction(txFunc func(tx contractsorm.Transaction) error) error {
	tx, err := r.Query().Begin()
	if err != nil {
		return err
	}

	if err := txFunc(tx); err != nil {
		if err := tx.Rollback().Error; err != nil {
			return errors.Wrapf(err, "rollback error: %v", err)
		}

		return err
	} else {
		return tx.Commit().Error
	}
}

func (r *Orm) WithContext(ctx context.Context) contractsorm.Orm {
	return NewOrm(ctx, r.Name, r.config, r.disableLog)
}
