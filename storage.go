package main

import (
	"github.com/asdine/storm"
	"github.com/satori/go.uuid"
)

type storage struct {
	db *storm.DB
}

func (sto *storage) CreateUser(config *UserConfig) (id *uuid.UUID, er error) {
	user := config.User()
	transaction, err := sto.db.Begin(true)
	if err = transaction.Save(user); err != nil {
		return nil, err
	}
	if err = transaction.Commit(); err != nil {
		return nil, err
	}
	return &user.UUID, nil
}
