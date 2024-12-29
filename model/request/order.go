package request

import "time"

type OrderPage struct {
	Page
	Username string `query:"username" form:"username"`
	Gender   *int   `query:"gender" form:"gender"`
	Enable   *int   `query:"enable" form:"enable"`
}

type OrderAdd struct {
	TaskID            uint      `gorm:"primaryKey" json:"task_id"`
	Name              string    `json:"name"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	ResponsiblePerson string    `json:"responsible_person"`
	CreatedAt         time.Time `json:"created_at"`
	Orders            []Order   `gorm:"foreignKey:TaskID" json:"orders"` // 关联订单
}

type OrderUpdate struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
	Enable   bool   `json:"enable"`
	RoleIds  []int  `json:"roleIds"`
}
