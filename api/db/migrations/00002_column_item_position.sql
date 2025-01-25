-- +goose Up
-- +goose StatementBegin
ALTER TABLE column_items
ADD COLUMN position INT;

CREATE UNIQUE INDEX idx_column_position ON column_items(column_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE column_items
DROP COLUMN position;

DROP INDEX idx_column_position;
-- +goose StatementEnd
