import { useState } from "react";
import { getPosts } from "@/lib/api";

type Post = {
  id: string;
  user_id: string;
  username: string;
  content: string;
  created_at: string;
  likes: number;
  liked: boolean;
  reply_count: number;
};

type FeedResponse = {
  data: Post[];
  next_cursor: string | null;
  has_more: boolean;
};

export function useFeed() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [cursor, setCursor] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const load = async () => {
    if (loading || !hasMore) return;

    setLoading(true);

    try {
      const data: FeedResponse = await getPosts(cursor || undefined);

      setPosts((prev) => {
        const existingIds = new Set(prev.map((p) => p.id));
        const newPosts = data.data.filter((p) => !existingIds.has(p.id));
        return [...prev, ...newPosts];
      });

      setCursor(data.next_cursor);
      setHasMore(data.has_more);
    } finally {
      setLoading(false);
    }
  };

  return { posts, load, hasMore, loading };
}
