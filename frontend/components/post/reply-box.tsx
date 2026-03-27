"use client";

import { useState } from "react";

export default function ReplyBox({ postId }: { postId: string }) {
  const [text, setText] = useState("");

  const send = () => {
    console.log("reply:", text);
    setText("");
  };

  return (
    <div className="fixed bottom-0 left-0 w-full bg-[#111] border-t border-green-400 p-4 flex gap-4">
      <span className="text-green-400">{">"}</span>

      <input
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="SYSTEM.EXECUTE_REPLY(...)"
        className="flex-1 bg-transparent text-green-400 outline-none text-sm"
      />

      <button
        onClick={send}
        className="bg-green-400 text-black px-4 py-1 text-xs font-bold"
      >
        SEND_PKT
      </button>
    </div>
  );
}
