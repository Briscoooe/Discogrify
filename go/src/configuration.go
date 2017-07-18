package main

import (
	"os"
	"encoding/json"
)

type Configuration struct {
	Redis struct {
		Host		string	`json:"host"`
		Port		string	`json:"port"`
		Password	string	`json:"password"`
		Db			int	`json:"db"`

	} `json:"redis"`
	RequestDelay int `json:"requestDelay"`
}

func loadConfiguration(file string) Configuration {
	var config Configuration

	configFile, err := os.Open(file)

	defer configFile.Close()
	if err != nil {
		rollingLog.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	rollingLog.Println("Configurations loaded")
	return config
}