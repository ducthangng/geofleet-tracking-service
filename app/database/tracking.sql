CREATE DATABASE ride_tracking_service;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE SCHEMA IF NOT EXISTS tracking_service;

CREATE TABLE IF NOT EXISTS tracking_service.coordinate (
    id BIGSERIAL,
    user_id uuid not null,
    coordinate GEOGRAPHY(POINT, 4326) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE TABLE IF NOT EXISTS tracking_service.ride_coordinate (
    id BIGSERIAL,
    ride_id uuid not null,
    user_id uuid not null,
    coordinate GEOGRAPHY(POINT, 4326) NOT NULL,
    velocity FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE TABLE tracking_service.ride_coordinate_partition_2026_01_21 PARTITION OF tracking_service.ride_coordinate
FOR VALUES FROM ('2026-01-21 00:00:00') TO ('2026-01-21 23:59:59');

-- children partition
CREATE TABLE tracking_service.coordinate_partition_2026_01_21 PARTITION OF tracking_service.coordinate
FOR VALUES FROM ('2026-01-21 00:00:00') TO ('2026-01-21 23:59:59');



