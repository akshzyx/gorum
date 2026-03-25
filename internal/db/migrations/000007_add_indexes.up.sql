-- Feed
CREATE INDEX IF NOT EXISTS idx_posts_feed 
ON posts(created_at DESC)
WHERE deleted_at IS NULL AND parent_post_id IS NULL;

-- User posts
CREATE INDEX IF NOT EXISTS idx_posts_user_feed
ON posts(user_id, created_at DESC)
WHERE parent_post_id IS NULL AND deleted_at IS NULL;

-- Replies
CREATE INDEX IF NOT EXISTS idx_posts_parent_created
ON posts(parent_post_id, created_at ASC)
WHERE deleted_at IS NULL;

-- Thread
CREATE INDEX IF NOT EXISTS idx_posts_root_created
ON posts(root_post_id, created_at ASC)
WHERE deleted_at IS NULL;

-- Likes
CREATE INDEX IF NOT EXISTS idx_post_likes_post_id ON post_likes(post_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_post_likes_user_post ON post_likes(user_id, post_id);