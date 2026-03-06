package usecase

import (
	"context"

	"github.com/ducthangng/GeoFleet/app/internal/domain/entity"
	"github.com/ducthangng/GeoFleet/app/internal/usecase/usecase_dto"
)

type TrackingUsecaseService interface {
	UploadLocationHistory(ctx context.Context, data usecase_dto.DriverLocationEvent) (insertedId int64, err error)

	GetNearbyDriver(ctx context.Context, location entity.Point) (data []usecase_dto.DriverLocationEvent, err error)
}
