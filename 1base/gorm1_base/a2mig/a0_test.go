package a2mig

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
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
	return db.Debug()
}
func TestMigrateTable1(t *testing.T) {
	db := getDB()
	err := db.AutoMigrate(&Cafeteria{}, &SimpleRestaurant{})
	if err != nil {
		panic(err)
	}
}

func TestAddValue(t *testing.T) {
	db := getDB()
	db = db.WithContext(context.Background())
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
	var err error
	t.Logf("new cafeteria id: %d, r1 id: %d, r2 id: %d", cafeNewValue.ID, r1.ID, r2.ID)
	//err := updateRestaurantSettings2(db, r1.ID, []string{"code"}, "update code")
	err = db.Model(SimpleRestaurant{}).Where("id=?", r1.ID).Update("settings", gorm.Expr(`jsonb_set(settings, '{code}', '"update code"'::jsonb, true)`)).Error
	assert.Nil(t, err)

	//err = updateRestaurantSettings2(db, r1.ID, []string{"contractInfo", "contact"}, "update contact")
	err = db.Model(SimpleRestaurant{}).Where("id=?", r1.ID).Update("settings", gorm.Expr(`jsonb_set(settings, ?, '"update contact"'::jsonb, true)`, "{contractInfo, contact}")).Error
	assert.Nil(t, err)
	err = db.Model(SimpleRestaurant{}).Where("id=?", r1.ID).Update("settings", gorm.Expr(`jsonb_set(settings, ?, ?::jsonb, true)`, "{contractInfo, value}", 123)).Error
	assert.Nil(t, err)

	err = db.Exec(`UPDATE simple_restaurants SET settings = jsonb_set(settings, '{contractInfo, value}', to_jsonb(1111111), true) WHERE id = ?`, r1.ID).Error

	//err = updateRestaurantSettings2(db, r1.ID, []string{"contractInfo", "value"}, 1213)
	assert.Nil(t, err)
	r1WithCafeteria := &RestaurantWithCafeteria{}
	err = db.Preload("Cafeteria").First(r1WithCafeteria, r1.ID).Error
	assert.Nil(t, err)
	spew.Dump(r1WithCafeteria)
	assert.NotNil(t, r1WithCafeteria.Cafeteria)
	assert.Equal(t, r1WithCafeteria.Cafeteria.Name, cafeNewValue.Name)

	rev = db.Model(SimpleRestaurant{}).Where("id=?", r1.ID).Update("settings", gorm.Expr(`jsonb_set(settings, ?, ?::jsonb, true)`, "{contractInfo}", RestaurantContractInfo{
		Contact:         "c",
		ContactPerson:   "p",
		TransferBank:    "b",
		TransferName:    "n",
		TransferAccount: "a",
		PaymentMethod:   "p",
		Values:          987,
	}))
	assert.Nil(t, rev.Error)

	cwithr := &CafeteriaWithRestaurants{}
	err = db.Preload("Restaurants").First(cwithr, cafeNewValue.ID).Error
	assert.Nil(t, err)
	t.Logf("cafeteria with restaurnts: %#v, r len: %d", cwithr, len(cwithr.Restaurants))
}

func updateRestaurantSettings1(db *gorm.DB, id int64, keys []string, value any) error {
	// Path to the subfield inside the JSONB column (here preferences.theme)
	jsonPath := fmt.Sprintf("{%s}", strings.Join(keys, ","))

	// GORM Raw SQL to update the subfield using PostgreSQL's jsonb_set function
	return db.Exec(`
        UPDATE simple_restaurants 
        SET settings = jsonb_set(settings, ?, to_jsonb(?)::jsonb, true)
        WHERE id = ?`, jsonPath, value, id).Error
}

func updateRestaurantSettings2(db *gorm.DB, id int64, keys []string, value any) error {
	// Create the JSONB path
	jsonPath := fmt.Sprintf("{%s}", strings.Join(keys, ","))

	return db.Model(SimpleRestaurant{}).Where("id=?", id).Update("settings", gorm.Expr("jsonb_set(settings, ?, ?::jsonb, true)", jsonPath, value)).Error
}
