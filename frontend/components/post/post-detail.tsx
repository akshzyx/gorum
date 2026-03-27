"use client";

type Props = {
  post: any;
};

export default function PostDetail({ post }: Props) {
  return (
    <div className="flex flex-col gap-6 pb-6">
      {/* TOP BAR */}
      <div className="flex justify-between items-center border-b border-neutral-800 pb-3">
        <div className="flex gap-3 text-xs font-mono">
          <span className="text-green-400 font-bold">
            @{post.username || post.user_id}
          </span>

          <span className="text-neutral-500">
            [{new Date(post.created_at).toLocaleString()}]
          </span>
        </div>

        <div className="text-green-400 text-xs font-mono">PRIORITY: HIGH</div>
      </div>

      {/* CONTENT ONLY */}
      <div className="text-green-300 text-sm leading-relaxed whitespace-pre-line">
        {post.content}
      </div>

      {/* ACTION BAR */}
      <div className="flex gap-6 text-xs text-neutral-500 border-t border-neutral-800 pt-3 font-mono">
        <span className="flex items-center gap-2 cursor-pointer hover:text-green-400">
          <i className="fa-regular fa-thumbs-up"></i>
          {post.likes || 0}_UP
        </span>

        <span className="flex items-center gap-2 cursor-pointer hover:text-green-400">
          <i className="fa-regular fa-share-nodes"></i>
          EXPORT
        </span>

        <span className="flex items-center gap-2 cursor-pointer hover:text-green-400">
          <i className="fa-regular fa-flag"></i>
          REPORT
        </span>
      </div>
    </div>
  );
}
