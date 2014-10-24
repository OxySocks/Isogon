package main

import (
	"github.com/jinzhu/gorm"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
)

// Initialization function that handles basic migration functions for the domotisocks system.
// TODO: Make database settings configurable
func init() {
	databaseString := "port=" + Settings.DatabasePort + " host=" + Settings.DatabaseHost + " user=" + Settings.DatabaseUsername + " password=" + Settings.DatabasePassword + " dbname=" + Settings.DatabaseName + " sslmode=disable"
	db, err := gorm.Open("postgres", databaseString)

	if err != nil {
		panic(err)
	}

	db.CreateTable(&Measurement{})
	db.CreateTable(&Node{})
	db.AutoMigrate(&Measurement{}, &Node{})
}

// Martini handler to couple the GORM database to all route handlers.
// Allows gorm.DB to be used as a parameter
func dbMiddleware() martini.Handler {
	databaseString := "port=" + Settings.DatabasePort + " host=" + Settings.DatabaseHost + " user=" + Settings.DatabaseUsername + " password=" + Settings.DatabasePassword + " dbname=" + Settings.DatabaseName + " sslmode=disable"
	db, err := gorm.Open("postgres", databaseString)

	if(err != nil) {
		panic(err)
	}

	return func(c martini.Context) {
		c.Map(&db)
	}
}
