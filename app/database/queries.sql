-- name: InsertCoordinate :one
INSERT INTO tracking_service.coordinate (
    user_id, coordinate
) VALUES (
    $1, $2
) 
RETURNING *;

-- name: InsertRideCoordinate :one
INSERT INTO tracking_service.ride_coordinate (
    user_id, ride_id, coordinate, velocity
) VALUES (
    $1, $2, $3, $4
) 
RETURNING *;

-- name: GetCoordinateByUser :many
SELECT * 
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

