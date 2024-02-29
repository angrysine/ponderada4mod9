package main

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	godotenv "github.com/joho/godotenv"
)

var messagePubHandlerSub mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var text = fmt.Sprintf("Recebido: %s do tópico: %s com QoS: %d\n", msg.Payload(), msg.Topic(), msg.Qos())
	Writer("./logs/subscriber_logs.txt",  text+ "\n")
}

var connectHandlerSub mqtt.OnConnectHandler = func(client mqtt.Client) {
	Writer("subscriber_logs.txt", "connected" + "\n")
}

var connectLostHandlerSub mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	var text = fmt.Sprintf("Connection lost: %v", err)
	Writer("subscriber_logs.txt",  text+ "\n")
}

func Subscriber() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
	}

	var broker = os.Getenv("BROKER_ADDR")
	var port = 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker, port))
	opts.SetClientID("Subscriber")
	opts.SetUsername(os.Getenv("HIVE_USER"))
	opts.SetPassword(os.Getenv("HIVE_PSWD"))
	opts.SetDefaultPublishHandler(messagePubHandlerSub)
	opts.OnConnect = connectHandlerSub
	opts.OnConnectionLost = connectLostHandlerSub

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe("test/topic", 1, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
}