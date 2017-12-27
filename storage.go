package main

import (
	"github.com/asdine/storm"
)

type storage struct {
	db *storm.DB
}

func (sto *storage) CreateUser(config *UserConfig) (string, error) {
	user := config.User()
	transaction, err := sto.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.Save(user); err != nil {
		return "", err
	}
	if err = transaction.Commit(); err != nil {
		return "", err
	}
	return user.UUID, nil
}

func (sto *storage) FindByLogin(login string) ([]*User, error) {
	var users []*User
	transaction, err := sto.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.Find("Login", login, &users); err != nil {
		return nil, err
	}
	if err = transaction.Commit(); err != nil {
		return nil, err
	}
	return users, nil
}
