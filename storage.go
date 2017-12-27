package main

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	errplus "github.com/pkg/errors"
	"time"
)

type Storage struct {
	db *storm.DB
}

func NewStorage(filepath string) (*Storage, error) {
	db, err := storm.Open(filepath,
		storm.Codec(protobuf.Codec),
		storm.Batch())
	if err != nil {
		return nil, errplus.Wrap(err, "error while opening DB file")
	}
	return &Storage{
		db: db,
	}, nil
}

func (sto *Storage) CreateUser(config *UserConfig) (string, error) {
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

func (sto *Storage) FindByLogin(login string) ([]*User, error) {
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

func (sto *Storage) FindByID(ID string) (*User, error) {
	var user User
	transaction, err := sto.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.One("id", ID, &user); err != nil {
		return nil, err
	}
	if err = transaction.Commit(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (sto *Storage) FindeByDate(date time.Time) ([]*User, error) {
	var users []*User
	transaction, err := sto.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.Find("RegistrationDate", date, &users); err != nil {
		return nil, err
	}
	if err = transaction.Commit(); err != nil {
		return nil, err
	}
	return users, nil
}
