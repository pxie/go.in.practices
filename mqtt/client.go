package main

import (
	"log"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

const (
	keepAliveTimeout = 60 * time.Second
	pingTimeout      = 50 * time.Second
)

type client struct {
	c mqtt.Client
}

func (client *client) init(broker string, clientID string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(keepAliveTimeout)
	opts.SetPingTimeout(pingTimeout)
	opts.SetAutoReconnect(true)
	client.c = mqtt.NewClient(opts)
	if token := client.c.Connect(); token.Wait() && token.Error() != nil {
		log.Panicln(token.Error())
	}
	log.Printf("init client successfully. broker=%s, clientID=%s", broker, clientID)
}

type message struct {
	ts   time.Time
	body []byte
}

type producer client

func (p *producer) pub(topic string, qos byte, payload interface{}) {
	size := getPayloadSize(payload)
	log.Printf("publish message. topic=%s, qos=%v, payloadSize=%v", topic, qos, size)
	token := p.c.Publish(topic, qos, true, payload)

	token.Wait()
	if token.Error() != nil {
		log.Panicln(token.Error())
	}
}

func getPayloadSize(payload interface{}) int {
	var size int
	switch payload.(type) {
	case string:
		size = len(payload.(string))
	case []byte:
		size = len(payload.([]byte))
	default:
		log.Println("unknown payload type, set size to zero")
		size = 0
	}
	return size
}
