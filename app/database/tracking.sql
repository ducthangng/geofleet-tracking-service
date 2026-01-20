CREATE SCHEMA IF NOT EXISTS tracking_service;

CREATE TABLE IF NOT EXISTS tracking_service.coordinate (
    id BIGSERIAL PRIMARY KEY,
    user_id uuid not null, 
    coordinate GEOGRAPHY(POINT, 4326) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

CREATE TABLE IF NOT EXISTS tracking_service.ride_coordinate (
    id BIGSERIAL PRIMARY KEY,
    ride_id uuid not null,
    user_id uuid not null,
    coordinate GEOGRAPHY(POINT, 4326) NOT NULL,
    velocity FLOAT, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

