"use client";

import { useEffect, useRef } from "react";
import { useFeed } from "@/hooks/use-feed";
import PostCard from "./post-card";
import Loader from "./loader";

export default function FeedList() {
  const { posts, load, hasMore, loading } = useFeed();
  const ref = useRef<HTMLDivElement | null>(null);

  // initial load
  useEffect(() => {
    load();
  }, []);

  // intersection observer
  useEffect(() => {
    if (!ref.current) return;

    const el = ref.current;

    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];

        if (entry.isIntersecting && hasMore && !loading) {
          load();
        }
      },
      {
        rootMargin: "200px", // preload earlier
      },
    );

    observer.observe(el);

    return () => {
      observer.unobserve(el);
      observer.disconnect();
    };
  }, [hasMore, loading, load]);

  return (
    <>
      {posts.map((p) => (
        <PostCard key={p.id} post={p} />
      ))}

      {hasMore && <div ref={ref} className="h-10" />}

      {loading && <Loader />}
    </>
  );
}
