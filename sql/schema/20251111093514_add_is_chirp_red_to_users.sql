-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD is_chirpy_red boolean DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN is_chirpy_red;
-- +goose StatementEnd
