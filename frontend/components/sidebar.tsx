export default function Sidebar() {
  return (
    <aside className="fixed left-0 top-0 h-screen w-64 border-r border-neutral-700 bg-[#0e0e0e] hidden md:flex flex-col p-6 gap-6 text-xs uppercase">
      <div className="text-green-400 font-bold">GORUM_SYS</div>

      <div className="flex flex-col gap-4">
        <span className="text-green-400">HOME</span>
        <span className="text-neutral-500">EXPLORE</span>
        <span className="text-neutral-500">COMMUNITIES</span>
        <span className="text-neutral-500">PROFILE</span>
      </div>

      <div className="mt-auto text-red-400">EXIT</div>
    </aside>
  );
}
