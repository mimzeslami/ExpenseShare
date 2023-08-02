-- Drop the "fellow_travelers" table and its indexes
DROP INDEX IF EXISTS idx_fellow_travelers_trip_id;
DROP TABLE IF EXISTS fellow_travelers;

-- Drop the "trips" table and its indexes
DROP INDEX IF EXISTS idx_trips_user_id;
DROP TABLE IF EXISTS trips;

-- Drop the "users" table and its indexes
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
