"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { getPostById, getReplies } from "@/lib/api";
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
      const p = await getPostById(id as string);
      const r = await getReplies(id as string);

      setPost(p);
      setReplies(r.data); // ✅ FIX
    };

    load();
  }, [id]);

  if (!post) return null;

  return (
    <div className="pt-20 md:pl-64 flex justify-center">
      <div className="w-full max-w-3xl flex flex-col gap-10 px-4">
        <div className="text-xs text-neutral-500 font-mono">
          LOCAL_FS / THREAD_ID: {id} / LIVE_VIEW
        </div>

        <PostDetail post={post} />

        <ReplyList replies={replies} />
      </div>

      <ReplyBox postId={id as string} />
    </div>
  );
}
