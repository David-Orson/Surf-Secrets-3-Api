package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type DbConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

func LoadConfigDb(jsonConfig string) DbConfig {
	var configuration DbConfig
	err := json.Unmarshal([]byte(jsonConfig), &configuration)
	if err != nil {
		log.Println("0007: Could not unmarshal configuration from the passed jsonConfig")
		log.Fatal(err)
	}
	return configuration
}

func LoadFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Println("0005: Could not load database configuration")
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("0006: Could not read file '" + path + "'")
	}
	return string(bytes)
}
