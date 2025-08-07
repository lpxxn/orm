package a1query

import (
	"reflect"
	"slices"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/lpxxn/orm/1base/gorm1_base/model"
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
func TestMigrate(t *testing.T) {
	db := getDB()
	// db.Debug().AutoMigrate(&model.User{}, &model.Address{})
	db.AutoMigrate(&model.OrderUser{}, &model.Order{}, &model.Product{}, &model.OrderItem{})
}

func TestQuery0(t *testing.T) {
	db := getDB()
	// db.Debug().AutoMigrate(&model.User{}, &model.Address{})
	//db.Debug().AutoMigrate(&model.OrderUser{}, &model.Order{}, &model.Product{}, &model.OrderItem{})

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

	u4 := &model.OrderUser{}
	db.Debug().Select(orderUserAllColumns).First(u4, 3)

	o1 := &model.Order{}
	db.Debug().Select(orderAllColumns).First(o1, 1)
}

var orderUserAllColumns = GetTableColumns(getDB(), &model.OrderUser{}, "id", "deleted_at")
var orderAllColumns = GetTableColumns(getDB(), &model.Order{}, "deleted_at")

func extractTagValues(input interface{}, tagName string, ignoreColumns ...string) []string {
	var tagValues []string
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	// Ensure the input is a struct
	if val.Kind() != reflect.Struct {
		return tagValues
	}

	// Traverse the struct fields and collect tag values
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		// Check if it's an embedded (anonymous) field
		if field.Anonymous {
			// Recursively collect tag values from the embedded struct
			embeddedTagValues := extractTagValues(val.Field(i).Interface(), tagName, ignoreColumns...)
			tagValues = append(tagValues, embeddedTagValues...)
		} else {
			// Get the tag value for the specified tag
			tagValue := field.Tag.Get(tagName)
			// check ignoreColumns
			if slices.Contains(ignoreColumns, tagName) {
				continue
			}
			if tagValue != "" {
				tagValues = append(tagValues, tagValue)
			}
		}
	}

	return tagValues
}

func GetTableColumns(db *gorm.DB, model interface{}, ignoreColumns ...string) []string {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(model)

	var columns []string
	for _, field := range stmt.Schema.Fields {
		//if !field.AutoIncrement && field.DBName != "" {
		if field.DBName != "" {
			// ignoreColumns
			if slices.Contains(ignoreColumns, field.DBName) {
				continue
			}
			columns = append(columns, field.DBName)
		}
	}

	return columns
}
