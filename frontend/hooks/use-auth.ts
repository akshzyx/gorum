"use client";

import { useEffect, useState } from "react";
import { getMe } from "@/lib/api";

type User = {
  id: string;
  username: string;
  avatar_url?: string | null;
};

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchMe = async () => {
    try {
      const data = await getMe();
      setUser(data);
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMe();
  }, []);

  return {
    user,
    setUser,
    loading,
    refetch: fetchMe,
  };
}
