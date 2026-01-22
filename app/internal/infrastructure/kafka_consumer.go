package infrastructure

import (
	"context"
	"log"

	"github.com/bytedance/sonic"
	"github.com/ducthangng/GeoFleet/app/internal/usecase"
	"github.com/ducthangng/GeoFleet/app/internal/usecase/usecase_dto"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	TrackingService usecase.TrackingUsecaseService
	KafkaReader     *kafka.Reader
}

func NewKafkaConsumer(trackingService usecase.TrackingUsecaseService, reader *kafka.Reader) *KafkaConsumer {
	return &KafkaConsumer{
		TrackingService: trackingService,
		KafkaReader:     reader,
	}
}

func (consumer *KafkaConsumer) Listen(ctx context.Context) error {
	// log.Println("DEBUG: Kafka Consumer started, waiting for messages...")
	for {
		message, err := consumer.KafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Printf("ERROR: Kafka Read error: %v", err)
			return err
		}

		// log.Printf("DEBUG: Received Raw Message from Partition %d: %s", message.Partition, string(message.Value))

		var newMessage usecase_dto.DriverLocationEvent
		if err = sonic.Unmarshal(message.Value, &newMessage); err != nil {
			log.Printf("ERROR: Unmarshall failed: %v. Raw data: %s", err, string(message.Value))
			continue
		}

		// log.Printf("DEBUG: Successfully parsed message for User: %s", newMessage.UserID)
		go consumer.TrackingService.UploadLocationHistory(ctx, newMessage)
	}
}
