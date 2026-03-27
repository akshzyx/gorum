"use client";

type Props = {
  replies: any[];
};

export default function ReplyList({ replies }: Props) {
  return (
    <div className="flex flex-col gap-8">
      <div className="text-xs text-neutral-500 border-b border-neutral-700 pb-2">
        REPLIES [{replies.length}]
      </div>

      {replies.map((r) => (
        <div key={r.id} className="flex">
          <div className="w-[2px] bg-green-400 mr-4" />

          <div className="flex flex-col gap-2">
            <div className="flex gap-3 text-xs">
              <span className="text-green-400">@{r.username || r.user_id}</span>
              <span className="text-neutral-500">
                {new Date(r.created_at).toLocaleTimeString()}
              </span>
            </div>

            <div className="text-sm text-neutral-200">{r.content}</div>
          </div>
        </div>
      ))}
    </div>
  );
}
