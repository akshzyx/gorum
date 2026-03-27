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
  const [open, setOpen] = useState(false);

  const [post, setPost] = useState<any>(null);
  const [replies, setReplies] = useState<any[]>([]);

  useEffect(() => {
    if (!id) return;

    getThread(id as string).then((thread) => {
      setPost(thread[0]);
      setReplies(thread.slice(1));
    });
  }, [id]);

  if (!post) return null;

  return (
    <div className="min-h-screen bg-[#0e0e0e] text-white">
      <Navbar onMenuClick={() => setOpen(true)} />
      <Sidebar open={open} setOpen={setOpen} />

      <main className="pt-16 md:ml-[180px] flex justify-center pb-24">
        <div className="w-full max-w-2xl flex flex-col gap-6 px-4">
          <div className="text-[10px] text-neutral-500 font-mono border-b border-neutral-800 pb-2">
            LOCAL_FS / THREAD_ID: {id} /{" "}
            <span className="text-green-400">LIVE_VIEW</span>
          </div>

          <PostDetail post={post} />

          <ReplyBox
            postId={id as string}
            onSuccess={() => window.location.reload()}
          />

          <ReplyList replies={replies} postId={id as string} />
        </div>
      </main>
    </div>
  );
}
