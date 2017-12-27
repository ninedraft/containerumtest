// +build integration

package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path"
	"testing"
)

const (
	test_dir = "test_files"
	addr     = ":8222"
)

var (
	logins = []string{"crook", "whale", "estrogen", "koala", "cursor",
		"populist", "gym", "server", "garden", "game", "bankbook",
		"purse", "prosecution", "desert", "forearm", "knuckle"}
)

func init() {
	err := os.Mkdir(test_dir, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("error while creating test_file dir: %v\n", err)
	}
}

func createTestStorage(test *testing.T) (*Storage, func()) {
	filename := path.Join(test_dir, test.Name()+".db")
	storage, err := NewStorage(filename)
	switch {
	case storage == nil && err != nil:
		test.Fatalf("db is nil, error while creating storage: %v\n", err)
	case storage != nil && err != nil:
		test.Fatalf("db is not nil, error while creating storage: %v\n", err)
	case storage == nil && err == nil:
		test.Fatalf("db is nil and err is nil\n")
	default: // storage != nil && err == nil
	}
	return storage, func() {
		if err := storage.db.Close(); err != nil {
			test.Errorf("error while closing test storage: %v\n", err)
		}

		if err = os.Remove(filename); err != nil {
			test.Errorf("error while removing test db file: %v\n", err)
		}
	}
}

func TestNewStorage(test *testing.T) {
	filename := path.Join(test_dir, test.Name()+".db")
	storage, err := NewStorage(filename)
	defer func() {
		storage.db.Close()
		if err := os.Remove(filename); err != nil {
			test.Errorf("error while removing test db file: %v\n", err)
		}
	}()

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
	storage, drop := createTestStorage(test)
	defer drop()

	for _, login := range logins {
		var id string
		err := storage.CreateUser(UserConfig{
			Login:            login,
			RegistrationDate: nil,
			UUID:             "",
		}, &id)
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

func TestRPC(test *testing.T) {
	storage, drop := createTestStorage(test)
	defer drop()
	service := testServiceCreateUser(test, storage)
	testMethod(test, service, findTestUser)
	testMethod(test, service, updateUser)
}

func testServiceCreateUser(test *testing.T, storage *Storage) *Service {
	service, err := NewService(storage)
	if err != nil {
		test.Fatalf("error while creating service: %v\n", err)
	}
	go service.Start(addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		test.Fatalf("error while dialing connection: %v\n", err)
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)
	createUser(test, client)
	return service
}

func testMethod(test *testing.T, service *Service, method func(*testing.T, *rpc.Client)) {
	conn := dialTestConnection(test)
	defer conn.Close()
	client := jsonrpc.NewClient(conn)
	method(test, client)
}

func dialTestConnection(test *testing.T) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		test.Fatalf("error while dialing connection: %v\n", err)
	}
	return conn
}

func createUser(test *testing.T, client *rpc.Client) {
	var repl string
	for _, login := range logins {
		err := client.Call("user.CreateUser",
			&UserConfig{
				Login:            login,
				RegistrationDate: nil,
				UUID:             "",
			}, &repl)
		if err != nil {
			test.Fatalf("error while calling server: %v\n", err)
		}
		//test.Logf("reply: %v\n", repl)
	}

}

func findTestUser(test *testing.T, client *rpc.Client) {
	var replFind []User
	for _, login := range logins {
		err := client.Call("user.FindByLogin", login, &replFind)
		if err != nil {
			test.Fatalf("error while calling user.FindByLogin: %v\n", err)
		}
		var lg string
		for n, user := range replFind {
			lg += fmt.Sprintf("%d %s %s\n", n, user.Login, user.UUID)
		}
		//test.Logf(lg)
	}
}

func updateUser(test *testing.T, client *rpc.Client) {
	var ID string
	err := client.Call("user.CreateUser",
		&UserConfig{
			Login:            "oldLogin",
			RegistrationDate: nil,
			UUID:             "",
		}, &ID)

	if err != nil {
		test.Fatalf("error while calling user.CreateUser: %v\n", err)
	}

	var status string
	newUserData := &UserConfig{
		Login:            "newLogin",
		RegistrationDate: nil,
		UUID:             ID,
	}
	err = client.Call("user.UpdateUser", newUserData, &status)
	if err != nil {
		test.Fatalf("error while calling user.UpdateUser: %v\n", err)
	}

	var replFind User
	err = client.Call("user.FindByID", ID, &replFind)
	if err != nil {
		test.Fatalf("error while calling user.FindByID: %v\n", err)
	}
	if replFind.Login != newUserData.Login {
		test.Fatalf("login is not updated!")
	}
}
