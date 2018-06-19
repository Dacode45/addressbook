package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dacode45/addressbook/crypto"
	"github.com/Dacode45/addressbook/server"
	"github.com/Dacode45/addressbook/storage"
)

const (
	mongoURL           = "localhost"
	dbName             = "addressbook"
	userCollectionName = "user"
)

var config = server.ServerConfig{
	JWTSecret: "secret",
}

func main() {

	session, err := storage.NewMongoSession(mongoURL)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}

	hash := crypto.Hash{}
	uStorage := storage.NewMongoUserStorage(session.Copy(), dbName, userCollectionName, &hash)

	errChan := make(chan error)

	s := server.NewServer(uStorage, config)

	go func() {
		errChan <- s.Start()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			os.Exit(0)
		}
	}
}
