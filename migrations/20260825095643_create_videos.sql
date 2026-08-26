-- +goose Up
CREATE TABLE videos (
	video_id UUID PRIMARY KEY,
	owner_id INT REFERENCES users(user_id),
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'draft',
	size BIGINT,
	created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
	CONSTRAINT videos_status_check CHECK (status IN ('draft','uploaded','processing','ready','failed','deleted'))
);

-- +goose Down
DROP TABLE video_statuses, videos CASCADE;
