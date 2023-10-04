-- +migrate Up
ALTER TABLE posts ADD COLUMN draft BOOLEAN NOT NULL DEFAULT TRUE;
UPDATE posts SET draft = FALSE;

-- +migrate Down
ALTER TABLE posts DROP COLUMN draft;
