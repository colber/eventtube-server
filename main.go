package main

import (
	"os"
	"log"
	"path/filepath"
	"encoding/json"
	
	models "./models"

	server "./websoket"
)

func main() {
	config,err:=getConfig()
	if err != nil {
		log.Fatal(err)
    } 
	log.Println("Starting service",config.Service.Name,"(",config.Service.Mode,"mode) on",config.Service.Host,":",config.Service.Port,"...")
	wsServer := getServer(config)
	wsServer.Start()
	
	log.Println("Exiting...")
}

func getConfig() (*models.Config, error) {
	config := new(models.Config)

	_, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.Open("./config.json")
    defer configFile.Close()
    if err != nil {
        return nil,err
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)
	return config,nil
}
func getServer(config *models.Config) server.Runner {
	srv, err := server.Create(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accepting requests:")
	return srv
}
