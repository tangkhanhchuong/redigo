package main

import (
	"log"
	"strconv"

	"redigo/internal/config"
	"redigo/internal/server"
)

func main() {
	log.Println("starting redigo on :3000")
	if err := server.Run(config.Protocol, config.Host+":"+strconv.Itoa(config.Port)); err != nil {
		log.Fatal(err)
	}
}
