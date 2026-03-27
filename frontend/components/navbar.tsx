"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

type User = {
  id: string;
  username: string;
  avatar_url?: string | null;
};

export default function Navbar() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadUser = async () => {
      try {
        const res = await fetch(`${BASE_URL}/me`, {
          credentials: "include",
        });

        if (!res.ok) {
          setUser(null);
          return;
        }

        const data = await res.json();
        setUser(data);
      } catch (err) {
        console.error(err);
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    loadUser();
  }, []);

  const avatarUrl = user?.avatar_url || null;

  return (
    <nav className="fixed top-0 w-full border-b border-neutral-800 bg-[#0e0e0e] flex justify-between items-center h-14 px-6 z-50">
      {/* LOGO */}
      <Link href="/" className="text-xl font-bold text-green-400">
        GORUM
      </Link>

      {/* EMPTY CENTER */}
      <div />

      {/* PROFILE */}
      <div
        onClick={() => router.push("/profile")}
        className="group w-9 h-9 border border-green-400 flex items-center justify-center cursor-pointer hover:bg-green-400 hover:text-black transition-none overflow-hidden"
      >
        {!loading && avatarUrl ? (
          <Image
            src={avatarUrl}
            alt="avatar"
            className="w-full h-full object-cover"
          />
        ) : (
          <i className="fa-regular fa-user text-sm text-neutral-400 group-hover:text-black" />
        )}
      </div>
    </nav>
  );
}
