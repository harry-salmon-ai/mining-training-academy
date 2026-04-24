import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

const API_BASE = "/api";

export async function apiFetch<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const token = localStorage.getItem("token");
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options?.headers as Record<string, string>),
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  if (!res.ok) {
    const text = await res.text();
    let body: { error?: string; message?: string; detail?: string } = {};
    try {
      body = JSON.parse(text) as typeof body;
    } catch {
      const trimmed = text.trim();
      if (trimmed) {
        throw new Error(trimmed.slice(0, 400));
      }
      throw new Error(`Request failed (${res.status})`);
    }
    const msg =
      (typeof body.error === "string" && body.error) ||
      (typeof body.message === "string" && body.message) ||
      (typeof body.detail === "string" && body.detail) ||
      `Request failed (${res.status})`;
    throw new Error(msg);
  }

  return res.json();
}
