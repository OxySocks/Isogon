package main

import (
	"github.com/jinzhu/gorm"
	"github.com/go-martini/martini"
	"fmt"
	_ "github.com/lib/pq"
)

// Initialization function that handles basic database functions for the isogon system.
// This includes automatic database table creation and migration.
func init() {
	databaseString := "port=" + Settings.Db.Port + " host=" + Settings.Db.Host + " user=" + Settings.Db.User + " password=" + Settings.Db.Password + " dbname=" + Settings.Db.DatabaseName + " sslmode=disable"
	db, err := gorm.Open("postgres", databaseString)
	defer db.Close()

	if err != nil {
		panic(err)
	}

	// Create a 'measurements' table if it does not exist.
	if !db.HasTable(&Measurement{}) {
		fmt.Println("No measurement table found, creating one.")
		db.CreateTable(&Measurement{})
	}

	// Create a 'node' table if it does not exist.
	if !db.HasTable(&Node{}) {
		fmt.Println("No node table found, creating one.")
		db.CreateTable(&Node{})
	}

	// GORM automatic migration.
	db.AutoMigrate(&Measurement{}, &Node{})
}

// Creates a gorm.Db database handler for martini
func GormMiddleware() martini.Handler {
	databaseString := "port=" + Settings.Db.Port + " host=" + Settings.Db.Host + " user=" + Settings.Db.User + " password=" + Settings.Db.Password + " dbname=" + Settings.Db.DatabaseName + " sslmode=disable"
	db, err := gorm.Open("postgres", databaseString)

	if(err != nil) {
		panic(err)
	}

	return func(c martini.Context) {
		c.Map(&db)
	}
}
