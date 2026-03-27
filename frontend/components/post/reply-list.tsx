"use client";

import { useState, useEffect, useRef } from "react";
import ReplyBox from "./reply-box";
import Link from "next/link";

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
  targetId?: string | null;
};

const VISIBLE_COUNT = 2;

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

// check if subtree contains target
function containsTarget(node: Reply, targetId: string | null): boolean {
  if (!targetId) return false;
  if (node.id === targetId) return true;
  return (node.children || []).some((child) => containsTarget(child, targetId));
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
  targetId,
}: any) {
  const hasTarget = containsTarget(reply, targetId);

  const [collapsed, setCollapsed] = useState(false);
  const [expanded, setExpanded] = useState(false);

  const ref = useRef<HTMLDivElement>(null);
  const hasScrolled = useRef(false);

  const children = reply.children || [];
  const hasChildren = children.length > 0;

  const isActive = activeReplyId === reply.id;

  // AUTO OPEN PATH TO TARGET
  useEffect(() => {
    if (hasTarget) {
      setCollapsed(false);
      setExpanded(true);
    }
  }, [hasTarget]);

  // SCROLL ONLY ONCE
  useEffect(() => {
    if (targetId === reply.id && ref.current && !hasScrolled.current) {
      hasScrolled.current = true;

      ref.current.scrollIntoView({
        behavior: "smooth",
        block: "center",
      });
    }
  }, [targetId, reply.id]);

  // SHOW TARGET EVEN IF HIDDEN
  let visibleChildren = expanded ? children : children.slice(0, VISIBLE_COUNT);

  if (!expanded && targetId) {
    const targetChild = children.find((c) => containsTarget(c, targetId));

    if (targetChild && !visibleChildren.find((c) => c.id === targetChild.id)) {
      visibleChildren = [...visibleChildren, targetChild];
    }
  }

  const hiddenCount = children.length - visibleChildren.length;

  const handleCopy = (e: any) => {
    e.stopPropagation();
    const url = `${window.location.origin}/post/${reply.id}`;
    navigator.clipboard.writeText(url);
  };

  return (
    <div className="relative pl-6">
      <div
        className="absolute left-2 top-0 bottom-0 w-[2px]"
        style={{ backgroundColor: getLineColor(depth) }}
      />

      <div
        ref={ref}
        onClick={() => {
          if (hasChildren) setCollapsed((p) => !p);
        }}
        className={`group p-2 transition cursor-pointer
          ${
            targetId === reply.id
              ? "bg-green-400/10 border border-green-400 shadow-[0_0_8px_rgba(74,222,128,0.2)]"
              : "hover:bg-neutral-900/40"
          }
        `}
      >
        <div className="flex items-center gap-3 text-xs font-mono">
          <span className="text-green-400 w-4">
            {!hasChildren ? "[·]" : collapsed ? "[+]" : "[-]"}
          </span>

          <Link
            href={`/user/${reply.username || reply.user_id}`}
            onClick={(e) => e.stopPropagation()}
            className="text-green-400 font-bold hover:underline"
          >
            @{reply.username || reply.user_id}
          </Link>

          <span className="text-neutral-500">
            {new Date(reply.created_at).toLocaleTimeString()}
          </span>
        </div>

        <div className="text-sm text-neutral-200 mt-2">{reply.content}</div>

        <div
          className="flex gap-4 text-[10px] text-neutral-500 mt-2 font-mono"
          onClick={(e) => e.stopPropagation()}
        >
          <span
            onClick={() => setActiveReplyId(isActive ? null : reply.id)}
            className="cursor-pointer hover:text-green-400"
          >
            [ REPLY ]
          </span>

          {/* <span
            onClick={handleCopy}
            className="cursor-pointer hover:text-green-400"
          >
            [ PROPAGATE ]
          </span> */}
        </div>

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

      {!collapsed && hasChildren && (
        <div className="mt-3 flex flex-col gap-3">
          {visibleChildren.map((child) => (
            <ReplyNode
              key={child.id}
              reply={child}
              depth={depth + 1}
              activeReplyId={activeReplyId}
              setActiveReplyId={setActiveReplyId}
              postId={postId}
              targetId={targetId}
            />
          ))}

          {!expanded && hiddenCount > 0 && (
            <div
              onClick={(e) => {
                e.stopPropagation();
                setExpanded(true);
              }}
              className="font-mono text-[10px] text-neutral-400 cursor-pointer hover:text-green-400"
            >
              [+] EXPAND ({hiddenCount} HIDDEN_NODES)
            </div>
          )}

          {expanded && hiddenCount > 0 && (
            <div
              onClick={(e) => {
                e.stopPropagation();
                setExpanded(false);
              }}
              className="font-mono text-[10px] text-neutral-400 cursor-pointer hover:text-green-400"
            >
              [-] COLLAPSE
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default function ReplyList({ replies, postId, targetId }: Props) {
  const [activeReplyId, setActiveReplyId] = useState<string | null>(null);

  if (!replies.length) return null;

  const rootID = replies[0]?.parent_post_id;
  const tree = buildTree(replies, rootID);

  return (
    <div className="flex flex-col gap-6">
      <div className="text-xs text-neutral-500 border-b border-neutral-700 pb-2 font-mono flex justify-between">
        <span>REPLIES [{replies.length}]</span>
        <span className="text-green-400/50">THREAD_DENSITY: NOMINAL</span>
      </div>

      {tree.map((r) => (
        <ReplyNode
          key={r.id}
          reply={r}
          activeReplyId={activeReplyId}
          setActiveReplyId={setActiveReplyId}
          postId={postId}
          targetId={targetId}
        />
      ))}
    </div>
  );
}
