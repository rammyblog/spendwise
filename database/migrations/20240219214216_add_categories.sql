-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE categories (
    id UUID DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY (id)

);
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
