package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderUser struct {
	gorm.Model
	Name   string   `gorm:"type:varchar(100);default:'Anonymous'"`
	Email  string   `gorm:"type:varchar(100);default:'a@b.com'"`
	Desc   string   `gorm:"size:10"`
	Orders []*Order `gorm:"foreignKey:UserID;references:ID"` // Explicit foreign key definition
	// // Order.UserID 关联到 OrderUser.ID
	// 这个关联的 SQL 查询时，Order 表里的 UserID 字段要对应 OrderUser 表里的 ID 字段来进行匹配。
	//foreignKey:UserID 表示 Order 这个模型中有个字段叫 UserID
	//references:ID 表示要去找的主模型（OrderUser）里的字段是 ID
}

type Order struct {
	gorm.Model
	UserID    uint
	User      *OrderUser //`gorm:"-"`
	OrderDate time.Time
	Desc      string
	Items     []*OrderItem
}

type Product struct {
	gorm.Model
	Name  string
	Price float64
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	Order     Order //`gorm:"-"`
	ProductID uint
	Product   *Product //`gorm:"-"`
	Quantity  int
}
