-- +goose Up
-- +goose StatementBegin
create table room (
    uuid uuid primary key,
    name text,
    dt_create timestamp without time zone default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd