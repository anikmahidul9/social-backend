ALTER TABLE comments
ADD COLUMN parent_comment_id BIGINT NULL;

ALTER TABLE comments
ADD CONSTRAINT comments_parent_comment_fk
FOREIGN KEY (parent_comment_id)
REFERENCES comments(id)
ON DELETE CASCADE;