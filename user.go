package main

import (
	"github.com/satori/go.uuid"
	"time"
)

type user struct {
	UUID             uuid.UUID
	Login            string
	RegistrationDate time.Time
}
