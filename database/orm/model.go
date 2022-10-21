package orm

import (
	"time"

	"gorm.io/gorm"

	"github.com/sujit-baniya/framework/facades"
)

type Model struct {
	ID uint `gorm:"primaryKey"`
	Timestamps
}

type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Relationship struct {
}

func (r *Relationship) HasOne(dest, id any, foreignKey string) error {
	return facades.Orm.Query().Where(foreignKey+" = ?", id).Find(dest)
}

func (r *Relationship) HasMany(dest, id any, foreignKey string) error {
	return facades.Orm.Query().Where(foreignKey+" in ?", id).Find(dest)
}

func (r *Relationship) belongsTo(dest, id any) error {
	return facades.Orm.Query().Find(dest, id)
}
