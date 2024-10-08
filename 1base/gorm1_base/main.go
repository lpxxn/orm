package main

import (
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
	// db.Debug().AutoMigrate(&model.User{}, &model.Address{})
	db.Debug().AutoMigrate(&model.OrderUser{}, &model.Order{}, &model.Product{}, &model.OrderItem{})

	orderUser := &model.OrderUser{
		Name: "John Doe",
		Orders: []*model.Order{
			{
				Desc: "abc",
				Items: []*model.OrderItem{
					{
						ProductID: 1,
						Product: &model.Product{
							Name: "haha",
						},
					},
					{
						Product: &model.Product{
							Name: "test",
						},
					},
				},
			},
		},
	}
	_ = orderUser
	uDb := db.Create(orderUser)
	if uDb.Error != nil {
		panic(uDb.Error)
	}

	u1 := &model.OrderUser{}
	db.First(u1, 3)
	spew.Dump(u1)
	u2 := &model.OrderUser{}
	// db.Debug().Preload("orders").Preload("order_items").Preload("products").First(u2, 3)
	// The Preload("Orders") method tells GORM to load the related records for the field Orders in the User struct.
	db.Debug().Preload("Orders").First(u2, 3)
	spew.Dump(u2)

	u3 := &model.OrderUser{}
	db.Debug().Preload("Orders").Preload("Orders.Items").Preload("Orders.Items.Product").First(u3, 3)
	spew.Dump(u3)

}
