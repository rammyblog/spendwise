-- +goose Up
-- +goose StatementBegin
ALTER TABLE expenses ADD COLUMN name VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE expenses DROP COLUMN name;
-- +goose StatementEnd
