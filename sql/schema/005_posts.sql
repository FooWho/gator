-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP,
    feed_id uuid NOT NULL,
    CONSTRAINT fk_feed_id
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE posts;