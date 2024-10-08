package query

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/lpxxn/orm/1base/gorm1_base/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestQuery1(t *testing.T) {
	println("hello world")
	dsn := "host=localhost dbname=myorm sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
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
	t.Log("=============")
	u4 := []*model.OrderUser{}
	db.Debug().Preload("Orders").Preload("Orders.Items").Preload("Orders.Items.Product").Find(&u4)
	spew.Dump(u4)
}
