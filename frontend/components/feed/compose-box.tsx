"use client";

import { useState } from "react";

export default function ComposeBox() {
  const [content, setContent] = useState("");

  return (
    <div className="border border-neutral-700 p-4">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="INPUT NEW BROADCAST..."
        className="w-full bg-transparent outline-none text-sm resize-none"
      />

      <div className="flex justify-end mt-2">
        <button className="bg-green-400 text-black px-4 py-1 text-xs">
          EXECUTE_POST
        </button>
      </div>
    </div>
  );
}
