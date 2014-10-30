package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"code.google.com/p/go-uuid/uuid"
	_ "github.com/lib/pq"
	"fmt"
)

type DatabaseConfiguration struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	DatabaseName string `json:"db_name"`
}

type Configuration struct {
	CookieHash string `json:"cookiehash"`
	Db DatabaseConfiguration `json:"database"`
}

var Settings *Configuration = LoadConfiguration()

func LoadConfiguration() *Configuration {
	_, err := os.OpenFile("isogon.json", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("isogon.json")

	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		var settings Configuration
		settings.Db.Host = "localhost"
		settings.Db.User = "postgres"
		settings.Db.Port = "5432"
		settings.CookieHash = uuid.New()

		jsonConfig, err := json.MarshalIndent(settings, "", "    ")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("isogon.json", jsonConfig, 0600)
		if err != nil {
			panic(err)
		}

		fmt.Println("Please set all appropriate settings in isogon.json")
		os.Exit(1);
	}

	var settings *Configuration
	if err := json.Unmarshal(data, &settings); err != nil {
		panic(err)
	}
	return settings
}
