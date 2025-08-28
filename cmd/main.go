package main

import (
	"log"

	"redigo/internal/config"
	"redigo/internal/server"
)

func main() {
	log.Println("starting redigo server on", config.Port)
	if err := server.RunIoMultiplexingServer(); err != nil {
		log.Fatal(err)
	}
}
