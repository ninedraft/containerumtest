package main

import (
	"github.com/satori/go.uuid"
	"time"
)

type user struct {
	UUID             uuid.UUID `storm:"unique"`
	Login            string    `storm:"index"`
	RegistrationDate time.Time `storm:"index"`
}

type userConfig struct {
	Login            string     `json:"login"`
	RegistrationDate *time.Time `json:"registration_date"`
}

func (config *userConfig) User() *user {
	if config.RegistrationDate == nil {
		now := time.Now()
		config.RegistrationDate = &now
	}
	return &user{
		UUID:             uuid.NewV4(),
		Login:            config.Login,
		RegistrationDate: *config.RegistrationDate,
	}
}
