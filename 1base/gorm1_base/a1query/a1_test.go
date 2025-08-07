package a1query

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/lpxxn/orm/1base/gorm1_base/model"
)

func TestQuery1(t *testing.T) {
	db := getDB()

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
	//spew.Dump(u3)
	t.Log(u3)
	t.Log("=============")
	u4 := []*model.OrderUser{}
	db.Debug().Preload("Orders").Preload("Orders.Items").Preload("Orders.Items.Product").Find(&u4)
	//spew.Dump(u4)
	t.Log(u4)
	t.Log("====Orders.Items=========")
	u3 = &model.OrderUser{}
	db.Debug().Preload("Orders.Items").First(u3, 3)
	spew.Dump(u3)
}

/*
Preload("A.B") // 表示：预加载 A 的同时，预加载 A 下的 B
不能有 gorm:"-" 标签否则字段会被 GORM 忽略，无法 preload

*/
