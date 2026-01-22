package registry

import (
	"context"
	"errors"

	"github.com/ducthangng/GeoFleet/app/internal/infrastructure"
	"github.com/ducthangng/GeoFleet/app/internal/interface/postgresql"
	"github.com/ducthangng/GeoFleet/app/internal/usecase"
	"github.com/ducthangng/GeoFleet/singleton"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

func ProvideDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	conn := singleton.GetConn()

	if conn == nil || conn.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	return conn.DB, nil
}

func ProvideRepository(db *pgxpool.Pool) *postgresql.Queries {
	return postgresql.New(db)
}

func ProvideRedis(ctx context.Context) *redis.Client {
	return singleton.GetRedisClient()
}

// ProvideUserUsecase provides the usecase struct
func ProvideUserUsecase(querier *postgresql.Queries, redis *redis.Client) *usecase.TrackingService {
	return usecase.NewTrackingService(querier, redis)
}

// ProvideUserUsecase provides the usecase struct
func ProvideKafkaReader(ctx context.Context) *kafka.Reader {
	kafkaReader := singleton.InitilizeKafkaReader()
	return kafkaReader
}

// ProviderSet groups these together (Optional, but clean)
var KafkaListenerSet = wire.NewSet(
	ProvideDBPool,
	ProvideRepository,
	ProvideRedis,
	ProvideUserUsecase,
	wire.Bind(new(usecase.TrackingUsecaseService), new(*usecase.TrackingService)),
	ProvideKafkaReader,
	infrastructure.NewKafkaConsumer,
)
