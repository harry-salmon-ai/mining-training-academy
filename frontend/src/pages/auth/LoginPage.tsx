import { useState } from "react";
import { Shield, GraduationCap, Loader2 } from "lucide-react";
import { useAuth } from "../../lib/auth";

const roles = [
  {
    key: "admin",
    label: "Admin",
    description: "Full access to manage modules, users, and platform settings",
    email: "admin@miningacademy.com",
    password: "admin123!",
    icon: Shield,
    color: "from-amber-500 to-orange-600",
    ring: "ring-amber-400",
  },
  {
    key: "user",
    label: "Learner",
    description: "Browse and complete training modules",
    email: "user@miningacademy.com",
    password: "user123!",
    icon: GraduationCap,
    color: "from-blue-500 to-indigo-600",
    ring: "ring-blue-400",
  },
] as const;

export default function LoginPage() {
  const { login } = useAuth();
  const [loggingIn, setLoggingIn] = useState<string | null>(null);
  const [error, setError] = useState("");

  const handleRoleLogin = async (role: (typeof roles)[number]) => {
    setError("");
    setLoggingIn(role.key);
    try {
      await login(role.email, role.password);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Login failed";
      setError(message);
    } finally {
      setLoggingIn(null);
    }
  };

  return (
    <div className="space-y-6">
      {error && (
        <div className="bg-red-900/30 border border-red-700/50 text-red-300 text-sm p-3 rounded-lg text-center">
          {error}
        </div>
      )}

      <p className="text-center text-gray-400 text-sm">Select a role to sign in</p>

      <div className="grid gap-4">
        {roles.map((role) => {
          const Icon = role.icon;
          const isLoading = loggingIn === role.key;
          return (
            <button
              key={role.key}
              onClick={() => handleRoleLogin(role)}
              disabled={loggingIn !== null}
              className={`relative group w-full text-left rounded-xl border border-gray-700 bg-gray-800/50 p-5 transition-all hover:border-gray-500 hover:bg-gray-800 focus:outline-none focus:ring-2 ${role.ring} focus:ring-offset-2 focus:ring-offset-gray-900 disabled:opacity-60 disabled:cursor-wait`}
            >
              <div className="flex items-center gap-4">
                <div className={`flex-shrink-0 h-12 w-12 rounded-lg bg-gradient-to-br ${role.color} flex items-center justify-center`}>
                  {isLoading ? (
                    <Loader2 size={22} className="text-white animate-spin" />
                  ) : (
                    <Icon size={22} className="text-white" />
                  )}
                </div>
                <div>
                  <div className="font-semibold text-white text-lg">{role.label}</div>
                  <div className="text-sm text-gray-400 mt-0.5">{role.description}</div>
                </div>
              </div>
            </button>
          );
        })}
      </div>

      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-gray-700" />
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-gray-900 px-2 text-gray-500">or</span>
        </div>
      </div>

      <button
        onClick={() => (window.location.href = "/api/auth/google")}
        disabled={loggingIn !== null}
        className="w-full flex items-center justify-center gap-2 rounded-lg border border-gray-700 bg-gray-800/50 px-4 py-3 text-sm text-gray-300 transition-colors hover:bg-gray-800 hover:border-gray-500 disabled:opacity-60"
      >
        <svg className="w-4 h-4" viewBox="0 0 24 24">
          <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" />
          <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
          <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
          <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
        </svg>
        Continue with Google
      </button>
    </div>
  );
}
