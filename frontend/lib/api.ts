const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getPosts(cursor?: string) {
  const url = new URL(`${BASE_URL}/post`);

  if (cursor) {
    url.searchParams.append("cursor", cursor);
  }

  const res = await fetch(url.toString(), {
    credentials: "include",
  });

  if (!res.ok) {
    throw new Error("failed to fetch posts");
  }

  return res.json();
}

export async function getPostById(id: string) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/post/${id}`);
  const json = await res.json();
  return json.data || json; // ✅ FIX
}

export async function getReplies(id: string) {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/post/${id}/replies`,
  );
  return res.json(); // already handled in page
}
