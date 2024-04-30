package main

import (
	"log"

	"github.com/ryanzola/toll-calculator/aggregator/client"
)

const (
	kafkaTopic         = "obu_data"
	aggregatorEndpoint = "http://127.0.0.1:3002/aggregate"
)

// Transport (HTTP, gRPC, Kafka) -> attach business logic to this transport

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
}
