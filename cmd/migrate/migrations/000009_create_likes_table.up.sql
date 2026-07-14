CREATE TABLE IF NOT EXISTS post_likes (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, post_id)
);

CREATE TABLE IF NOT EXISTS comment_likes (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, comment_id)
);

CREATE INDEX idx_post_likes_post_id
ON post_likes(post_id);

CREATE INDEX idx_comment_likes_comment_id
ON comment_likes(comment_id);