package handler

import (
	"context"

	"github.com/ducthangng/GeoFleet/app/internal/domain/entity"
	"github.com/ducthangng/GeoFleet/app/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "github.com/ducthangng/geofleet-proto/gen/go/common/v1"
	tracking_v1 "github.com/ducthangng/geofleet-proto/gen/go/tracking/v1"
)

type TrackingHandler struct {
	TrackingUsecase usecase.TrackingUsecaseService
	tracking_v1.UnimplementedTrackingServiceServer
}

func NewTrackingHandler(usecase usecase.TrackingUsecaseService) *TrackingHandler {
	return &TrackingHandler{
		TrackingUsecase: usecase,
	}
}

func (trackingHandler *TrackingHandler) GetDriverNearby(ctx context.Context, input *tracking_v1.GetDriverNearbyRequest) (*tracking_v1.GetDriverNearbyResponse, error) {
	result, err := trackingHandler.TrackingUsecase.GetNearbyDriver(ctx, entity.Point{
		Latitude:  input.GetLocation().GetLat(),
		Longitude: input.GetLocation().GetLng(),
	})

	if err != nil {
		return nil, err
	}

	var driverLocations []*tracking_v1.DriverLocation
	for _, value := range result {
		driverLocations = append(driverLocations, &tracking_v1.DriverLocation{
			DriverId: value.UserID.String(),
			Role:     value.UserRole,
			Location: &commonv1.LatLng{
				Lat: value.Lat,
				Lng: value.Lng,
			},
			UpdatedAt: timestamppb.New(value.Timestamp),
		})
	}

	response := &tracking_v1.GetDriverNearbyResponse{
		Drivers: driverLocations,
	}

	return response, nil
}

func (trackingHandler *TrackingHandler) UploadLocationHistory(grpc.ClientStreamingServer[tracking_v1.UploadLocationHistoryRequest, tracking_v1.UploadLocationHistoryResponse]) error {
	return nil
}
