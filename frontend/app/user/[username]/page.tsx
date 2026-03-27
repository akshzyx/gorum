"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

import Navbar from "@/components/navbar";
import Sidebar from "@/components/sidebar";
import PostCard from "@/components/feed/post-card";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

// safely extract array from any API shape
function extractArray(data: any): any[] {
  if (Array.isArray(data)) return data;
  if (Array.isArray(data?.data)) return data.data;
  if (Array.isArray(data?.posts)) return data.posts;
  if (Array.isArray(data?.replies)) return data.replies;
  return [];
}

// safe fetch wrapper
async function safeFetch(url: string) {
  try {
    const res = await fetch(url);

    if (!res.ok) {
      console.error("API ERROR:", url);
      return null;
    }

    return await res.json();
  } catch (err) {
    console.error("FETCH FAILED:", url, err);
    return null;
  }
}

export default function ProfilePage() {
  const { username } = useParams();

  const [open, setOpen] = useState(false);

  const [user, setUser] = useState<any>(null);
  const [posts, setPosts] = useState<any[]>([]);
  const [replies, setReplies] = useState<any[]>([]);
  const [tab, setTab] = useState<"posts" | "replies">("posts");

  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!username) return;

    const load = async () => {
      setLoading(true);

      const userData = await safeFetch(`${BASE_URL}/user/${username}`);

      const postsData = await safeFetch(`${BASE_URL}/user/${username}/posts`);

      const repliesData = await safeFetch(
        `${BASE_URL}/user/${username}/replies`,
      );

      setUser(userData || null);
      setPosts(extractArray(postsData));
      setReplies(extractArray(repliesData));

      setLoading(false);
    };

    load();
  }, [username]);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#0e0e0e] text-green-400 flex items-center justify-center font-mono text-sm">
        LOADING_PROFILE...
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-[#0e0e0e] text-red-400 flex items-center justify-center font-mono text-sm">
        USER_NOT_FOUND
      </div>
    );
  }

  const created = new Date(user.created_at || Date.now());
  const uptimeDays = Math.floor(
    (Date.now() - created.getTime()) / (1000 * 60 * 60 * 24),
  );

  const activeList = tab === "posts" ? posts : replies;

  return (
    <div className="min-h-screen bg-[#0e0e0e] text-white">
      <Navbar onMenuClick={() => setOpen(true)} />
      <Sidebar open={open} setOpen={setOpen} />

      <main className="pt-16 md:ml-[180px] flex justify-center">
        <div className="w-full max-w-3xl flex flex-col px-4 py-8 gap-8">
          {/* PROFILE HEADER */}
          <div className="border border-neutral-700 bg-[#111] p-6">
            <div className="flex flex-col md:flex-row gap-6">
              {/* AVATAR */}
              <div className="w-24 h-24 border border-green-400">
                {user.avatar_url ? (
                  <img
                    src={user.avatar_url}
                    alt="avatar"
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-green-400 font-bold">
                    {user.username?.[0]?.toUpperCase() || "U"}
                  </div>
                )}
              </div>

              {/* INFO */}
              <div className="flex-1 flex flex-col gap-4">
                {/* USERNAME */}
                <h1 className="text-2xl md:text-3xl font-bold text-green-400 font-mono uppercase">
                  {user.username || "UNKNOWN_USER"}
                </h1>

                {/* BIO */}
                <div className="text-sm text-neutral-300 border-l-2 border-neutral-700 pl-3">
                  {user.bio || "NO_DESCRIPTION_PROVIDED"}
                </div>

                {/* STATS */}
                <div className="flex gap-10 pt-4 border-t border-neutral-800 font-mono text-xs">
                  <div>
                    <div className="text-green-400 text-lg">{posts.length}</div>
                    <div className="text-neutral-500">POSTS</div>
                  </div>

                  <div>
                    <div className="text-green-400 text-lg">
                      {replies.length}
                    </div>
                    <div className="text-neutral-500">REPLIES</div>
                  </div>

                  <div>
                    <div className="text-green-400 text-lg">{uptimeDays}</div>
                    <div className="text-neutral-500">UPTIME_DAYS</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* TABS */}
          <div className="flex border-b border-neutral-800 font-mono text-xs">
            <button
              onClick={() => setTab("posts")}
              className={`px-6 py-3 border-b-2 ${
                tab === "posts"
                  ? "border-green-400 text-green-400"
                  : "border-transparent text-neutral-500"
              }`}
            >
              POSTS
            </button>

            <button
              onClick={() => setTab("replies")}
              className={`px-6 py-3 border-b-2 ${
                tab === "replies"
                  ? "border-green-400 text-green-400"
                  : "border-transparent text-neutral-500"
              }`}
            >
              REPLIES
            </button>
          </div>

          {/* FEED */}
          <div className="flex flex-col gap-6">
            {activeList.length === 0 && (
              <div className="text-center text-neutral-500 font-mono text-xs py-10">
                NO_DATA_FOUND
              </div>
            )}

            {activeList.map((item) => (
              <PostCard key={item.id} post={item} />
            ))}
          </div>
        </div>
      </main>
    </div>
  );
}
