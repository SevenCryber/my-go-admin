package request

import "time"

type OrderPage struct {
	Page
	Username string `query:"username" form:"username"`
	Gender   *int   `query:"gender" form:"gender"`
	Enable   *int   `query:"enable" form:"enable"`
}

type OrderAdd struct {
	Name             string         `json:"name"`
	StartTime        time.Time      `json:"start_time"`
	EndTime          time.Time      `json:"end_time"`
	OrderAuditPerson string         `json:"order_audit_person"`
	Kol              []KolDetailAdd `json:"kol"` // 关联订单
}

type KolDetailAdd struct {
	KolName  string `json:"kol_name"`
	KolPrice string `json:"kol_price"`
}

type OrderUpdate struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
	Enable   bool   `json:"enable"`
	RoleIds  []int  `json:"roleIds"`
}
