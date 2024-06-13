-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE public.verification_tokens (
  id BIGINT PRIMARY KEY,
  user_id public.xid NOT NULL UNIQUE,
  token VARCHAR(16) NOT NULL,
  send_to VARCHAR(255) NOT NULL UNIQUE,
  created_at BIGINT NOT NULL,
  expires_at BIGINT NOT NULL
);
ALTER TABLE public.verification_tokens
ADD CONSTRAINT unique_verification_token_and_user UNIQUE (user_id, send_to);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE public.verification_tokens DROP CONSTRAINT unique_verification_token_and_user;
DROP TABLE public.verification_tokens;