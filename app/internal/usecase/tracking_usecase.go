package usecase

import (
	"context"
	"log"

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

func (service *TrackingService) UploadLocationHistory(ctx context.Context, data usecase_dto.DriverLocationEvent) (err error) {
	// check if user currently in ride
	rideId := service.RedisService.Get(ctx, data.UserID.String()).Val()
	log.Println("DEBUG: rideId = ", rideId)

	rideUUID, err := cast.CastUUID(rideId)
	log.Println("DEBUG: uuid = ", rideId, "-- err: ", err)

	// invalid ride
	if len(rideId) == 0 || err != nil {
		log.Println("DEBUG: err != nil")
		// currently not in ride
		coordinate := postgresql.InsertCoordinateParams{
			UserID:    data.UserID,
			Longitude: data.Lng,
			Latitude:  data.Lat,
		}

		res, err := service.DataService.InsertCoordinate(ctx, coordinate)
		if err != nil {
			log.Println("DEBUG: insert: ", err)
			return err
		}

		if res.ID == 0 {
			log.Println("DEBUG: resID = 0")
		}

		log.Println("DEBUG: resID = ", res.ID)
		return nil
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
		log.Println("DEBUG: insert ride: ", err)
		return err
	}

	if res.ID == 0 {
		log.Println("DEBUG: UploadLocationHistory resID = 0")
	}

	log.Println("DEBUG: resID ride = ", res.ID)
	return nil
}
