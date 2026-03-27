"use client";

import { useState } from "react";

export default function ReplyBox({ postId }: { postId: string }) {
  const [text, setText] = useState("");

  const send = () => {
    console.log(text);
    setText("");
  };

  return (
    <div className="fixed bottom-0 left-0 w-full border-t border-green-400 bg-[#111]">
      <div className="max-w-3xl mx-auto flex items-center gap-4 px-4 py-3">
        <span className="text-green-400 font-mono">{">"}</span>

        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="SYSTEM.EXECUTE_REPLY(CONTENT='TYPE HERE...')"
          className="flex-1 bg-transparent text-green-400 outline-none text-sm font-mono"
        />

        <button
          onClick={send}
          className="bg-green-400 text-black px-4 py-1 text-xs font-bold"
        >
          SEND_PKT
        </button>
      </div>
    </div>
  );
}
