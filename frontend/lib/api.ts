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

  const data = await res.json();

  return {
    thread: Array.isArray(data) ? data : data.data || [],
    target_id: data.target_id || id,
  };
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

export async function login(email: string, password: string) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) {
    throw new Error("invalid credentials");
  }

  return res.json();
}

export async function getMe() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/me`, {
    credentials: "include",
  });

  if (!res.ok) return null;

  return res.json();
}

export async function logout() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/logout`, {
    method: "POST",
    credentials: "include",
  });

  if (!res.ok) {
    throw new Error("logout failed");
  }
}
