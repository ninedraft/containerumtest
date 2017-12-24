package main

import (
	"github.com/asdine/storm"
)

type storage struct {
	db *storm.DB
}
