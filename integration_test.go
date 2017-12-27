// +build integration

package main

import (
	"log"
	"os"
	"path"
	"testing"
)

const (
	test_dir = "test_files"
)

func init() {
	err := os.Mkdir(test_dir, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("error while creating test_file dir: %v\n", err)
	}
}

func TestNewStorage(test *testing.T) {
	filename := path.Join(test_dir, test.Name()+".db")
	storage, err := NewStorage(filename)
	switch {
	case storage == nil && err != nil:
		test.Fatalf("db is nil, error while creating storage: %v\n", err)
	case storage != nil && err != nil:
		test.Fatalf("db is not nil, error while creating storage: %v\n", err)
	case storage == nil && err == nil:
		test.Fatalf("db is nil and err is nil\n")
	default:
	}
}

func TestCreateUser(test *testing.T) {
	filename := path.Join(test_dir, test.Name()+".db")
	storage, err := NewStorage(filename)
	if err != nil {
		test.Fatalf("error while creating storage: %v\n", err)
	}
	defer storage.db.Close()

	logins := []string{"crook", "whale", "estrogen", "koala", "cursor",
		"populist", "gym", "server", "garden", "game", "bankbook",
		"purse", "prosecution", "desert", "forearm", "knuckle"}
	for _, login := range logins {
		id, err := storage.CreateUser(&UserConfig{
			login,
			nil,
		})
		switch {
		case id == "" && err == nil:
			test.Fatalf("nil id and err while creating user %s\n", login)
		case id == "" && err != nil:
			test.Fatalf("error while creating user %s: %v\n", login, err)
		case id != "" && err != nil:
			test.Fatalf("error while creating user %s, id is not nil: %s %v\n",
				login, id, err)
		default: // id != "" && err == nil
		}
	}
}
