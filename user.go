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
