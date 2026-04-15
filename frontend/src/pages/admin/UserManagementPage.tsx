import { useEffect, useState } from "react";
import { apiFetch } from "../../lib/utils";
import { Input } from "../../components/ui/input";
import { Badge } from "../../components/ui/badge";
import { Card, CardContent } from "../../components/ui/card";

export default function UserManagementPage() {
  const [users, setUsers] = useState<any[]>([]);
  const [search, setSearch] = useState("");

  useEffect(() => {
    apiFetch<{ users: any[] }>(`/users?search=${search}`).then((d) => setUsers(d.users || [])).catch(() => {});
  }, [search]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">User Management</h1>
      </div>
      <Input placeholder="Search users..." value={search} onChange={(e) => setSearch(e.target.value)} className="max-w-md" />
      <Card>
        <CardContent className="p-0">
          <table className="w-full text-sm">
            <thead><tr className="border-b bg-gray-50"><th className="text-left px-4 py-3">Name</th><th className="text-left px-4 py-3">Email</th><th className="text-left px-4 py-3">Role</th><th className="text-left px-4 py-3">Status</th><th className="text-left px-4 py-3">Joined</th></tr></thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium">{u.name || `${u.firstName || ""} ${u.lastName || ""}`.trim()}</td>
                  <td className="px-4 py-3 text-gray-500">{u.email}</td>
                  <td className="px-4 py-3"><Badge variant="secondary" className="text-xs capitalize">{u.role?.toLowerCase()}</Badge></td>
                  <td className="px-4 py-3"><Badge variant={u.isActive ? "success" : "destructive"} className="text-xs">{u.isActive ? "Active" : "Inactive"}</Badge></td>
                  <td className="px-4 py-3 text-gray-500 text-xs">{new Date(u.createdAt).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
