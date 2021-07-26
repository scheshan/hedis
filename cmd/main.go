package main

import (
	"hedis/server"
	"log"
	"time"
)

func main() {
	srv := server.NewStandardServer(nil)

	err := srv.Start()
	if err != nil {
		log.Fatal(err)
	}

	<-time.After(1000 * time.Minute)
}
