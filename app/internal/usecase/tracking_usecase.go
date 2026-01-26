package usecase

import (
	"context"

	"github.com/ducthangng/GeoFleet/app/internal/domain/entity"
	"github.com/ducthangng/GeoFleet/app/internal/interface/postgresql"
	"github.com/ducthangng/GeoFleet/app/internal/usecase/usecase_dto"
	"github.com/ducthangng/GeoFleet/service/cast"
	"github.com/go-redis/redis/v8"
)

type TrackingService struct {
	DataService  postgresql.Querier
	RedisService *redis.Client
}

func NewTrackingService(service *postgresql.Queries, client *redis.Client) *TrackingService {
	return &TrackingService{
		DataService:  service,
		RedisService: client,
	}
}

func (service *TrackingService) UploadLocationHistory(ctx context.Context, data usecase_dto.DriverLocationEvent) (insertedId int64, err error) {
	// check if user currently in ride
	rideId := service.RedisService.Get(ctx, data.UserID.String()).Val()

	rideUUID, err := cast.CastUUID(rideId)
	if len(rideId) == 0 || err != nil {
		// currently not in ride
		coordinate := postgresql.InsertCoordinateParams{
			UserID:    data.UserID,
			Longitude: data.Lng,
			Latitude:  data.Lat,
		}

		res, err := service.DataService.InsertCoordinate(ctx, coordinate)
		if err != nil {
			return insertedId, err
		}

		return res.ID, nil
	}

	coordinate := postgresql.InsertRideCoordinateParams{
		UserID: data.UserID,
		RideID: rideUUID,
		Coordinate: entity.Point{
			Latitude:  data.Lat,
			Longitude: data.Lng,
		},
	}

	res, err := service.DataService.InsertRideCoordinate(ctx, coordinate)
	if err != nil {
		return insertedId, err
	}

	return res.ID, err
}
