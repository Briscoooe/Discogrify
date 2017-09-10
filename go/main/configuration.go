package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Logger struct {
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxSize"`
		MaxBackups int    `json:"maxBackups"`
		MaxAge     int    `json:"maxAge"`
	}
	Redis struct {
		Host            string `json:"host"`
		Port            string `json:"port"`
		Password        string `json:"password"`
		Db              int    `json:"db"`
		HoursExpiration int    `json:"expiration"`
	} `json:"redis"`
	Spotify struct {
		RedirectURI  string `json:"redirectUri"`
		RequestDelay int    `json:"requestDelay"`
	}
	Cookie struct {
		CookieName string `json:"cookieName"`
		Expiration int    `json:"expiration"`
	}
	FilePath string `json:"filePath"`
}

func LoadConfiguration(file string) Configuration {
	var config Configuration

	configFile, err := os.Open(file)

	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
