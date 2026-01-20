package usecase

import (
	"github.com/ducthangng/GeoFleet/app/internal/interface/postgresql"
	"github.com/ducthangng/GeoFleet/app/internal/usecase/usecase_dto"
)

type TrackingService struct {
	DataService postgresql.Querier
}

func NewTrackingService(service *postgresql.Queries) *TrackingService {
	return &TrackingService{
		DataService: service,
	}
}

func (service *TrackingService) UploadLocationHistory(data usecase_dto.DriverLocationEvent) (err error) {

	return nil
}
