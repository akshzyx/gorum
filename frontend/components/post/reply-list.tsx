"use client";

import { useState } from "react";

type Reply = {
  id: string;
  user_id: string;
  username?: string;
  content: string;
  created_at: string;
  parent_post_id?: string;
  children?: Reply[];
};

type Props = {
  replies: Reply[];
};

function buildTree(replies: Reply[], rootID: string) {
  const map = new Map<string, Reply>();
  const roots: Reply[] = [];

  replies.forEach((r) => {
    map.set(r.id, { ...r, children: [] });
  });

  replies.forEach((r) => {
    if (r.parent_post_id === rootID) {
      roots.push(map.get(r.id)!);
    } else if (r.parent_post_id && map.has(r.parent_post_id)) {
      map.get(r.parent_post_id)!.children!.push(map.get(r.id)!);
    }
  });

  return roots;
}

function getLineColor(depth: number) {
  const opacity = Math.max(1 - depth * 0.25, 0.2);
  return `rgba(74, 222, 128, ${opacity})`;
}

function ReplyNode({ reply, depth = 0 }: { reply: Reply; depth?: number }) {
  const [expanded, setExpanded] = useState(false);

  const children = reply.children || [];
  const hasChildren = children.length > 0;

  const VISIBLE_COUNT = 2;

  const visibleChildren = expanded
    ? children
    : children.slice(0, VISIBLE_COUNT);

  const hiddenCount = children.length - VISIBLE_COUNT;

  return (
    <div className="relative pl-6">
      {/* CONTINUOUS LINE */}
      <div
        className="absolute left-2 top-0 bottom-0 w-[2px]"
        style={{ backgroundColor: getLineColor(depth) }}
      />

      {/* NODE */}
      <div className="group hover:bg-neutral-900/40 p-2 transition cursor-pointer">
        {/* HEADER */}
        <div className="flex items-center gap-3 text-xs font-mono">
          <span className="text-green-400 w-4">
            {!hasChildren ? "[·]" : expanded ? "[-]" : "[+]"}
          </span>

          <span className="text-green-400 font-bold">
            @{reply.username || reply.user_id}
          </span>

          <span className="text-neutral-500">
            {new Date(reply.created_at).toLocaleTimeString()}
          </span>
        </div>

        {/* CONTENT */}
        <div className="text-sm text-neutral-200 mt-2 leading-relaxed">
          {reply.content}
        </div>

        {/* ACTIONS */}
        <div className="flex gap-4 text-[10px] text-neutral-500 mt-2 font-mono">
          <span className="cursor-pointer hover:text-green-400">[ REPLY ]</span>
          <span className="cursor-pointer hover:text-green-400">
            [ VOTE_UP ]
          </span>
        </div>
      </div>

      {/* CHILDREN */}
      {hasChildren && (
        <div className="mt-4 flex flex-col gap-4">
          {visibleChildren.map((child) => (
            <ReplyNode key={child.id} reply={child} depth={depth + 1} />
          ))}

          {/* EXPAND / COLLAPSE */}
          {!expanded && hiddenCount > 0 && (
            <div
              className="font-mono text-[10px] text-neutral-400 cursor-pointer hover:text-green-400"
              onClick={() => setExpanded(true)}
            >
              [+] EXPAND ({hiddenCount} HIDDEN_NODES)
            </div>
          )}

          {expanded && hiddenCount > 0 && (
            <div
              className="font-mono text-[10px] text-neutral-400 cursor-pointer hover:text-green-400"
              onClick={() => setExpanded(false)}
            >
              [-] COLLAPSE
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default function ReplyList({ replies }: Props) {
  if (!replies.length) return null;

  const rootID = replies[0]?.parent_post_id;
  const tree = buildTree(replies, rootID);

  return (
    <div className="flex flex-col gap-6">
      {/* HEADER */}
      <div className="text-xs text-neutral-500 border-b border-neutral-700 pb-2 font-mono flex justify-between">
        <span>REPLIES [{replies.length}]</span>
        <span className="text-green-400/50">THREAD_DENSITY: NOMINAL</span>
      </div>

      {/* TREE */}
      {tree.map((r) => (
        <ReplyNode key={r.id} reply={r} />
      ))}
    </div>
  );
}
