-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_board_position ON kbcolumns(board_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_board_position;
-- +goose StatementEnd
