-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE public.reset_tokens (
  id BIGINT PRIMARY KEY,
  user_id public.xid NOT NULL UNIQUE,
  token VARCHAR(16) NOT NULL,
  created_at BIGINT NOT NULL,
  expires_at BIGINT NOT NULL
);
ALTER TABLE public.reset_tokens
ADD CONSTRAINT unique_reset_token_per_user UNIQUE (user_id, token);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE public.reset_tokens DROP CONSTRAINT unique_reset_token_per_user;
DROP TABLE public.reset_tokens;