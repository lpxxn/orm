package mig1

import (
	"math/rand"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	dsn := "host=localhost dbname=myorm sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
	return db
}
func TestMigrateTable1(t *testing.T) {
	db := getDB().Debug()
	err := db.AutoMigrate(&Cafeteria{}, &SimpleRestaurant{})
	if err != nil {
		panic(err)
	}
}

func TestAddValue(t *testing.T) {
	db := getDB().Debug()

	newCafeteria := Cafeteria{
		SnowflakeID: rand.Int63(),
		Name:        "东北大食堂",
		Settings: &CafeteriaSettings{
			Status:  true,
			Address: "黑龙江",
			Code:    "cxcsd",
		},
		CityID:   1,
		Subtitle: "sub",
	}
	r1 := &SimpleRestaurant{
		SnowflakeID: rand.Int63(),
		RefID:       rand.Int63(),
		Name:        "KFC",
		Settings: RestaurantSettings{
			Code:        "A",
			Keyword:     []string{"a", "1"},
			EnableImage: true,
			ContractInfo: RestaurantContractInfo{
				Contact:         "张",
				ContactPerson:   "三",
				TransferBank:    "中银",
				TransferName:    "abcd",
				TransferAccount: "",
				PaymentMethod:   "",
			},
		},
	}
	r2 := &SimpleRestaurant{
		SnowflakeID: rand.Int63(),
		RefID:       rand.Int63(),
		Name:        "火🔥旺",
		Settings: RestaurantSettings{
			Code: "B",
			ContractInfo: RestaurantContractInfo{
				Contact:         "",
				ContactPerson:   "李四",
				TransferBank:    "建设",
				TransferName:    "银行",
				TransferAccount: "A1223",
				PaymentMethod:   "Afsdf",
			},
		},
	}
	cafeNewValue := &CafeteriaWithRestaurants{
		Cafeteria:   newCafeteria,
		Restaurants: []*SimpleRestaurant{r1, r2},
	}
	rev := db.Create(cafeNewValue)
	if rev.Error != nil {
		panic(rev.Error)
	}

	t.Logf("new cafeteria id: %d, r1 id: %d, r2 id: %d", cafeNewValue.ID, r1.ID, r2.ID)

	r1WithCafeteria := &RestaurantWithCafeteria{}
	err := db.Preload("Cafeteria").First(r1WithCafeteria, r1.ID).Error
	assert.Nil(t, err)
	spew.Dump(r1WithCafeteria)
	assert.NotNil(t, r1WithCafeteria.Cafeteria)
	assert.Equal(t, r1WithCafeteria.Cafeteria.Name, cafeNewValue.Name)

	cwithr := &CafeteriaWithRestaurants{}
	err = db.Preload("Restaurants").First(cwithr, cafeNewValue.ID).Error
	assert.Nil(t, err)
	t.Logf("cafeteria with restaurnts: %#v, r len: %d", cwithr, len(cwithr.Restaurants))
}

func TestQueryPreload(t *testing.T) {
	//dsn := "host=localhost dbname=myorm sslmode=disable TimeZone=Asia/Shanghai"
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	//if err != nil {
	//	panic(err)
	//}
	//cafeteriaRestaurants := &CafeteriaWithRestaurants{}
	//err = db.Debug().Preload("Restaurants").Where("")
}
