-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments(
  id bigserial PRIMARY KEY,
  user_id bigserial NOT NULL,
  post_id bigserial NOT NULL,
  content text NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comments IF EXISTS;
-- +goose StatementEnd
