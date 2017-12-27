// +build integration

package main

import (
	"log"
	"net"
	"net/rpc/jsonrpc"
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
		var id string
		err = storage.CreateUser(&UserConfig{login, nil}, &id)
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

func TestServiceCreateUser(test *testing.T) {
	addr := ":8222"
	filename := path.Join(test_dir, test.Name()+".db")
	storage, err := NewStorage(filename)
	if err != nil {
		test.Fatalf("error while creating storage: %v\n", err)
	}
	defer storage.db.Close()

	service, err := NewService(storage)
	if err != nil {
		test.Fatalf("error while creating service: %v\n", err)
	}
	go service.Start(":8222")

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		test.Fatalf("error while dialing connection: %v\n", err)
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)
	var repl interface{}
	err = client.Call("User.CreateUser",
		&UserConfig{
			"Merlin",
			nil,
		}, &repl)
	if err != nil {
		test.Fatalf("error while calling server: %v\n", err)
	}
	test.Logf("repy: %v\n", repl)
}
