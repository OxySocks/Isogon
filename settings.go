package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"fmt"
)

type Configuration struct {
	DatabaseHost               string          `json:"db_host"`
	DatabasePort           string          `json:"db_port"`
	DatabaseUsername           string            `json:"db_username"`
	DatabasePassword         string          `json:"db_password"`
	DatabaseName	string `json:"db_name"`
}

var Settings *Configuration = LoadConfiguration()

func LoadConfiguration() *Configuration {
	_, err := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		var settings Configuration
		settings.DatabaseHost = "localhost"
		settings.DatabaseUsername = "postgres"
		settings.DatabasePort = "5432"

		jsonconfig, err := json.MarshalIndent(settings, "", "    ")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("config.json", jsonconfig, 0600)
		if err != nil {
			panic(err)
		}

		fmt.Println("Please set all appropriate settings in config.json")
		os.Exit(1);
	}

	var settings *Configuration
	if err := json.Unmarshal(data, &settings); err != nil {
		panic(err)
	}
	return settings
}
