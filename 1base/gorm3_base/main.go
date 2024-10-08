package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/lpxxn/orm/1base/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	println("hello world")
	dsn := "host=localhost dbname=myorm sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
	// 查询指定的列
	type OrderUser struct {
		Name string `gorm:"type:varchar(100);default:'Anonymous'"`
	}

	u1 := &OrderUser{}
	db.Debug().Model(&model.OrderUser{}).First(u1, 3)
	// SELECT "order_users"."name" FROM "order_users" WHERE "order_users"."id" = 3 AND "order_users"."deleted_at" IS NULL ORDER BY "order_users"."id" LIMIT 1
	spew.Dump(u1)

	fmt.Println("=============")
	u4 := []*model.OrderUser{}

	db.Debug().Preload("Orders", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, user_id")
	}).Preload("Orders.Items").Preload("Orders.Items.Product", func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc")
		// return db.Order("products.id desc")
	}).Select("id, email").Find(&u4)
	spew.Dump(u4)
}

type MyOrderUser struct {
	gorm.Model
	Name   string         `gorm:"type:varchar(100);default:'Anonymous'"`
	Orders []*model.Order `gorm:"foreignKey:UserID;references:ID"` // Explicit foreign key definition
}
