package main

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ryanzola/toll-calculator/aggregator/client"
	"github.com/ryanzola/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

// this can also be called KafkaTransport
// because the business logic is going to be abstracted
type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   *client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("Kafka consumer started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) Stop() {
	logrus.Info("Kafka consumer stopped")
	c.isRunning = false
	c.consumer.Close()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Kafka consume error %s", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err)
			continue
		}

		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("Error calculating distance: %s", err)
			continue
		}

		req := types.Distance{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUID: data.OBUID,
		}

		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregator error: %s", err)
			continue
		}
	}
}
