"use client";

import { useState } from "react";
import { login } from "@/lib/api";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export default function LoginModal({
  open,
  onClose,
  onSuccess,
}: {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}) {
  const [mode, setMode] = useState<"login" | "signup">("login");

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  if (!open) return null;

  const handleSubmit = async () => {
    setLoading(true);
    setError("");

    try {
      if (mode === "login") {
        await login(email, password);
      } else {
        const res = await fetch(`${BASE_URL}/auth/signup`, {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            email,
            password,
            username,
          }),
        });

        if (!res.ok) {
          throw new Error("signup failed");
        }
      }

      onSuccess();
      onClose();
    } catch {
      setError(mode === "login" ? "Invalid credentials" : "Signup failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="fixed inset-0 bg-black/70 flex items-center justify-center z-[100]"
      onClick={onClose}
    >
      <div
        className="bg-[#0e0e0e] border border-neutral-800 p-6 w-[340px]"
        onClick={(e) => e.stopPropagation()}
      >
        {/* TOGGLE */}
        <div className="flex mb-4 border border-neutral-700">
          <button
            onClick={() => setMode("login")}
            className={`flex-1 py-2 text-sm ${
              mode === "login" ? "bg-green-400 text-black" : "text-neutral-400"
            }`}
          >
            Sign In
          </button>
          <button
            onClick={() => setMode("signup")}
            className={`flex-1 py-2 text-sm ${
              mode === "signup" ? "bg-green-400 text-black" : "text-neutral-400"
            }`}
          >
            Sign Up
          </button>
        </div>

        {/* USERNAME (signup only) */}
        {mode === "signup" && (
          <input
            className="w-full mb-3 p-2 bg-black border border-neutral-700 text-white"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        )}

        {/* EMAIL */}
        <input
          className="w-full mb-3 p-2 bg-black border border-neutral-700 text-white"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />

        {/* PASSWORD */}
        <input
          type="password"
          className="w-full mb-3 p-2 bg-black border border-neutral-700 text-white"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        {/* ERROR */}
        {error && <p className="text-red-500 text-sm mb-2">{error}</p>}

        {/* SUBMIT */}
        <button
          onClick={handleSubmit}
          disabled={loading}
          className="w-full bg-green-400 text-black py-2 font-semibold"
        >
          {loading
            ? mode === "login"
              ? "Logging in..."
              : "Creating account..."
            : mode === "login"
              ? "Login"
              : "Create Account"}
        </button>
      </div>
    </div>
  );
}
