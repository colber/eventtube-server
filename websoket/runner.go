package websoket

import (
	"../models"
	"./gorilla"
)

// Runner ...
type Runner interface {
	Start()
}

// Create ...
func Create(config *models.Config) (Runner, error) {
	var srv Runner
	var err error

	srv,err = gorilla.NewRunner(config)

	return srv, err
}