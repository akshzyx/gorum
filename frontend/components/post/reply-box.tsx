"use client";

import { useEffect, useRef, useState } from "react";

export default function ReplyBox({
  postId,
  parentId,
  placeholder,
  small = false,
  onSuccess,
  autoFocus = false,
}: any) {
  const [text, setText] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  // AUTO FOCUS
  useEffect(() => {
    if (autoFocus && inputRef.current) {
      inputRef.current.focus();

      // scroll into view nicely
      inputRef.current.scrollIntoView({
        behavior: "smooth",
        block: "center",
      });
    }
  }, [autoFocus]);

  const send = async () => {
    if (!text.trim()) return;
    console.log(text);
    setText("");
    onSuccess?.();
  };

  return (
    <div className="border border-green-400 bg-[#1a1a1a] px-4 py-3">
      <div className="flex items-center gap-4">
        <span className="text-green-400 font-mono">{">"}</span>

        <input
          ref={inputRef}
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder={
            placeholder || "SYSTEM.EXECUTE_REPLY(CONTENT='TYPE HERE...')"
          }
          className="flex-1 bg-transparent text-green-400 outline-none text-sm font-mono"
        />

        <button
          onClick={send}
          className="border border-green-400 px-4 py-1 text-xs font-bold"
        >
          SEND_PKT
        </button>
      </div>
    </div>
  );
}
