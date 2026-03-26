export default function Navbar() {
  return (
    <nav className="fixed top-0 w-full border-b border-neutral-700 bg-[#0e0e0e] flex justify-between items-center h-16 px-6 z-50 uppercase text-sm">
      <div className="text-2xl font-bold text-green-400">GORUM</div>

      <div className="hidden md:flex gap-8">
        <span className="text-green-400 border-b-2 border-green-400">FEED</span>
        <span className="text-neutral-500">EXPLORE</span>
        <span className="text-neutral-500">COMMUNITIES</span>
      </div>
    </nav>
  );
}
