-- +migrate Up
PRAGMA foreign_keys = ON;
CREATE TEMPORARY TABLE post_tags_backup AS SELECT * FROM post_tags;
DROP TABLE post_tags;
CREATE TABLE post_tags (
    post_id TEXT NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (post_id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag)
);
INSERT INTO post_tags (post_id, tag) SELECT post_id, tag FROM post_tags_backup;
DROP TABLE post_tags_backup;

-- +migrate Down

