package main

import (
	"github.com/asdine/storm"
	"github.com/satori/go.uuid"
)

type storage struct {
	db *storm.DB
}

func (sto *storage) CreateUser(config *UserConfig) (*uuid.UUID, error) {
	user := config.User()
	err := sto.db.Save(user)
	if err != nil {
		return nil, err
	}
	return &user.UUID, nil
}
