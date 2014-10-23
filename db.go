package main

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/go-martini/martini"
)

// Initialization function that handles basic migration functions for the domotisocks system.
// TODO: Make database settings configurable
func init() {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/domotisocks?charset=utf8&parseTime=True")

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
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/domotisocks?charset=utf8&parseTime=True")

	if(err != nil) {
		panic(err)
	}

	return func(c martini.Context) {
		c.Map(&db)
	}
}
