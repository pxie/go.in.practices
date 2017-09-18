package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/pxie/paho.mqtt.golang"
)

const (
	keepAliveTimeout = 2 * time.Second
	pingTimeout      = 1 * time.Second
)

type client struct {
	c  mqtt.Client
	ID string
}

var count = 0

func handler(client mqtt.Client, msg mqtt.Message) {
	count++
	log.Println("count", count)
	payload := msg.Payload()
	size := getPayloadSize(payload)
	if size > 100 {
		log.Printf("default handler. clientID=%s, topic=%s, payloadSize=%v",
			client.GetClientID(), msg.Topic(), size)
	} else {
		log.Printf("default handler. clientID=%s, topic=%s, payload=%v",
			client.GetClientID(), msg.Topic(), string(payload))
	}
}

func (client *client) init(broker string, clientID string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(keepAliveTimeout)
	opts.SetPingTimeout(pingTimeout)
	// opts.SetAutoReconnect(true)
	// opts.SetCleanSession(true)
	opts.SetDefaultPublishHandler(handler)

	client.c = mqtt.NewClient(opts)
	client.ID = client.c.GetClientID()
	if token := client.c.Connect(); token.Wait() && token.Error() != nil {
		log.Panicln("connect broker error.", token.Error())
	}
	log.Printf("init client successfully. broker=%s, clientID=%s", broker, clientID)
}

func (client *client) pub(topic string, qos byte, payload interface{}) {
	size := getPayloadSize(payload)
	log.Printf("publish topic. clientID=%s, topic=%s, qos=%v, payloadSize=%v", client.ID, topic, qos, size)
	token := client.c.Publish(topic, qos, false, payload)

	token.Wait()
	if token.Error() != nil {
		log.Panicln("publish topic error.", token.Error())
	}
}

func (client *client) sub(topic string, qos byte, callback mqtt.MessageHandler) {
	text := ""
	if callback != nil {
		text = " with callback function"
	}
	log.Printf("subcribe topic%s. clientID=%s, topic=%s, qos=%v", text, client.ID, topic, qos)
	if token := client.c.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		log.Panicln("subcribe topic error.", token.Error())
	}
}

func (client *client) unsub(topic string) {
	log.Printf("unsubcribe topic. clientID=%s, topic=%s", client.ID, topic)
	if token := client.c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Panicln("unsubcribe topic error.", token.Error())
	}
}

func (client *client) distroy() {
	log.Printf("distroy mqtt client. clientID=%s", client.ID)
	client.c.Disconnect(10)
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

type message struct {
	ClientID string
	TS       time.Time
	Body     []byte
}

func getRandomBytes(size int) []byte {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		log.Panicln("generate random bytes error.", err)
	}
	log.Printf("generate rand bytes. bufSize=%v", len(buf))
	return buf
}

func createMsg(clientID string, size int) []byte {
	msg := &message{}
	msg.ClientID = clientID
	msg.TS = time.Now()
	msg.Body = getRandomBytes(size)
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Panicln("json marshal error.", err)
	}
	log.Printf("create message. contents=%v", string(bytes))

	return bytes
}
