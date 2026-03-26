import Navbar from "@/components/navbar";
import Sidebar from "@/components/sidebar";
import FeedHeader from "@/components/feed/feed-header";
import ComposeBox from "@/components/feed/compose-box";
import FeedList from "@/components/feed/feed-list";

export default function Page() {
  return (
    <div className="min-h-screen bg-[#0e0e0e] text-white">
      <Navbar />
      <Sidebar />

      <main className="pt-20 md:pl-64 flex justify-center">
        <div className="w-full max-w-2xl flex flex-col gap-6 px-4">
          <FeedHeader />
          <ComposeBox />
          <FeedList />
        </div>
      </main>
    </div>
  );
}
