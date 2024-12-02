package model

import (
	"github.com/SevenCryber/my-go-admin/initialize/datetime"
)

type User struct {
	Id         int `gorm:"autoIncrement"`
	Username   string
	Password   string
	Enable     bool
	CreateTime datetime.Datetime `gorm:"autoCreateTime"`
	UpdateTime datetime.Datetime `gorm:"autoUpdateTime"`
}
