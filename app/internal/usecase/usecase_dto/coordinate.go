package usecase_dto

import "github.com/google/uuid"

/*
 Tiêu chuẩn RFC mới là cái "Universal" (phổ quát).
 Thư viện của Google chỉ đơn giản là một implementation (bản thực thi) phổ biến nhất,
 đáng tin cậy nhất và được cộng đồng Go coi là "chuẩn mực" (de-facto standard) để làm việc với cái 16 bytes đó.
 Nó không phải là một cái "vỏ" bọc cho các loại UUID khác nhau.
 Nó là một công cụ giúp Go hiểu và thao tác với định dạng 128-bit đó một cách an toàn (Type-safe).

 ADS:
 + Tiết kiệm bộ nhớ: Một string cho UUID tốn ít nhất 36 bytes. Một uuid.UUID chỉ tốn đúng 16 bytes. Nếu bạn có hàng triệu driver_id trong memory, con số này rất đáng kể.
 + Validation cực sớm: Nếu một chuỗi không phải UUID hợp lệ chui vào hệ thống, thư viện sẽ báo lỗi ngay khi bạn cố gắng parse nó, thay vì đợi đến khi DB từ chối.
 + So sánh nhanh (Performance): So sánh hai mảng 16 bytes (integer comparison) luôn nhanh hơn so với so sánh chuỗi (string comparison).
*/

type DriverLocationEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Timestamp int64     `json:"timestamp"`
	Geohash   string    `json:"geohash"`
}
