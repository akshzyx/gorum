-- Drop indexes first
DROP INDEX IF EXISTS idx_posts_parent_post_id;
DROP INDEX IF EXISTS idx_posts_root_post_id;

-- Then drop columns
ALTER TABLE posts
DROP COLUMN IF EXISTS parent_post_id,
DROP COLUMN IF EXISTS root_post_id;
