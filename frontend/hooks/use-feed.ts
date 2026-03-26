import { useState } from "react";
import { getPosts } from "@/lib/api";

export function useFeed() {
  const [posts, setPosts] = useState<any[]>([]);
  const [cursor, setCursor] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const load = async () => {
    if (loading || !hasMore) return;

    setLoading(true);

    const data = await getPosts(cursor || undefined);

    setPosts((prev) => [...prev, ...data.data]);
    setCursor(data.next_cursor);
    setHasMore(data.has_more);

    setLoading(false);
  };

  return { posts, load, hasMore, loading };
}
