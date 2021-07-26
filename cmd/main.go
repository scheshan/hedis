package main

import (
	"hedis"
	"log"
	"time"
)

func main() {
	srv := hedis.NewStandardServer(nil)

	err := srv.Start()
	if err != nil {
		log.Fatal(err)
	}

	<-time.After(1000 * time.Minute)
}
