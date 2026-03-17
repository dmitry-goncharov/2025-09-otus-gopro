CREATE DATABASE otus;
\c otus
CREATE TABLE events (
	id uuid PRIMARY KEY,
    title text NOT NULL,
    date timestamptz NOT NULL,
    user_id uuid NOT NULL
);
CREATE TABLE event_notifications (
    id uuid PRIMARY KEY,
    event_id uuid NOT NULL,
    date timestamptz NOT NULL
);