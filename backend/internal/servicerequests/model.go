package servicerequests

import "time"

type ServiceRequest struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"column:user_id;not null" json:"user_id"`
	PlanID    int64      `gorm:"column:plan_id;not null" json:"plan_id"`
	Status    string     `gorm:"column:status;size:30;not null" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (ServiceRequest) TableName() string {
	return "service_requests"
}
