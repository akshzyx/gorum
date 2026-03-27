"use client";

import { useState } from "react";
import ReplyBox from "./reply-box";

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
  postId: string;
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

function ReplyNode({
  reply,
  depth = 0,
  activeReplyId,
  setActiveReplyId,
  postId,
}: any) {
  const [collapsed, setCollapsed] = useState(false);

  const children = reply.children || [];
  const hasChildren = children.length > 0;

  const isActive = activeReplyId === reply.id;

  return (
    <div className="relative pl-6">
      {/* LINE */}
      <div
        className="absolute left-2 top-0 bottom-0 w-[2px]"
        style={{ backgroundColor: getLineColor(depth) }}
      />

      {/* NODE (CLICKABLE) */}
      <div
        onClick={() => {
          if (hasChildren) setCollapsed((prev: boolean) => !prev);
        }}
        className="group hover:bg-neutral-900/40 p-2 transition cursor-pointer"
      >
        {/* HEADER */}
        <div className="flex items-center gap-3 text-xs font-mono">
          <span className="text-green-400 w-4">
            {!hasChildren ? "[·]" : collapsed ? "[+]" : "[-]"}
          </span>

          <span className="text-green-400 font-bold">
            @{reply.username || reply.user_id}
          </span>

          <span className="text-neutral-500">
            {new Date(reply.created_at).toLocaleTimeString()}
          </span>
        </div>

        {/* CONTENT */}
        <div className="text-sm text-neutral-200 mt-2">{reply.content}</div>

        {/* ACTIONS */}
        <div
          className="flex gap-4 text-[10px] text-neutral-500 mt-2 font-mono"
          onClick={(e) => e.stopPropagation()} // 🔥 IMPORTANT
        >
          <span
            onClick={() => setActiveReplyId(isActive ? null : reply.id)}
            className="cursor-pointer hover:text-green-400"
          >
            [ REPLY ]
          </span>
        </div>

        {/* INLINE REPLY BOX */}
        {isActive && (
          <div onClick={(e) => e.stopPropagation()}>
            <ReplyBox
              postId={postId}
              parentId={reply.id}
              small
              placeholder={`REPLYING TO @${reply.username || reply.user_id}...`}
              onSuccess={() => window.location.reload()}
            />
          </div>
        )}
      </div>

      {/* CHILDREN */}
      {!collapsed && hasChildren && (
        <div className="mt-3 flex flex-col gap-3">
          {children.map((child) => (
            <ReplyNode
              key={child.id}
              reply={child}
              depth={depth + 1}
              activeReplyId={activeReplyId}
              setActiveReplyId={setActiveReplyId}
              postId={postId}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export default function ReplyList({ replies, postId }: Props) {
  const [activeReplyId, setActiveReplyId] = useState<string | null>(null);

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
        <ReplyNode
          key={r.id}
          reply={r}
          activeReplyId={activeReplyId}
          setActiveReplyId={setActiveReplyId}
          postId={postId}
        />
      ))}
    </div>
  );
}
