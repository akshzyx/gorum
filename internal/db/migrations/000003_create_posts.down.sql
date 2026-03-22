-- Drop indexes first
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_user_id;

-- Then drop the table
DROP TABLE IF EXISTS posts;
