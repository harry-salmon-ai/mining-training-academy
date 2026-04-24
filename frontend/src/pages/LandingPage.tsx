import { Link, Navigate } from "react-router-dom";
import { useAuth } from "../lib/auth";
import { GraduationCap } from "lucide-react";

export default function LandingPage() {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="animate-spin rounded-full h-10 w-10 border-2 border-amber-500 border-t-transparent" />
      </div>
    );
  }

  if (user) {
    return <Navigate to="/dashboard" replace />;
  }

  return (
    <div className="min-h-screen bg-gray-900 flex flex-col items-center justify-center px-6 text-center">
      <div className="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-amber-500 to-orange-600 shadow-lg">
        <GraduationCap className="h-9 w-9 text-white" aria-hidden />
      </div>
      <h1 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
        Mining Training Academy
      </h1>
      <p className="mt-3 max-w-md text-gray-400">
        Operational excellence through learning. Sign in to access modules,
        track progress, and earn certificates.
      </p>
      <div className="mt-10 flex flex-wrap items-center justify-center gap-4">
        <Link
          to="/login"
          className="inline-flex items-center justify-center rounded-lg bg-amber-600 px-6 py-3 text-sm font-semibold text-white shadow hover:bg-amber-500 transition-colors"
        >
          Sign in
        </Link>
        <Link
          to="/register"
          className="inline-flex items-center justify-center rounded-lg border border-gray-600 bg-gray-800/50 px-6 py-3 text-sm font-medium text-gray-200 hover:border-gray-500 hover:bg-gray-800 transition-colors"
        >
          Create account
        </Link>
      </div>
    </div>
  );
}
