package main

import (
	"chats-api/internal/config"
	"chats-api/internal/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if os.Getenv("APP_ENV") != "docker" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file ", err.Error())
			return
		}
	}

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal("error loading config ", err.Error())
	}
	srv, err := server.NewServer(conf)
	if err != nil {
		log.Fatal("error initializing server ", err.Error())
	}

	if err := srv.Start(); err != nil {
		log.Fatal("error starting server ", err.Error())
	}

}
