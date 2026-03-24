CREATE TABLE follows (
    follower_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY (follower_id, following_id),

    CHECK (follower_id <> following_id)
);

-- Index for: "who follows this user"
CREATE INDEX idx_follows_following_id ON follows(following_id);

-- Index for: "who this user follows"
CREATE INDEX idx_follows_follower_id ON follows(follower_id);