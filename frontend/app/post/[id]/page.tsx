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

      <main className="pt-16 md:ml-[180px] flex justify-center pb-24">
        <div className="w-full max-w-2xl flex flex-col gap-6 px-4">
          {/* HEADER */}
          <div className="text-[10px] text-neutral-500 font-mono border-b border-neutral-800 pb-2 tracking-wider">
            LOCAL_FS / THREAD_ID: {id} /{" "}
            <span className="text-green-400">LIVE_VIEW</span>
          </div>

          {/* POST */}
          <PostDetail post={post} />

          {/* REPLY BOX */}
          <ReplyBox
            postId={id as string}
            onSuccess={() => window.location.reload()}
          />

          {/* REPLIES */}
          <ReplyList replies={replies} postId={id as string} />
        </div>
      </main>
    </div>
  );
}
