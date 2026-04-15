import { Search, Bell, LogOut } from "lucide-react";
import { useAuth } from "../../lib/auth";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function TopBar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [search, setSearch] = useState("");

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (search.trim()) {
      navigate(`/modules?search=${encodeURIComponent(search)}`);
    }
  };

  return (
    <header className="h-14 border-b border-gray-200 bg-white flex items-center justify-between px-6">
      <form onSubmit={handleSearch} className="flex items-center gap-2 flex-1 max-w-md">
        <div className="relative flex-1">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" size={16} />
          <Input
            placeholder="Search modules..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9 h-8"
          />
        </div>
      </form>

      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={() => navigate("/notifications")} className="relative">
          <Bell size={18} />
        </Button>

        <div className="flex items-center gap-2 text-sm">
          <div className="w-8 h-8 rounded-full bg-primary-100 text-primary-700 flex items-center justify-center font-medium text-xs">
            {user?.name?.[0]?.toUpperCase() || user?.email?.[0]?.toUpperCase() || "?"}
          </div>
          <span className="text-gray-700 hidden md:block">{user?.name || user?.email}</span>
        </div>

        <Button variant="ghost" size="icon" onClick={logout}>
          <LogOut size={18} />
        </Button>
      </div>
    </header>
  );
}
