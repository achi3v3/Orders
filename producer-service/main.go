package main

import (
	"fmt"
	"log"
	"producer-service/internal/config"
	"producer-service/internal/kafka/producer"
)

func main() {
	kafkaCfg, err := config.LoadKafkaConfig()
	if err != nil {
		log.Printf("main: %v", err)
		return
	}
	producer.ExternalSend(*kafkaCfg)
	fmt.Println("Success")
}
