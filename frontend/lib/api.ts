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
