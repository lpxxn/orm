package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

}

func getDB() *gorm.DB {
	dsn := "host=localhost dbname=myorm sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
	return db.Debug()
}
