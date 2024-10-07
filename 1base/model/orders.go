package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderUser struct {
	gorm.Model
	Name   string
	Email  string
	Orders []*Order `gorm:"foreignKey:UserID;references:ID"` // Explicit foreign key definition
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
