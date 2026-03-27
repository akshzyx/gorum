"use client";

type Props = {
  replies: any[];
};

export default function ReplyList({ replies }: Props) {
  return (
    <div className="flex flex-col gap-8">
      {/* HEADER */}
      <div className="text-xs text-neutral-500 border-b border-neutral-700 pb-2 font-mono">
        REPLIES [{replies.length}]
      </div>

      {replies.map((r) => (
        <div key={r.id} className="flex">
          {/* LEFT GREEN LINE */}
          <div className="w-[2px] bg-green-400 mr-4" />

          {/* CONTENT */}
          <div className="flex flex-col gap-2">
            <div className="flex gap-3 text-xs font-mono">
              <span className="text-green-400">@{r.username || r.user_id}</span>

              <span className="text-neutral-500">
                {new Date(r.created_at).toLocaleTimeString()}
              </span>
            </div>

            <div className="text-sm text-neutral-200 leading-relaxed">
              {r.content}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
