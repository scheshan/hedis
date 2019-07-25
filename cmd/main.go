package main

import (
	"fmt"
	"time"

	"github.com/scheshan/hedis"
)

func main() {

	c := new(hedis.ServerConfig)
	c.Addr = ":16379"
	ser := hedis.NewServer(c)
	if err := ser.Start(); err != nil {
		fmt.Println(err)
	}

	<-time.After(100000 * time.Second)
}
