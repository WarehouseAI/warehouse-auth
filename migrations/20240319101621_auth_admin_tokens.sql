-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE public.admin_tokens (
  user_id public.xid NOT NULL DEFAULT xid(),
  number BIGINT NOT NULL,
  purpose INTEGER NOT NULL,
  secret CHAR(64) NOT NULL,
  expires_at BIGINT NOT NULL
);
ALTER TABLE public.admin_tokens
ADD CONSTRAINT admin_tokens_pkey PRIMARY KEY (user_id, number, purpose);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE public.admin_tokens;