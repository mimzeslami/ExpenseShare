
-- Create uuid extension 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the "users" table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(100) NOT NULL
);

-- Create an index on the "email" column for faster user lookup
CREATE UNIQUE INDEX idx_users_email ON users (email);

-- Create the "trips" table
CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_name VARCHAR(100) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users(id)
);

-- Create an index on the "user_id" column for faster trip retrieval
CREATE INDEX idx_trips_user_id ON trips (user_id);

-- Create the "fellow_travelers" table
CREATE TABLE fellow_travelers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID REFERENCES trips(id),
    fellow_first_name VARCHAR(50) NOT NULL,
    fellow_last_name VARCHAR(50) NOT NULL,
    fellow_email VARCHAR(100) NOT NULL
);

-- Create an index on the "trip_id" column for faster fellow traveler retrieval
CREATE INDEX idx_fellow_travelers_trip_id ON fellow_travelers (trip_id);
