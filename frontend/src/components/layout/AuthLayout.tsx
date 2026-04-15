import { Outlet, Navigate } from "react-router-dom";
import { useAuth } from "../../lib/auth";

export default function AuthLayout() {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="h-screen flex items-center justify-center bg-gray-50">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
      </div>
    );
  }

  if (user) {
    return <Navigate to="/dashboard" replace />;
  }

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-white">Mining Training Academy</h1>
          <p className="text-gray-400 text-sm mt-1">
            Operational Excellence Through Learning
          </p>
        </div>
        <Outlet />
      </div>
    </div>
  );
}
