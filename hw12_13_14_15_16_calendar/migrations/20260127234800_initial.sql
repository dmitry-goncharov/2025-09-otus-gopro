-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
	id uuid PRIMARY KEY,    
    title text NOT NULL,
    date timestamptz NOT NULL,
    user_id uuid NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
