package main

import (
	"flag"
	"log"
)

func main() {
	filename := flag.String("db", "user.db", "DB file")
	storage, err := NewStorage(*filename)
	if err != nil {
		log.Fatalf("error while creating storage: %v\n", err)
	}
	service, err := NewService(storage)
	if err != nil {
		log.Fatalf("error while creating service: %v\n", err)
	}
	err = service.Start(":8222")
	if err != nil {
		log.Fatalf("error while startinggo service: %v\n", err)
	}
}
