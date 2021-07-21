package main

import (
	"hedis/server"
	"log"
	"time"
)

func main() {
	srv := server.NewStandard(nil)

	err := srv.Run()
	if err != nil {
		log.Fatal(err)
	}

	<-time.After(1000 * time.Minute)
}
