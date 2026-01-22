package usecase

import (
	"context"

	"github.com/ducthangng/GeoFleet/app/internal/usecase/usecase_dto"
)

type TrackingUsecaseService interface {
	UploadLocationHistory(context.Context, usecase_dto.DriverLocationEvent) error
}
