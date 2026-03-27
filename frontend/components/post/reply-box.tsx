"use client";

import { useState } from "react";
import { createReply } from "@/lib/api";

type Props = {
  postId: string;
  parentId?: string; // for nested reply
  placeholder?: string;
  small?: boolean;
  onSuccess?: () => void;
};

export default function ReplyBox({
  postId,
  parentId,
  placeholder,
  small = false,
  onSuccess,
}: Props) {
  const [text, setText] = useState("");
  const [loading, setLoading] = useState(false);

  const send = async () => {
    if (!text.trim()) return;

    setLoading(true);
    try {
      await createReply(parentId || postId, text);
      setText("");
      onSuccess?.();
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // 🔥 INLINE VERSION (for replies)
  if (small) {
    return (
      <div className="border-l-2 border-green-400 bg-[#151515] px-3 py-2 mt-2 flex items-center gap-3">
        <span className="text-green-400 font-mono text-xs">{">"}</span>

        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder={placeholder || "REPLY..."}
          className="flex-1 bg-transparent text-green-400 outline-none text-xs font-mono placeholder:text-neutral-600"
        />

        <button
          onClick={send}
          disabled={loading}
          className="text-green-400 border border-green-400 px-2 py-0.5 text-[10px] font-bold hover:bg-green-400 hover:text-black"
        >
          {loading ? "..." : "SEND"}
        </button>
      </div>
    );
  }

  // 🔥 MAIN VERSION (under post)
  return (
    <div className="border border-green-400 bg-[#1a1a1a] px-4 py-3">
      <div className="flex items-center gap-4">
        <span className="text-green-400 font-mono font-bold text-sm">
          {">"}
        </span>

        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="SYSTEM.EXECUTE_REPLY(CONTENT='TYPE HERE...')"
          className="flex-1 bg-transparent text-green-400 outline-none text-sm font-mono placeholder:text-neutral-600"
        />

        <button
          onClick={send}
          disabled={loading}
          className="border border-green-400 text-green-400 px-4 py-1 text-xs font-bold font-mono hover:bg-green-400 hover:text-black"
        >
          {loading ? "..." : "SEND_PKT"}
        </button>
      </div>
    </div>
  );
}
