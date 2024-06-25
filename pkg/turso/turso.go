package turso

import (
	//"database/sql/driver"
	//"errors"
	"fmt"
	"os"
	"watermelon_king-api/models"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	sqlite "github.com/ytsruh/gorm-libsql"
	"gorm.io/gorm"
)

func Init() {
	var err error
	database := GetDB()

	err = database.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("Auto migrate failed")
		panic(err)
	}

	fmt.Println("Database initialized and migration completed")
}

func GetDB() *gorm.DB {
	dburl := os.Getenv("DB_URL")
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        dburl,
		DriverName: "libsql",
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database")
		panic(err)
	}
	return db
}
