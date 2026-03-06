package usecase

import (
	"context"
	"log"
	"time"

	"github.com/ducthangng/GeoFleet/app/internal/domain/entity"
	"github.com/ducthangng/GeoFleet/app/internal/interface/postgresql"
	"github.com/ducthangng/GeoFleet/app/internal/usecase/usecase_dto"
	"github.com/ducthangng/GeoFleet/service/cast"
	"github.com/ducthangng/GeoFleet/singleton"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	identity_v1 "github.com/ducthangng/geofleet-proto/gen/go/identity/v1"
	ridev1 "github.com/ducthangng/geofleet-proto/gen/go/ride/v1"
)

type TrackingService struct {
	DataService  postgresql.Querier
	RedisService *redis.Client
}

func NewTrackingService(service postgresql.Querier) *TrackingService {
	return &TrackingService{
		DataService:  service,
		RedisService: singleton.GetRedisClient(),
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

	// push to redis geohash
	if err = service.addGeoHash(ctx, data); err != nil {
		return
	}

	return res.ID, err
}

func (service *TrackingService) addGeoHash(ctx context.Context, location usecase_dto.DriverLocationEvent) error {

	// TODO: add const for these fields
	defaultKey := "drivers"
	if location.UserRole == identity_v1.UserRole_USER_ROLE_PASSENGER {
		defaultKey = "passengers"
	}

	defaultStatus := "on_ride_"
	if location.RideStatus == ridev1.RideStatus_RIDE_STATUS_UNSPECIFIED {
		defaultStatus = "active_"
	}

	if err := service.RedisService.GeoAdd(ctx, defaultStatus+defaultKey, &redis.GeoLocation{
		Name:      location.UserID.String(),
		Latitude:  location.Lat,
		Longitude: location.Lng,
	}); err != nil {
		log.Println("encounter error when adding geofleet: ", err)
		return nil
	}

	return nil

}

func (service *TrackingService) GetNearbyDriver(ctx context.Context, location entity.Point) (data []usecase_dto.DriverLocationEvent, err error) {
	searchResult := service.RedisService.GeoSearchLocation(ctx, "active_drivers", &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  location.Longitude,
			Latitude:   location.Latitude,
			Radius:     3,
			RadiusUnit: "km",
			Sort:       "ASC", // by distance
			Count:      10,    // Limit results
		},
		WithCoord: true,
		WithDist:  true,
		WithHash:  true,
	})

	listDriverIDs, err := searchResult.Result()
	if err != nil {
		return nil, err
	}

	for _, r := range listDriverIDs {
		log.Println("found: ", r.Name, " -- at: ", r.Latitude, " - ", r.Longitude)

		id, _ := uuid.Parse(r.Name)
		data = append(data, usecase_dto.DriverLocationEvent{
			UserID:     id,
			Lat:        r.Latitude,
			Lng:        r.Longitude,
			RideStatus: ridev1.RideStatus_RIDE_STATUS_UNSPECIFIED,
			UserRole:   identity_v1.UserRole_USER_ROLE_DRIVER,
			Timestamp:  time.Now(),
		})
	}

	return data, nil
}
