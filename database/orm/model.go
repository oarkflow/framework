package orm

import (
	"time"

	"gorm.io/gorm"

	"github.com/oarkflow/framework/facades"
)

type Model struct {
	ID uint `gorm:"primaryKey"`
	Timestamps
}

type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Relationship struct {
}

func (r *Relationship) HasOne(dest, id any, foreignKey string) *gorm.DB {
	return facades.Orm.Query().Where(foreignKey+" = ?", id).Find(dest)
}

func (r *Relationship) HasMany(dest, id any, foreignKey string) *gorm.DB {
	return facades.Orm.Query().Where(foreignKey+" in ?", id).Find(dest)
}

func (r *Relationship) belongsTo(dest, id any) *gorm.DB {
	return facades.Orm.Query().Find(dest, id)
}
