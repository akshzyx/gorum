"use client";

type Props = {
  post: any;
};

export default function PostDetail({ post }: Props) {
  return (
    <div className="flex flex-col gap-6 border-b border-neutral-700 pb-6">
      <div className="flex justify-between items-center border-b border-neutral-800 pb-2">
        <div className="flex gap-3 text-xs">
          <span className="text-green-400 font-bold">
            @{post.username || post.user_id} {/* ✅ fallback */}
          </span>
          <span className="text-neutral-500">
            [{new Date(post.created_at).toLocaleString()}]
          </span>
        </div>

        <div className="text-green-400 text-xs">PRIORITY: HIGH</div>
      </div>

      <div className="text-2xl font-bold uppercase tracking-wide">
        {post.content}
      </div>

      <div className="text-green-300 text-sm leading-relaxed">
        {post.content}
      </div>

      <div className="flex gap-6 text-xs text-neutral-500 border-t border-neutral-800 pt-3">
        <span>👍 {post.likes || 0}</span>
        <span>↗ EXPORT</span>
        <span>⚑ REPORT</span>
      </div>
    </div>
  );
}
