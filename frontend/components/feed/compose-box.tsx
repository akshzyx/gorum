"use client";

import { useState } from "react";
import { createPost } from "@/lib/api";

export default function ComposeBox({ onPost }: any) {
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(false);

  const submit = async () => {
    if (!content.trim() || loading) return;

    setLoading(true);

    try {
      await createPost(content);
      setContent("");
      onPost?.(); // refresh feed
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="border border-neutral-700 p-4">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="INPUT NEW BROADCAST..."
        className="w-full bg-transparent outline-none text-sm resize-none"
      />

      <div className="flex justify-end mt-2">
        <button
          onClick={submit}
          className="bg-green-400 text-black px-4 py-1 text-xs"
        >
          {loading ? "EXECUTING..." : "EXECUTE_POST"}
        </button>
      </div>
    </div>
  );
}
