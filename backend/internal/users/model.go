package users

import "time"

type User struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string     `gorm:"column:name;size:120;not null"`
	PhoneNumber  string     `gorm:"column:phone_number;size:10;not null"`
	PasswordHash string     `gorm:"column:password_hash"`
	Role         string     `gorm:"column:role;size:20;not null"`
	Status       string     `gorm:"column:status;size:20;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;not null"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}
