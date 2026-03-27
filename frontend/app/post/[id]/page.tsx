"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

import { getThread } from "@/lib/api";

import Navbar from "@/components/navbar";
import Sidebar from "@/components/sidebar";

import PostDetail from "@/components/post/post-detail";
import ReplyList from "@/components/post/reply-list";
import ReplyBox from "@/components/post/reply-box";

export default function Page() {
  const { id } = useParams();

  const [post, setPost] = useState<any>(null);
  const [replies, setReplies] = useState<any[]>([]);

  useEffect(() => {
    if (!id) return;

    const load = async () => {
      const thread = await getThread(id as string);

      // first item = main post
      setPost(thread[0]);

      // rest = replies
      setReplies(thread.slice(1));
    };

    load();
  }, [id]);

  if (!post) return null;

  return (
    <div className="min-h-screen bg-[#0e0e0e] text-white">
      <Navbar />
      <Sidebar />

      <main className="pt-20 md:pl-64 flex justify-center pb-32">
        <div className="w-full max-w-3xl flex flex-col gap-10 px-4">
          <div className="text-xs text-neutral-500 font-mono">
            LOCAL_FS / THREAD_ID: {id} /{" "}
            <span className="text-green-400">LIVE_VIEW</span>
          </div>

          <PostDetail post={post} />

          <ReplyList replies={replies} />
        </div>
      </main>

      <ReplyBox postId={id as string} />
    </div>
  );
}
