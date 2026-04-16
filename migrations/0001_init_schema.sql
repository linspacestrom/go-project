-- +goose Up
-- +goose StatementBegin

CREATE TABLE user_roles (
                            id   SERIAL PRIMARY KEY,
                            name TEXT NOT NULL UNIQUE
);

INSERT INTO user_roles (name) VALUES ('admin'), ('user'), ('mentor');

CREATE TABLE users (
                       id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email         TEXT NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       role_id       INT NOT NULL REFERENCES user_roles(id),
                       created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd