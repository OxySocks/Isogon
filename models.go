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
	RegistrationTime time.Time
}

