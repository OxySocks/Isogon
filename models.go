package main

import "time"

type Node struct {
	Id int64 `json:"-"`
	HardwareAddress string `json:"hw_address"`
	Measurements []Measurement `json:"measurements"`
	CanonicalName string `json:"canonical_name"`
}

type Measurement struct {
	Id int64 `json:"-"`
	Temperature uint `json:"temperature"`
	Humidity uint `json:"humidity"`
	NodeId int64 `json:"-"`
	RegistrationTime time.Time `json:"time"`
}

type User struct {
	ID int64 `json:"id" gorm:"primary_key:yes"`
	Name string `json:"name" form:"name" binding:"required"`
	Password string `json:"-" form:"password" sql:"-"`
	Recovery string `json:"-"`
	PasswordHash []byte `json:"-"`
	Email string `json:"email" form:"email" binding:"required" sql:"unique"`
}
