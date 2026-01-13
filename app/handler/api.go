package handler

import tracking_v1 "github.com/ducthangng/geofleet-proto/gen/go/tracking/v1"

type TrackingHandler struct {
	tracking_v1.UnimplementedTrackingServiceServer
}

func NewTrackingHandler() *TrackingHandler {
	return &TrackingHandler{}
}

func (th *TrackingHandler) UploadLocationHistory(data tracking_v1.TrackingService_UploadLocationHistoryServer) error {

	return nil
}
