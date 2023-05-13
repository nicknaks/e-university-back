-- +goose Up
-- +goose StatementBegin
INSERT INTO users (type, name, login, password, ownerid, token) VALUES (3, 'Администратор', 'admin', 'admin', '190ec55e-3e67-49f5-87b5-41bd0ab5dd51', '08bfbdf2-85d2-4ffd-b0ea-9b0ac0be12ff');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
