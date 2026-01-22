-- name: InsertCoordinate :one
INSERT INTO tracking_service.coordinate (
    user_id, coordinate
) VALUES (
    $1, 
    ST_SetSRID(ST_MakePoint(@longitude::float8, @latitude::float8), 4326)
) 
RETURNING id, ST_AsBinary(coordinate) as coordinate, created_at;

-- name: InsertRideCoordinate :one
INSERT INTO tracking_service.ride_coordinate (
    user_id, ride_id, coordinate, velocity
) VALUES (
    $1, $2, $3, $4
) 
RETURNING id, coordinate::bytea as coordinate, created_at;

-- name: GetCoordinateByUser :many
-- cast to bytea so that GO knows it must be []byte
SELECT id, user_id, coordinate::bytea as coordinate, created_at
FROM tracking_service.coordinate 
WHERE user_id = $1;

-- name: GetCoordinateInRide :many
SELECT * 
FROM tracking_service.ride_coordinate 
WHERE ride_id = $1; 

-- name: GetCoordinateInRideByUser :many
SELECT * 
FROM tracking_service.ride_coordinate 
WHERE ride_id = $1 and user_id = $2; 


-- Tại sao lại dùng ::bytea?
-- Trong Postgres, kiểu geography hoặc geometry là kiểu nội bộ. 
-- Khi sqlc sinh code, nếu bạn không ép kiểu về bytea, driver pgx có thể cố gắng 
-- parse nó thành các kiểu phức tạp hơn. Việc ép về bytea giúp dữ liệu truyền qua 
-- network ở dạng thô (binary), giúp hàm Scan của bạn hoạt động với hiệu suất cao nhất.
