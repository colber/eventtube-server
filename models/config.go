package models

// Config ...
type Config struct {
	Service *Service 		`json:"service"`
}

// Service ...
type Service struct {
	Name 		string `json:"name"`
	Description string `json:"description"`
	Version 	string `json:"version"`
	Host 		string `json:"host"`
	Port 		string `json:"port"`
	Mode  		string `json:"mode"`
}



