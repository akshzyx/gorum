"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";

export default function Navbar() {
  const router = useRouter();

  return (
    <nav className="fixed top-0 w-full border-b border-neutral-700 bg-[#0e0e0e] flex justify-between items-center h-14 px-6 z-50">
      {/* LEFT → LOGO */}
      <Link
        href="/"
        className="text-xl font-bold text-green-400 tracking-tight"
      >
        GORUM
      </Link>

      {/* CENTER → EMPTY */}
      <div />

      {/* RIGHT → PROFILE */}
      <div className="flex items-center gap-4">
        <i
          onClick={() => router.push("/profile")}
          className="fa-regular fa-user text-neutral-400 cursor-pointer hover:text-green-400"
        />
      </div>
    </nav>
  );
}
