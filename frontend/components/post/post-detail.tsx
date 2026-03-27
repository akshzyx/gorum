"use client";

import Link from "next/link";
import { useState } from "react";
import { likePost, unlikePost } from "@/lib/api";

type Props = {
  post: any;
};

export default function PostDetail({ post }: Props) {
  const [copied, setCopied] = useState(false);

  const [liked, setLiked] = useState(post.liked ?? false);
  const [likes, setLikes] = useState(post.likes ?? 0);
  const [loading, setLoading] = useState(false);

  const handleShare = async () => {
    const url = `${window.location.origin}/post/${post.id}`;
    await navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  const handleLike = async () => {
    if (loading) return;

    setLoading(true);

    const prevLiked = liked;
    const prevLikes = likes;

    if (liked) {
      setLiked(false);
      setLikes((l: number) => l - 1);
    } else {
      setLiked(true);
      setLikes((l: number) => l + 1);
    }

    try {
      if (!prevLiked) {
        await likePost(post.id);
      } else {
        await unlikePost(post.id);
      }
    } catch {
      setLiked(prevLiked);
      setLikes(prevLikes);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col gap-6 pb-6">
      <div className="flex justify-between items-center border-b border-neutral-800 pb-3">
        <div className="flex gap-3 text-xs font-mono">
          <Link
            href={`/user/${post.username || post.user_id}`}
            className="text-green-400 font-bold hover:underline"
            onClick={(e) => e.stopPropagation()}
          >
            @{post.username || post.user_id}
          </Link>

          <span className="text-neutral-500">
            [{new Date(post.created_at).toLocaleString()}]
          </span>
        </div>

        <div className="text-green-400 text-xs font-mono">PRIORITY: HIGH</div>
      </div>

      <div className="text-green-300 text-sm leading-relaxed whitespace-pre-line">
        {post.content}
      </div>

      <div className="flex gap-6 text-xs text-neutral-500 border-t border-neutral-800 pt-3 font-mono">
        <span
          onClick={handleLike}
          className={`flex items-center gap-2 cursor-pointer transition ${
            liked ? "text-green-400" : "text-neutral-500 hover:text-green-400"
          } ${loading ? "opacity-50 pointer-events-none" : ""}`}
        >
          <i
            className={`${
              liked ? "fa-solid text-green-400" : "fa-regular"
            } fa-thumbs-up`}
          />
          <span className={liked ? "text-green-400" : ""}>
            VOTE_UP [{likes}]
          </span>
        </span>

        <span
          onClick={handleShare}
          className="flex items-center gap-2 cursor-pointer hover:text-green-400"
        >
          <i className="fa-regular fa-share-nodes" />
          {copied ? "COPIED" : "PROPAGATE"}
        </span>

        <span className="flex items-center gap-2 cursor-pointer hover:text-green-400">
          <i className="fa-regular fa-flag" />
          REPORT
        </span>
      </div>
    </div>
  );
}
