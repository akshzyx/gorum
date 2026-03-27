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
  return json.data || json;
}

export async function getThread(id: string) {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/post/${id}/thread`,
  );
  return res.json();
}

export async function createReply(postId: string, content: string) {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/post/${postId}/reply`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ content }),
    },
  );

  if (!res.ok) {
    throw new Error("failed to create reply");
  }

  return res.json();
}
