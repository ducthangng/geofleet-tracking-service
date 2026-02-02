package usecase

import (
	"context"

	"github.com/ducthangng/GeoFleet/app/internal/usecase/usecase_dto"
)

type TrackingUsecaseService interface {
	UploadLocationHistory(ctx context.Context, data usecase_dto.DriverLocationEvent) (insertedId int64, err error)
}
