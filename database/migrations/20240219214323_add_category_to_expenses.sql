-- +goose Up
-- +goose StatementBegin
ALTER TABLE expenses
ADD COLUMN category_id UUID,
ADD CONSTRAINT fk_expenses_categories
    FOREIGN KEY (category_id) REFERENCES categories(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE expenses
DROP FOREIGN KEY fk_expenses_categories,
DROP COLUMN category_id;
-- +goose StatementEnd
