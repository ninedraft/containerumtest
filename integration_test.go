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
	case db == nil && err != nil:
		test.Fatalf("db is nil, error while creating storage: %v\n", err)
	case db != nil && err != nil:
		test.Fatalf("db is not nil, error while creating storage: %v\n", err)
	case db == nil && err == nil:
		test.Fatalf("db is nil and err is nil\n")
	default:
	}
}
