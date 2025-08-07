package a1query

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/lpxxn/orm/1base/gorm1_base/model"
	"gorm.io/gorm"
)

func TestQueryFirst(t *testing.T) {
	db := getDB()
	// 查询指定的列
	type OrderUser struct {
		Name  string `gorm:"type:varchar(100);default:'Anonymous'"`
		Email string
	}
	u1 := &OrderUser{}
	db.Debug().Model(&model.OrderUser{}).First(u1, 3)
	// SELECT "order_users"."name" FROM "order_users" WHERE "order_users"."id" = 3 AND "order_users"."deleted_at" IS NULL ORDER BY "order_users"."id" LIMIT 1
	spew.Dump(u1)
	u1q := model.OrderUser{
		Model: gorm.Model{ID: 2},
		Name:  "heiheihei",
	}
	db.Debug().First(&u1q)
	// SELECT * FROM "order_users" WHERE "order_users"."deleted_at" IS NULL AND "order_users"."id" = 2 ORDER BY "order_users"."id" LIMIT 1
	spew.Dump(u1q)

	db.Debug().Where("id = 1").First(&u1q) // 会 id = 1 and id = 2， 要小心哇~~~~ 啥玩意儿
	// SELECT * FROM "order_users" WHERE id = 1 AND "order_users"."deleted_at" IS NULL AND "order_users"."id" = 2 ORDER BY "order_users"."id" LIMIT 1
	spew.Dump(u1q)

	u1q = model.OrderUser{
		Model: gorm.Model{ID: 2},
		Name:  "heiheihei",
		Email: "aaa@heihei.com",
	}
	db.Debug().Model(u1q).Select("Email").Updates(u1q)
	// 指定只更新字段 email，其他字段即使有值也不会更新。
	// Email
	// UPDATE "order_users" SET "updated_at"='2025-08-07 23:05:28.917',"email"='aaa@heihei.com' WHERE "order_users"."deleted_at" IS NULL AND "id" = 2
	// Select("name")  只更新 name 字段
	//Select("name", "email")  同时更新 name 和 email
	//Omit("desc") 更新除 desc 外的字段
	//不写 Select 或 Omit
	//GORM 会自动判断“非零值”字段来更新
}

func TestQueryFirst2(t *testing.T) {
	db := getDB()
	// 查询指定的列
	u1qResult := model.OrderUser{}
	db.Debug().Model(MyOU{ID: 2}).First(&u1qResult) //没啥用哇，这个u1q,只有update的时候好使，但文档里说是可以的，妈的
	spew.Dump(u1qResult)

	db.Debug().First(&u1qResult)
	//  SELECT * FROM "order_users" WHERE "order_users"."deleted_at" IS NULL AND "order_users"."id" = 1 ORDER BY "order_users"."id" LIMIT 1
	spew.Dump(u1qResult)
}

type MyOU struct {
	ID uint `gorm:"primarykey"`
}

func (MyOU) TableName() string {
	return "order_users"
}
