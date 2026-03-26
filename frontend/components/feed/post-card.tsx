"use client";

export default function PostCard({ post }: { post: any }) {
  return (
    <div className="border border-neutral-700 p-4 flex flex-col gap-3">
      <div className="text-green-400 text-xs font-bold">@{post.username}</div>

      <div className="border-l-2 border-green-400 pl-3">{post.content}</div>

      <div className="flex gap-4 text-xs text-neutral-500">
        <span>REPLY [{post.reply_count}]</span>
        <span>LIKE [{post.likes}]</span>
      </div>
    </div>
  );
}
