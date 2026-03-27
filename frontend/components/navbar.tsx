"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { useAuth } from "@/hooks/use-auth";
import LoginModal from "@/components/auth/login-modal";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

type User = {
  id: string;
  username: string;
  avatar_url?: string | null;
};

export default function Navbar({ onMenuClick }: { onMenuClick: () => void }) {
  const router = useRouter();

  // replaced manual state with auth hook
  const { user, refetch } = useAuth();

  // modal state (new)
  const [open, setOpen] = useState(false);

  const avatarUrl = user?.avatar_url;

  return (
    <>
      <nav className="fixed top-0 w-full h-14 border-b border-neutral-800 bg-[#0e0e0e] flex items-center justify-between px-4 md:px-6 z-50">
        {/* LEFT */}
        <div className="flex items-center gap-4">
          {/* HAMBURGER */}
          <i
            onClick={onMenuClick}
            className="fa-solid fa-bars text-neutral-400 md:hidden cursor-pointer"
          />

          <Link
            href="/"
            className="text-lg md:text-xl font-bold text-green-400"
          >
            GORUM
          </Link>
        </div>

        {/* RIGHT */}
        <div
          onClick={() => {
            if (!user) {
              setOpen(true);
            } else {
              router.push(`/user/${user.username}`);
            }
          }}
          className="group w-8 h-8 md:w-9 md:h-9 border border-green-400 flex items-center justify-center cursor-pointer hover:bg-green-400 overflow-hidden"
        >
          {avatarUrl ? (
            <Image
              src={avatarUrl}
              alt="avatar"
              width={40}
              height={40}
              className="w-full h-full object-cover"
            />
          ) : (
            <i className="fa-regular fa-user text-sm text-neutral-400 group-hover:text-black" />
          )}
        </div>
      </nav>

      {/* LOGIN MODAL */}
      <LoginModal
        open={open}
        onClose={() => setOpen(false)}
        onSuccess={refetch}
      />
    </>
  );
}
