"use client";

type Reply = {
  id: string;
  user_id: string;
  username?: string;
  content: string;
  created_at: string;
  parent_post_id?: string | null;
};

type Props = {
  replies: Reply[];
};

// TREE BUILDER

function buildTree(replies: Reply[], rootID: string) {
  const map = new Map<string, any>();
  const roots: any[] = [];

  replies.forEach((r) => {
    map.set(r.id, { ...r, children: [] });
  });

  replies.forEach((r) => {
    if (r.parent_post_id === rootID) {
      roots.push(map.get(r.id));
    } else if (r.parent_post_id && map.has(r.parent_post_id)) {
      map.get(r.parent_post_id).children.push(map.get(r.id));
    }
  });

  return roots;
}
//  RECURSIVE ITEM

function ReplyItem({ reply, depth = 0 }: { reply: any; depth?: number }) {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex">
        {/* LEFT LINE + INDENT */}
        <div className="mr-4 flex" style={{ marginLeft: depth * 20 }}>
          <div className="w-[2px] bg-green-400" />
        </div>

        {/* CONTENT */}
        <div className="flex flex-col gap-2">
          <div className="flex gap-3 text-xs font-mono">
            <span className="text-green-400">
              @{reply.username || reply.user_id}
            </span>

            <span className="text-neutral-500">
              {new Date(reply.created_at).toLocaleTimeString()}
            </span>
          </div>

          <div className="text-sm text-neutral-200 leading-relaxed whitespace-pre-line">
            {reply.content}
          </div>
        </div>
      </div>

      {/* CHILDREN */}
      {reply.children?.length > 0 && (
        <div className="flex flex-col gap-4">
          {reply.children.map((child: any) => (
            <ReplyItem key={child.id} reply={child} depth={depth + 1} />
          ))}
        </div>
      )}
    </div>
  );
}

// MAIN COMPONENT

export default function ReplyList({ replies }: Props) {
  if (replies.length === 0) return null;

  const rootID = replies[0]?.parent_post_id; // main post id
  const tree = buildTree(replies, rootID);

  return (
    <div className="flex flex-col gap-8">
      <div className="text-xs text-neutral-500 border-b border-neutral-700 pb-2 font-mono">
        REPLIES [{replies.length}]
      </div>

      {tree.map((r) => (
        <ReplyItem key={r.id} reply={r} />
      ))}
    </div>
  );
}
