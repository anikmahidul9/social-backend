ALTER TABLE comments
DROP CONSTRAINT comments_parent_comment_fk;

ALTER TABLE comments
DROP COLUMN parent_comment_id;