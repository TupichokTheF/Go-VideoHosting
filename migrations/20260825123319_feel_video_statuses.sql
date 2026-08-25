-- +goose Up
INSERT INTO video_statuses(status) 
VALUES ('draft'), ('uploaded'), ('ready'), ('deleted');

-- +goose Down
TRUNCATE TABLE video_statuses;
