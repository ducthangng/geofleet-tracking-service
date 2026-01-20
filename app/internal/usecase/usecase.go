package usecase

import "github.com/ducthangng/GeoFleet/app/internal/usecase/usecase_dto"

type TrackingUsecaseService interface {
	UploadLocationHistory(data usecase_dto.DriverLocationEvent) (err error)
}
