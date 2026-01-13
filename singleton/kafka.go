package singleton

import (
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
			GroupID:  "tracking-service-v1", // fix for Kafka cache
			MinBytes: 10e3,                  // 10KB
			MaxBytes: 10e6,                  // 10MB
			// Reader should be able to auto commit
			CommitInterval: time.Second,
		})
	})
	return kafkaReader
}

// CloseKafka cleans up the connection on shutdown
func CloseKafka() error {
	if kafkaReader != nil {
		return kafkaReader.Close()
	}
	return nil
}
