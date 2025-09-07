-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_location(
    login		TEXT,
	session_id	INTEGER NOT NULL,
	mountpoint 	TEXT NOT NULL,
	station		TEXT NOT NULL,
	ntrip_agent	TEXT NOT NULL,
	connect_time	INTEGER NOT NULL,
	time_span	INTEGER,
	recieved_data	DOUBLE PRECISION NOT NULL,
	sent_data	DOUBLE PRECISION NOT NULL,
	status_code	INTEGER,
	latency	INTEGER,
	sv_num	INTEGER,
	lat	DOUBLE PRECISION NOT NULL,
	lon	DOUBLE PRECISION NOT NULL,
	height	DOUBLE PRECISION,
	station_distance	DOUBLE PRECISION NOT NULL,
	created_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_login ON user_location(Login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_location;
-- +goose StatementEnd
