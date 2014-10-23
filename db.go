package main

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/go-martini/martini"
)

func init() {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gormtest?charset=utf8&parseTime=True")

	if err != nil {
		panic(err)
	}

	db.CreateTable(&Measurement{})
	db.CreateTable(&Node{})
	db.AutoMigrate(&Measurement{}, &Node{})
}

func dbMiddleware() martini.Handler {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gormtest?charset=utf8&parseTime=True")

	if(err != nil) {
		fmt.Println(err)
	}



	return func(c martini.Context) {
		c.Map(&db)
	}
}
