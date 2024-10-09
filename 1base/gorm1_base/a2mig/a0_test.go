package a2mig

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
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

func TestRestaurantRows(t *testing.T) {
	db := getDB()
	rows, err := db.Model(&SimpleRestaurant{}).Rows()
	assert.Nil(t, err)
	assert.NotNil(t, rows)
	defer rows.Close()

	columns, err := rows.Columns()
	assert.Nil(t, err)
	t.Log(columns)

	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the result into the slice of interface{}
		err := rows.Scan(valuePtrs...)
		assert.Nil(t, err)

		// Print column names and values
		for i, col := range columns {
			val := values[i]
			valType := reflect.TypeOf(val)

			if valType == nil {
				fmt.Printf("%s: %v (Type: nil)\n", col, val)
			} else {
				// Convert []byte to string if applicable
				if byteSlice, ok := val.([]byte); ok {
					val = string(byteSlice)
					fmt.Printf("%s: %v (Type: %v, converted from []byte)\n", col, val, reflect.TypeOf(val))
				} else {
					fmt.Printf("%s: %v (Type: %v)\n", col, val, valType)
				}
			}
		}
		fmt.Println("---") // Separator between rows	}
		if err := rows.Err(); err != nil {
			assert.Error(t, err)
		}
	}
}

/*
values := make([]interface{}, len(columns))

创建一个切片 values，其长度等于列的数量。
每个元素的类型是 interface{}，这允许它存储任何类型的值。
这个切片将用来存储从数据库行中实际读取的值。
valuePtrs := make([]interface{}, len(columns))

创建另一个切片 valuePtrs，长度也等于列的数量。
这个切片将存储指向 values 切片中每个元素的指针。
for i := range columns { valuePtrs[i] = &values[i] }

遍历所有列。
对于每一列，将 valuePtrs 中的对应元素设置为 values 中相应元素的地址。
这种设置的目的和重要性：

动态类型处理：

由于我们不知道每列的具体类型，使用 interface{} 允许我们处理任何数据类型。
指针的使用：
rows.Scan() 方法需要接收指针作为参数，以便直接修改这些位置的值。
通过使用 &values[i]，我们提供了 values 切片中每个元素的地址。
间接层：

valuePtrs 作为一个间接层，允许 Scan 方法填充 values 切片，而不需要知道 values 的具体结构。
灵活性：

这种方法使得代码可以处理任意数量的列，而不需要预先知道表的结构。
使用示例：

err := rows.Scan(valuePtrs...)
这行代码会将数据库行的值扫描到 values 切片中，通过 valuePtrs 中存储的指针。之后，values 切片将包含该行的所有列值，可以进行进一步处理或显示。
这种技术是处理未知结构的数据库查询结果的常用方法，特别是在需要动态处理不同表或查询结果的情况下。

-----
这段 Go 代码片段用于从数据库查询结果中动态地获取数据列，并将这些列的值存储到一个接口类型的切片 ([]interface{}) 中。下面我来详细解释每一行代码的作用：
values := make([]interface{}, len(columns))
values 是一个 []interface{} 类型的切片（slice），大小与 columns 一样大，用于存储从数据库中读取的每一列的数据值。由于 interface{} 是 Go 语言中的空接口类型，它能够存储任何类型的值（例如 int, string, float64 等）。使用 len(columns) 来初始化切片的大小，确保能够容纳所有的列值。
valuePtrs := make([]interface{}, len(columns))
	•	valuePtrs 也是一个 []interface{} 类型的切片，大小同样为 len(columns)。它将保存每一列的指针，以便数据库扫描（如 rows.Scan）将查询到的数据存入对应的内存位置。
for i := range columns {
	valuePtrs[i] = &values[i]
}
	•	这个循环遍历 columns 切片的索引，并将 values 中的每个元素的指针存入 valuePtrs 中。即 valuePtrs[i] 存储的是 values[i] 的地址（即指针）。
	•	通过这种方式，后续调用 rows.Scan 时，扫描结果会被填充到 values[i] 中，而 valuePtrs[i] 则作为 Scan 函数的参数，指向 values[i] 的位置。

用途
这段代码通常在数据库查询结果的动态处理场景中使用。通过这种方式，程序能够动态处理未知数量和类型的列，常用于 database/sql 包中处理 SQL 查询结果。使用 rows.Scan(valuePtrs...) 语句时，可以将每列的值自动填入到 values 中，然后进一步处理这些值。
执行步骤
	1.	values 切片用于存放查询的每列数据。
	2.	valuePtrs 切片用于存放指向 values 切片中每个元素的指针。
	3.	通过 for 循环将 values 中每个元素的地址存入 valuePtrs。
	4.	使用 valuePtrs 作为 Scan 的参数，以便动态接收查询的结果。
通过这种设计，程序能够灵活处理数据库的返回值，而不需要提前知道列的数据类型或数量。
	values 切片：存储每个列的具体值，类型为 []interface{}，因为我们不知道列的具体数据类型。
	•	valuePtrs 切片：存储指向 values 中各个位置的指针，这些指针用于告诉 Scan 方法应该把查询结果的数据存储在哪里。
*/
