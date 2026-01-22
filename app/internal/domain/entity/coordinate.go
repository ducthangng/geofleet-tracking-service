package entity

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// LatLng ~ tracking.proto
type Point struct {
	Longitude float64
	Latitude  float64
}

// Scan: DB ([]byte) -> Go (Struct)
func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for Point: %T", value)
	}

	// WKB cho Point thường là 21 bytes.
	// Nếu dùng PostGIS geography, đôi khi có thêm 4 bytes SRID (EWKB), tổng là 25 bytes.
	// Để an toàn, ta lấy 16 bytes cuối cùng.
	if len(b) < 21 {
		return errors.New("invalid WKB length for Point")
	}

	// Đọc 8 bytes cho Longitude và 8 bytes cho Latitude từ cuối mảng byte lên
	// Offset 5:9 cho Type, 9:21 cho Coordinates (nếu là EWKB thì offset sẽ khác)
	// Cách an toàn nhất là đọc từ cuối ngược lại:
	n := len(b)
	p.Latitude = math.Float64frombits(binary.LittleEndian.Uint64(b[n-8:]))
	p.Longitude = math.Float64frombits(binary.LittleEndian.Uint64(b[n-16 : n-8]))

	return nil
}

// Value: Go (Struct) -> DB ([]byte)
func (p Point) Value() (driver.Value, error) {
	// Chúng ta sẽ tạo chuẩn WKB (21 bytes) để Postgres hiểu
	buf := make([]byte, 21)

	buf[0] = 1                                 // LittleEndian
	binary.LittleEndian.PutUint32(buf[1:5], 1) // Type 1 = Point

	// Convert Float64 sang bits (Uint64) rồi mới Put vào buffer
	binary.LittleEndian.PutUint64(buf[5:13], math.Float64bits(p.Longitude))
	binary.LittleEndian.PutUint64(buf[13:21], math.Float64bits(p.Latitude))

	return buf, nil
}
