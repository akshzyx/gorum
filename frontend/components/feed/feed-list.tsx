"use client";

import { useEffect, useRef } from "react";
import { useFeed } from "@/hooks/use-feed";
import PostCard from "./post-card";
import Loader from "./loader";

export default function FeedList() {
  const { posts, load, hasMore, loading } = useFeed();
  const ref = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    load();
  }, []);

  useEffect(() => {
    if (!ref.current) return;

    const observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        load();
      }
    });

    const el = ref.current;
    observer.observe(el);

    return () => observer.unobserve(el);
  }, [hasMore]);

  return (
    <>
      {posts.map((p) => (
        <PostCard key={p.id} post={p} />
      ))}

      {hasMore && <div ref={ref} />}

      {loading && <Loader />}
    </>
  );
}
