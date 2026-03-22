ALTER TABLE posts
ADD COLUMN parent_post_id TEXT REFERENCES posts(id),
ADD COLUMN root_post_id TEXT REFERENCES posts(id);

CREATE INDEX idx_posts_parent_post_id ON posts(parent_post_id);
CREATE INDEX idx_posts_root_post_id ON posts(root_post_id);
