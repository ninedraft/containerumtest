package main

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	errplus "github.com/pkg/errors"
	"time"
)

var (
	ErrUUIDIsNotProvided = fmt.Errorf("user UUID is not provided")
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

func (sto *Storage) CreateUser(config UserConfig, id *string) error {
	user := config.User()
	transaction, err := sto.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.Save(user); err != nil {
		return err
	}
	if err = transaction.Commit(); err != nil {
		return err
	}
	*id = user.UUID
	return nil
}

func (sto *Storage) FindByLogin(login string, retUsers *[]*User) error {
	var users []*User
	transaction, err := sto.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.Find("Login", login, &users); err != nil {
		return err
	}
	if err = transaction.Commit(); err != nil {
		return err
	}
	*retUsers = users
	return nil
}

func (sto *Storage) FindByID(ID string, retUser *User) error {
	var user User
	transaction, err := sto.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.One("UUID", ID, &user); err != nil {
		return err
	}
	if err = transaction.Commit(); err != nil {
		return err
	}
	*retUser = user
	return nil
}

func (sto *Storage) FindeByDate(date time.Time, retUsers *[]*User) error {
	var users []*User
	transaction, err := sto.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.Find("RegistrationDate", date, &users); err != nil {
		return err
	}
	if err = transaction.Commit(); err != nil {
		return err
	}
	*retUsers = users
	return nil
}

func (storage *Storage) UpdateUser(user UserConfig, ret *string) error {
	if user.UUID == "" {
		return ErrUUIDIsNotProvided
	}
	transaction, err := storage.db.Begin(true)
	defer transaction.Rollback()
	if err = transaction.Update(user.User()); err != nil {
		return err
	}
	if err = transaction.Commit(); err != nil {
		return err
	}
	*ret = "ok"
	return nil
}
