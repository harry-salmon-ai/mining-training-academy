import { Link, useLocation } from "react-router-dom";
import {
  LayoutDashboard,
  BookOpen,
  GraduationCap,
  Award,
  User,
  Settings,
  Users,
  BarChart3,
  FolderOpen,
  ClipboardList,
  Shield,
  Bell,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { useState } from "react";
import { useAuth } from "../../lib/auth";
import { canAccess } from "../../lib/permissions";
import { cn } from "../../lib/utils";

const learnerNav = [
  { label: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { label: "Browse Modules", href: "/modules", icon: BookOpen },
  { label: "My Learning", href: "/my-learning", icon: GraduationCap },
  { label: "Certificates", href: "/certificates", icon: Award },
  { label: "Notifications", href: "/notifications", icon: Bell },
];

const adminNav = [
  { label: "Admin Dashboard", href: "/admin", icon: BarChart3 },
  { label: "Modules", href: "/admin/modules", icon: FolderOpen },
  { label: "Categories", href: "/admin/categories", icon: ClipboardList },
  { label: "Users", href: "/admin/users", icon: Users },
  { label: "Enrollments", href: "/admin/enrollments", icon: GraduationCap },
  { label: "Reports", href: "/admin/reports", icon: BarChart3 },
  { label: "Audit Log", href: "/admin/audit", icon: Shield },
  { label: "Settings", href: "/admin/settings", icon: Settings },
];

export default function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);
  const location = useLocation();
  const { user } = useAuth();

  const showAdmin = canAccess(user?.role, "modules", "create");

  return (
    <aside
      className={cn(
        "flex flex-col bg-gray-900 text-white transition-all duration-200",
        collapsed ? "w-16" : "w-64"
      )}
    >
      <div className="flex items-center justify-between p-4 border-b border-gray-800">
        {!collapsed && (
          <div>
            <h1 className="text-sm font-bold tracking-tight">Mining Training</h1>
            <p className="text-xs text-gray-400">Academy</p>
          </div>
        )}
        <button
          onClick={() => setCollapsed(!collapsed)}
          className="p-1 rounded hover:bg-gray-800"
        >
          {collapsed ? <ChevronRight size={16} /> : <ChevronLeft size={16} />}
        </button>
      </div>

      <nav className="flex-1 py-4 overflow-y-auto">
        <div className="px-3 mb-2">
          {!collapsed && (
            <p className="text-xs font-medium text-gray-500 uppercase tracking-wider mb-2">
              Learning
            </p>
          )}
          {learnerNav.map((item) => {
            const isActive = location.pathname === item.href;
            return (
              <Link
                key={item.href}
                to={item.href}
                className={cn(
                  "flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors mb-0.5",
                  isActive
                    ? "bg-primary-600 text-white"
                    : "text-gray-300 hover:bg-gray-800 hover:text-white"
                )}
              >
                <item.icon size={18} />
                {!collapsed && item.label}
              </Link>
            );
          })}
        </div>

        {showAdmin && (
          <div className="px-3 mt-4">
            {!collapsed && (
              <p className="text-xs font-medium text-gray-500 uppercase tracking-wider mb-2">
                Administration
              </p>
            )}
            {adminNav.map((item) => {
              const isActive = location.pathname === item.href;
              return (
                <Link
                  key={item.href}
                  to={item.href}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors mb-0.5",
                    isActive
                      ? "bg-primary-600 text-white"
                      : "text-gray-300 hover:bg-gray-800 hover:text-white"
                  )}
                >
                  <item.icon size={18} />
                  {!collapsed && item.label}
                </Link>
              );
            })}
          </div>
        )}
      </nav>

      <div className="p-4 border-t border-gray-800">
        <Link
          to="/profile"
          className="flex items-center gap-3 text-sm text-gray-300 hover:text-white"
        >
          <User size={18} />
          {!collapsed && (
            <div className="truncate">
              <p className="font-medium truncate">{user?.name || user?.email}</p>
              <p className="text-xs text-gray-500 capitalize">
                {user?.role?.toLowerCase()}
              </p>
            </div>
          )}
        </Link>
      </div>
    </aside>
  );
}
