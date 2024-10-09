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

	//for rows.Next() {
	//	valuePtrs := make([]interface{}, len(columns))
	//	// Scan the result into the slice of interface{}
	//	for i := range columns {
	//		valuePtrs[i] = new(interface{})
	//	}
	//	err := rows.Scan(valuePtrs...)
	//	assert.Nil(t, err)
	//	for i, col := range columns {
	//		val := valuePtrs[i]
	//		valType := reflect.TypeOf(val)
	//
	//		if valType == nil {
	//			fmt.Printf("%s: %v (Type: nil)\n", col, val)
	//		} else {
	//			// Convert []byte to string if applicable
	//			if byteSlice, ok := val.([]byte); ok {
	//				val = string(byteSlice)
	//				fmt.Printf("%s: %v (Type: %v, converted from []byte)\n", col, val, reflect.TypeOf(val))
	//			} else {
	//				fmt.Printf("%s: %v (Type: %v)\n", col, val, valType)
	//			}
	//		}
	//	}
	//	fmt.Println("---") // Separator between rows	}
	//}
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns)) // valuePtrs := [&values[0], &values[1], &values[2]]
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
---------

首先，让我们看看如果不使用这个间接层，我们可能会怎么做：
values := make([]interface{}, len(columns))
err := rows.Scan(&values[0], &values[1], &values[2], ...) // 这是不可行的
这种方法有两个问题：

我们需要手动为每一列写一个 &values[i]，这在列数未知的情况下是不可能的。
更重要的是，values 切片中的元素类型是 interface{}，而 &values[i] 的类型是 *interface{}，这不是 Scan 方法所期望的。
现在，让我们看看为什么需要间接层：

类型匹配：
rows.Scan() 期望接收 *string, *int, *bool 等类型的指针，而不是 *interface{}。

动态分配：
我们需要一种方法来动态地为每一列创建一个合适的指针，而不需要预先知道列的类型或数量。

这就是 valuePtrs 发挥作用的地方：

values := make([]interface{}, len(columns))
valuePtrs := make([]interface{}, len(columns))
for i := range columns {
    valuePtrs[i] = &values[i]
}
这里发生的事情是：

values[i] 是一个 interface{} 类型的变量。
&values[i] 是一个指向 interface{} 的指针，类型为 *interface{}。
当我们将 &values[i] 赋值给 valuePtrs[i]（也是 interface{}类型）时，我们实际上是在存储一个指向 interface{} 的指针。
当我们调用 rows.Scan(valuePtrs...) 时：

Scan 方法接收到的是一系列 interface{} 类型的值。
但这些 interface{} 值实际上包含的是指向 values 切片中元素的指针。
Scan 可以通过这些指针将实际的数据写入 values 切片中的相应位置。
这个间接层允许我们：

动态处理任意数量的列。
让 Scan 方法能够正确地填充数据，而不需要知道每列的具体类型。
在 Scan 操作后，我们可以直接从 values 切片中读取数据。
总的来说，这个间接层解决了类型匹配的问题，同时提供了足够的灵活性来处理未知结构的查询结果。

-----
values := make([]interface{}, len(columns))
err := rows.Scan(values...)
这段代码看起来更简单，但实际上它无法正常工作。原因如下：
类型不匹配：
rows.Scan() 需要接收指针参数。它期望得到 *string, *int, *float64 等类型的指针，而不是 interface{} 类型的值。
无法修改值：
即使 Scan 接受 interface{} 类型，它也无法修改 values 切片中的实际值。因为在 Go 中，当你传递一个 interface{} 值时，你传递的是一个包含类型信息和值的副本，而不是原始值的引用。
interface{} 的局限性：
interface{} 本身不能存储值，它只能包含指向实际值的引用。如果我们直接传递 interface{} 给 Scan，它无法知道应该在哪里存储实际的数据。

-----
为什么不直接把数据 Scan 到 values 里呢？答案是 不行，因为 Scan 函数需要的是指针，而不是普通的变量。具体原因如下：
1. Scan 需要的是指针
rows.Scan 函数的工作原理是把数据库查询结果中的每一列的值填充到你提供的变量里。要让 Scan 能正确地修改这些变量，它需要知道这些变量在内存中的地址，也就是这些变量的指针。
如果你直接传入 values，它只是一个存储数据的普通切片，并不是指向数据位置的指针，Scan 不会知道怎么去修改它们。
因此，你需要传递每个 values[i] 的地址，这样 Scan 才能把数据存进去。

举个简单的例子
假设你有一个 int 变量，Scan 想把数据库中的某个值填进去。如果你直接传这个变量，Go 不知道在哪里放数据，因为你没告诉它这个变量在内存中的位置。你必须传递这个变量的地址（指针）：
var age int
rows.Scan(&age) // 传递的是 age 的地址 &age，而不是直接传 age
通过 &age，Scan 才知道要把结果放到 age 这个变量里。

2. values 是 interface{} 类型
values 是 []interface{}，是一个可以存储任何类型数据的切片。但是，在 Go 里，interface{} 只能存储值，而不能存储值的地址（指针）。而 rows.Scan 需要的正是值的地址。
假设 values 是这样定义的：
values := []interface{}{1, "Alice", 25}
	•	values[0] 是 1，是一个值，不是指针。
	•	values[1] 是 "Alice"，也是一个值，不是指针。
Scan 需要的不是这些值，而是这些值的存储位置。如果你直接传递 values，Scan 不知道把查询结果放在哪，因为它不能通过普通的值修改你想要的数据。
3. 为什么要用 valuePtrs
为了让 Scan 把数据正确放入 values 中，你需要一个存储每个 values[i] 的指针的切片，这就是 valuePtrs。
你做的事情其实是：

	1.	创建 values：用来存储每一列的数据。
	2.	创建 valuePtrs：用来存放每个 values[i] 的指针，告诉 Scan 数据该存到哪里。
for i := range columns {
    valuePtrs[i] = &values[i]  // 把每个 values[i] 的地址存到 valuePtrs
}
然后，当你调用 rows.Scan(valuePtrs...) 时：
	•	Scan 会通过 valuePtrs 里的指针把数据放到 values 中。
总结
	•	不能直接传 values 是因为 Scan 需要的是存储位置（指针），而不是值本身。
	•	你用 valuePtrs 存储每个 values[i] 的指针，让 Scan 能正确地把查询结果放到 values 里。
*/
