package query

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/lpxxn/orm/1base/gorm1_base/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestQueryFirst(t *testing.T) {
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

	db.Debug().Model(u1q).Select("name").Updates(u1q)
}

func TestQueryFirst2(t *testing.T) {
	dsn := "host=localhost dbname=myorm sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
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
