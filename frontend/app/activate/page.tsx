"use client";

import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export default function ActivatePage() {
  const params = useSearchParams();
  const router = useRouter();

  const token = params.get("token");

  const [status, setStatus] = useState<"loading" | "success" | "error">(
    "loading",
  );

  useEffect(() => {
    if (!token) {
      setStatus("error");
      return;
    }

    const activate = async () => {
      try {
        const res = await fetch(`${BASE_URL}/auth/activate`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ token }),
        });

        if (!res.ok) throw new Error();

        setStatus("success");

        // optional redirect after 2s
        setTimeout(() => {
          router.push("/");
        }, 2000);
      } catch {
        setStatus("error");
      }
    };

    activate();
  }, [token, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#0e0e0e] text-white">
      <div className="border border-neutral-800 p-6 w-[320px] text-center">
        {status === "loading" && <p>Activating your account...</p>}
        {status === "success" && (
          <p className="text-green-400">Account activated! Redirecting...</p>
        )}
        {status === "error" && (
          <p className="text-red-500">Invalid or expired activation link</p>
        )}
      </div>
    </div>
  );
}
