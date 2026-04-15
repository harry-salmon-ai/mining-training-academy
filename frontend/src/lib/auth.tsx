import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from "react";
import { apiFetch } from "./utils";

export interface User {
  id: string;
  email: string;
  name: string | null;
  firstName: string | null;
  lastName: string | null;
  image: string | null;
  role: "SUPERADMIN" | "ADMIN" | "INSTRUCTOR" | "SUPERVISOR" | "LEARNER";
  department: string | null;
  jobTitle: string | null;
  siteLocation: string | null;
  employeeId: string | null;
  isActive: boolean;
  createdAt: string;
  lastLoginAt: string | null;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (
    email: string,
    password: string,
    firstName: string,
    lastName: string
  ) => Promise<void>;
  logout: () => void;
  isAdmin: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const autoLogin = useCallback(async () => {
    try {
      const data = await apiFetch<{ token: string; user: User }>("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email: "admin@miningacademy.com", password: "admin123!" }),
      });
      localStorage.setItem("token", data.token);
      setUser(data.user);
    } catch {
      // backend not reachable yet — leave as unauthenticated
    }
  }, []);

  const fetchUser = useCallback(async () => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        await autoLogin();
        return;
      }
      const data = await apiFetch<User>("/auth/me");
      setUser(data);
    } catch {
      localStorage.removeItem("token");
      await autoLogin();
    } finally {
      setLoading(false);
    }
  }, [autoLogin]);

  useEffect(() => {
    // Check for token in URL (from OAuth callback)
    const params = new URLSearchParams(window.location.search);
    const urlToken = params.get("token");
    if (urlToken) {
      localStorage.setItem("token", urlToken);
      window.history.replaceState({}, "", window.location.pathname);
    }
    fetchUser();
  }, [fetchUser]);

  const login = async (email: string, password: string) => {
    const data = await apiFetch<{ token: string; user: User }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    localStorage.setItem("token", data.token);
    setUser(data.user);
  };

  const register = async (
    email: string,
    password: string,
    firstName: string,
    lastName: string
  ) => {
    const data = await apiFetch<{ token: string; user: User }>(
      "/auth/register",
      {
        method: "POST",
        body: JSON.stringify({ email, password, firstName, lastName }),
      }
    );
    localStorage.setItem("token", data.token);
    setUser(data.user);
  };

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
    apiFetch("/auth/logout", { method: "POST" }).catch(() => {});
  };

  const isAdmin =
    user?.role === "ADMIN" || user?.role === "SUPERADMIN";

  return (
    <AuthContext.Provider
      value={{ user, loading, login, register, logout, isAdmin }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within AuthProvider");
  return context;
}
