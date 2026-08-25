-- +goose Up
CREATE TABLE video_statuses (
	status_id SERIAL PRIMARY KEY,
	status VARCHAR(50) NOT NULL
);

CREATE TABLE videos (
	video_id SERIAL PRIMARY KEY,
	owner_id INT REFERENCES users(user_id),
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	status_id INT REFERENCES video_statuses(status_id),
	created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE video_statuses, videos CASCADE;
