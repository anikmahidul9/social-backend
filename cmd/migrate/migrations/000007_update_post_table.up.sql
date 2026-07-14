ALTER TABLE posts
DROP COLUMN IF EXISTS tags;

CREATE TYPE post_visibility AS ENUM ('public', 'private');

ALTER TABLE posts
ADD COLUMN visibility post_visibility NOT NULL DEFAULT 'public';