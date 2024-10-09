package query

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/lpxxn/orm/1base/gorm1_base/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestQuery2(t *testing.T) {
	dsn := "host=localhost dbname=myorm sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}

	t.Log("======自定义Preload方法=======")
	u4 := []*model.OrderUser{}

	db.Debug().Preload("Orders", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, user_id")
	}).Preload("Orders.Items").Preload("Orders.Items.Product", func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc")
		// return db.Order("products.id desc")
	}).Select("id, email").Find(&u4)
	spew.Dump(u4)
}
