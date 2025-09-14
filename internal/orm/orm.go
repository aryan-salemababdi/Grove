package orm

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(dsn string) *gorm.DB {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect database: %v", err)
	}
	log.Println("✅ Database connected successfully")
	return db
}

func DB() *gorm.DB {
	if db == nil {
		log.Fatal("❌ Database not initialized. Call orm.Init() first.")
	}
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("DB not initialized, call orm.Init() first")
	}
	return db
}
