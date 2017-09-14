package main

import (
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

func main() {
	c := mqtt.NewClientOptions()
	log.Println(c)
}
