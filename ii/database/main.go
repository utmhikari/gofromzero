package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DBModel struct {
	*gorm.Model
	// You can add your own model components here
}

// DB database instance
var DBInstance *gorm.DB

// InitDB initialize DB
func InitDB() {
	log.Println("Initializing database...")
	db, err := gorm.Open("mysql", "root:123456@tcp(localhost:32773)/gofromzero?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	DBInstance = db
}

// CloseDB close DB
func CloseDB() {
	log.Println("Closing database...")
	err := DBInstance.Close()
	if err != nil {
		panic(err)
	}
}
