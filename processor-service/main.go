package main

import (
	"fmt"
	"processor/consumer"
)

func main() {

	topic := "logs-topic"
	workerCount := 5

	consumer.StartConsumer(topic, workerCount)

	fmt.Println("Processor service started...")

	select {}
}