package main

import (
	"log"

	"go.in.practices/pkg/client"
)

func main() {
	c := lib.NewClient()

	c.Connect()
	c.Disconnect(5)
	c.IsConnected()
	c.GetID()

	log.Println("end game")
}
