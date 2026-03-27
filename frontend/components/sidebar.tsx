"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  {
    href: "/",
    icon: "fa-house",
  },
  {
    href: "/explore",
    icon: "fa-magnifying-glass",
  },
  {
    href: "/profile",
    icon: "fa-user",
  },
  {
    href: "/settings",
    icon: "fa-gear",
  },
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed left-0 top-0 h-screen w-[72px] border-r border-neutral-800 bg-[#0e0e0e] hidden md:flex flex-col items-center py-6 z-40">
      {/* TOP SPACING (for navbar overlap) */}
      <div className="mt-10 flex flex-col gap-4 w-full items-center">
        {navItems.map((item) => {
          const active = pathname === item.href;

          return (
            <Link
              key={item.href}
              href={item.href}
              className={`w-full flex justify-center py-4 transition-none
                ${
                  active
                    ? "bg-green-400 text-black"
                    : "text-neutral-500 hover:text-green-400 hover:bg-neutral-900"
                }
              `}
            >
              <i className={`fa-solid ${item.icon}`} />
            </Link>
          );
        })}
      </div>

      {/* BOTTOM EXIT */}
      <div className="mt-auto mb-6">
        <button
          className="text-red-400 hover:text-red-300"
          onClick={() => {
            window.location.href = "/";
          }}
        >
          <i className="fa-solid fa-power-off text-lg" />
        </button>
      </div>
    </aside>
  );
}
