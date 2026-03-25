CREATE TABLE post_likes (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, post_id),
    CHECK (user_id <> post_id)
);

CREATE INDEX idx_post_likes_post_id ON post_likes(post_id);