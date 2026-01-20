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

	for {
		message, err := consumer.KafkaReader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var newMessage usecase_dto.DriverLocationEvent
		if err = sonic.Unmarshal(message.Value, newMessage); err != nil {
			log.Fatalf("encounter error when unmarshall kafka message: ", err)
		}

		go consumer.TrackingService.UploadLocationHistory(newMessage)
	}
}
