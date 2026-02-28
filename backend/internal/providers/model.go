package providers

import "time"

type Provider struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:name;size:120;not null" json:"name"`
	ImageURL  *string    `gorm:"column:image_url" json:"image_url,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (Provider) TableName() string {
	return "providers"
}
