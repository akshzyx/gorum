"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { name: "HOME", href: "/", icon: "fa-house" },
  { name: "EXPLORE", href: "/explore", icon: "fa-magnifying-glass" },
  { name: "PROFILE", href: "/profile", icon: "fa-user" },
  { name: "SETTINGS", href: "/settings", icon: "fa-gear" },
];

export default function Sidebar({
  open,
  setOpen,
}: {
  open: boolean;
  setOpen: (v: boolean) => void;
}) {
  const pathname = usePathname();

  return (
    <>
      {/* OVERLAY (mobile) */}
      {open && (
        <div
          onClick={() => setOpen(false)}
          className="fixed inset-0 bg-black/50 z-40 md:hidden"
        />
      )}

      <aside
        className={`
          fixed top-0 left-0 h-screen w-[180px] bg-[#0e0e0e] border-r border-neutral-800 z-50
          transform transition-transform duration-200
          ${open ? "translate-x-0" : "-translate-x-full"}
          md:translate-x-0 md:flex md:flex-col hidden md:block
          font-mono text-xs uppercase tracking-wider
        `}
      >
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
                onClick={() => setOpen(false)}
                className={`flex items-center gap-3 px-4 py-3
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
            onClick={() => (window.location.href = "/")}
          >
            <i className="fa-solid fa-power-off w-4" />
            <span>EXIT</span>
          </button>
        </div>
      </aside>
    </>
  );
}
