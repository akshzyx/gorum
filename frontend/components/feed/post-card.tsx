"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { likePost, unlikePost } from "@/lib/api";

type Post = {
  id: string;
  username: string;
  content: string;
  created_at: string;
  reply_count: number;
  likes: number;
  liked?: boolean;
};

export default function PostCard({ post }: { post: Post }) {
  const router = useRouter();

  const [copied, setCopied] = useState(false);
  const [timeAgo, setTimeAgo] = useState("");
  const [liked, setLiked] = useState(post.liked || false);
  const [likes, setLikes] = useState(post.likes);

  useEffect(() => {
    const date = new Date(post.created_at);

    const updateTime = () => {
      const diff = Math.floor((Date.now() - date.getTime()) / 1000);

      if (diff < 60) return setTimeAgo(`${diff}s ago`);
      if (diff < 3600) return setTimeAgo(`${Math.floor(diff / 60)}m ago`);
      if (diff < 86400) return setTimeAgo(`${Math.floor(diff / 3600)}h ago`);
      return setTimeAgo(`${Math.floor(diff / 86400)}d ago`);
    };

    updateTime();

    const interval = setInterval(updateTime, 60000);

    return () => clearInterval(interval);
  }, [post.created_at]);

  const handleShare = async (e: React.MouseEvent) => {
    e.stopPropagation();

    const url = `${window.location.origin}/post/${post.id}`;
    await navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  const handleLike = async (e: React.MouseEvent) => {
    e.stopPropagation();

    const prevLiked = liked;
    const prevLikes = likes;

    // optimistic update
    setLiked(!liked);
    setLikes(liked ? likes - 1 : likes + 1);

    try {
      if (!prevLiked) {
        await likePost(post.id);
      } else {
        await unlikePost(post.id);
      }
    } catch {
      // rollback
      setLiked(prevLiked);
      setLikes(prevLikes);
    }
  };

  return (
    <div
      onClick={() => router.push(`/post/${post.id}`)}
      className="border border-neutral-700 p-4 flex flex-col gap-3 bg-black hover:bg-neutral-900 transition cursor-pointer"
    >
      <div className="flex justify-between items-start">
        <Link
          href={`/user/${post.username}`}
          onClick={(e) => e.stopPropagation()}
          className="text-green-400 text-xs font-bold hover:underline"
        >
          @{post.username}
        </Link>

        <div
          className="text-xs text-neutral-500 cursor-default"
          title={new Date(post.created_at).toLocaleString()}
        >
          {timeAgo}
        </div>
      </div>

      <div className="border-l-2 border-green-400 pl-3 text-sm leading-relaxed">
        {post.content}
      </div>

      <div className="border-t border-neutral-800" />

      <div className="flex items-center justify-between text-xs text-neutral-500">
        <div className="flex gap-5 items-center">
          <span
            onClick={(e) => {
              e.stopPropagation();
              router.push(`/post/${post.id}?reply=1`);
            }}
            className="flex items-center gap-2 cursor-pointer hover:text-green-400"
          >
            <i className="fa-regular fa-reply" />
            REPLY [{post.reply_count}]
          </span>

          <span
            onClick={handleLike}
            className={`flex items-center gap-2 cursor-pointer ${
              liked ? "text-green-400" : "hover:text-green-400"
            }`}
          >
            <i
              className={`${liked ? "fa-solid" : "fa-regular"} fa-thumbs-up`}
            />
            VOTE_UP [{likes}]
          </span>
        </div>

        <span
          onClick={handleShare}
          className="flex items-center gap-2 cursor-pointer hover:text-green-400"
        >
          <i className="fa-regular fa-share-nodes" />
          {copied ? "COPIED" : "PROPAGATE"}
        </span>
      </div>
    </div>
  );
}
