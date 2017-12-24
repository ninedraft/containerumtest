package main

import (
	"github.com/satori/go.uuid"
	"time"
)

// User struct.
// To create User instance, use UserConfig
type User struct {
	UUID             string    `storm:"id"`
	Login            string    `storm:"index"`
	RegistrationDate time.Time `storm:"index"`
}

// UserConfig creates new User
type UserConfig struct {
	Login            string     `json:"login"`
	RegistrationDate *time.Time `json:"registration_date"`
}

// User method emits new User from UserConfig with same Login field.
// If RegistrationDate is nil, then time.Now() is used.
func (config *UserConfig) User() *User {
	if config.RegistrationDate == nil {
		now := time.Now()
		config.RegistrationDate = &now
	}
	return &User{
		UUID:             uuid.NewV4().String(),
		Login:            config.Login,
		RegistrationDate: *config.RegistrationDate,
	}
}
