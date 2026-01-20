package usecase_dto

type DriverLocationEvent struct {
	UserID    string  `json:"user_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Timestamp int64   `json:"timestamp"`
	Geohash   string  `json:"geohash"`
}
