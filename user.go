package main

import (
	"time"
)

type User struct {
	UUID             string
	Login            string
	RegistrationDate time.Time
}
