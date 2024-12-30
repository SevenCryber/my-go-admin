package model

import "time"

type Order struct {
	OrderID          uint        `gorm:"primaryKey" json:"order_id"`
	Name             string      `json:"name"`
	CreateTime       time.Time   `json:"create_time" gorm:"create_time"`
	StartTime        time.Time   `json:"start_time"`
	EndTime          time.Time   `json:"end_time"`
	OrderAuditPerson string      `json:"order_audit_person" gorm:"order_audit_person"`
	Kol              []KolDetail `gorm:"foreignKey:OrderID" json:"kol"` // 关联订单
}

func (Order) TableName() string {
	return "orders" // 自定义表名为 "orders"
}

type KolDetail struct {
	KolID    uint   `gorm:"primaryKey" json:"kol_id"`
	OrderID  uint   `json:"order_id"`
	KolName  string `json:"kol_name"`
	KolPrice string `json:"kol_price"`
}

func (KolDetail) TableName() string {
	return "kol_detail" // 自定义表名为 "kol_detail"
}
