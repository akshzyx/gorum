export default function FeedHeader() {
  return (
    <div className="flex justify-between border-b border-neutral-700 pb-2">
      <div>
        <h1 className="text-green-400 text-xl font-bold">CENTRAL_FEED</h1>
        <p className="text-xs text-neutral-500">SYNCHRONIZED</p>
      </div>

      <span className="text-green-400 border border-green-400 px-2 text-xs">
        PUBLIC_READ_ONLY
      </span>
    </div>
  );
}
