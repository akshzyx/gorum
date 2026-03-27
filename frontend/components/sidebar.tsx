"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  {
    name: "HOME",
    href: "/",
    icon: "fa-house",
  },
  {
    name: "EXPLORE",
    href: "/explore",
    icon: "fa-magnifying-glass",
  },
  {
    name: "PROFILE",
    href: "/profile",
    icon: "fa-user",
  },
  {
    name: "SETTINGS",
    href: "/settings",
    icon: "fa-gear",
  },
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed left-0 top-0 h-screen w-[180px] border-r border-neutral-800 bg-[#0e0e0e] hidden md:flex flex-col z-40 font-mono text-xs uppercase tracking-wider">
      {/* HEADER */}
      <div className="px-4 pt-6 pb-8">
        <div className="text-green-400 font-bold">GORUM_SYS</div>
        <div className="text-[10px] text-neutral-500">v1.0.4</div>
      </div>

      {/* NAV */}
      <div className="flex flex-col">
        {navItems.map((item) => {
          const active = pathname === item.href;

          return (
            <Link
              key={item.name}
              href={item.href}
              className={`flex items-center gap-3 px-4 py-3 transition-none
                ${
                  active
                    ? "text-green-400 border-l-2 border-green-400 bg-[#111]"
                    : "text-neutral-500 hover:text-green-400 hover:bg-[#111]"
                }
              `}
            >
              <i className={`fa-solid ${item.icon} w-4`} />
              <span>{item.name}</span>
            </Link>
          );
        })}
      </div>

      {/* EXIT */}
      <div className="mt-auto px-4 pb-6">
        <button
          className="flex items-center gap-3 text-red-400 hover:text-red-300"
          onClick={() => {
            window.location.href = "/";
          }}
        >
          <i className="fa-solid fa-power-off w-4" />
          <span>EXIT</span>
        </button>
      </div>
    </aside>
  );
}
