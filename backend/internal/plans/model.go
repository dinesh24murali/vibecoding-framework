package plans

import "time"

type Plan struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProviderID  int64      `gorm:"column:provider_id;not null" json:"provider_id"`
	Name        string     `gorm:"column:name;size:120;not null" json:"name"`
	Description string     `gorm:"column:description;type:text;not null" json:"description"`
	Price       float64    `gorm:"column:price;not null" json:"price"`
	Discount    float64    `gorm:"column:discount;not null" json:"discount"`
	IsActive    bool       `gorm:"column:is_active;not null" json:"is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (Plan) TableName() string {
	return "plans"
}
