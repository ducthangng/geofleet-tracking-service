package singleton

import (
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaReader *kafka.Reader
	kafkaOnce   sync.Once
)

// GetKafkaWriter returns a singleton instance of the Kafka Producer
func InitilizeKafkaReader() *kafka.Reader {
	kafkaOnce.Do(func() {
		brokers := GetGlobalConfig().KafkaBrokers
		kafkaReader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  "tracking-service-v2026", // fix for Kafka cache
			Topic:    "ride_tracking.raw_coordinate",
			MinBytes: 1,    // 1 byte is good. Higher throughput, then change this param
			MaxBytes: 10e6, // 10MB
			// Reader should be able to auto commit
			CommitInterval: time.Second,
		})
	})

	log.Println("connect to kafka reader successfully")

	return kafkaReader
}

// CloseKafka cleans up the connection on shutdown
func CloseKafka() error {
	if kafkaReader != nil {
		return kafkaReader.Close()
	}
	return nil
}
