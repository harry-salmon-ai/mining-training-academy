import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Plus, Edit, Trash2, Eye, Send } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { Button } from "../../components/ui/button";
import { Badge } from "../../components/ui/badge";
import { Card, CardContent } from "../../components/ui/card";

const statusBadge: Record<string, { variant: any; label: string }> = {
  DRAFT: { variant: "secondary", label: "Draft" },
  IN_REVIEW: { variant: "warning", label: "In Review" },
  PUBLISHED: { variant: "success", label: "Published" },
  ARCHIVED: { variant: "outline", label: "Archived" },
};

export default function ModuleManagementPage() {
  const [modules, setModules] = useState<any[]>([]);

  const loadModules = () => {
    apiFetch<{ modules: any[] }>("/modules?limit=100").then((d) => setModules(d.modules || [])).catch(() => {});
  };

  useEffect(loadModules, []);

  const handlePublish = async (id: string) => {
    await apiFetch(`/modules/${id}/publish`, { method: "POST" });
    loadModules();
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this module?")) return;
    await apiFetch(`/modules/${id}`, { method: "DELETE" });
    loadModules();
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Module Management</h1>
        <Link to="/admin/modules/new">
          <Button><Plus size={16} className="mr-2" /> New Module</Button>
        </Link>
      </div>

      <Card>
        <CardContent className="p-0">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-gray-50">
                <th className="text-left px-4 py-3 font-medium">Title</th>
                <th className="text-left px-4 py-3 font-medium">Category</th>
                <th className="text-left px-4 py-3 font-medium">Level</th>
                <th className="text-left px-4 py-3 font-medium">Status</th>
                <th className="text-left px-4 py-3 font-medium">Author</th>
                <th className="text-right px-4 py-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {modules.map((mod) => (
                <tr key={mod.id} className="border-b last:border-0 hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium">{mod.title}</td>
                  <td className="px-4 py-3 text-gray-500">{mod.category?.name}</td>
                  <td className="px-4 py-3">
                    <Badge variant="secondary" className="text-xs capitalize">{mod.level?.toLowerCase()}</Badge>
                  </td>
                  <td className="px-4 py-3">
                    <Badge variant={statusBadge[mod.status]?.variant} className="text-xs">
                      {statusBadge[mod.status]?.label || mod.status}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 text-gray-500">{mod.author?.name}</td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex items-center justify-end gap-1">
                      <Link to={`/modules/${mod.slug}`}>
                        <Button variant="ghost" size="icon"><Eye size={16} /></Button>
                      </Link>
                      <Link to={`/admin/modules/${mod.id}/builder`}>
                        <Button variant="ghost" size="icon"><Edit size={16} /></Button>
                      </Link>
                      {mod.status === "DRAFT" && (
                        <Button variant="ghost" size="icon" onClick={() => handlePublish(mod.id)}>
                          <Send size={16} />
                        </Button>
                      )}
                      <Button variant="ghost" size="icon" onClick={() => handleDelete(mod.id)}>
                        <Trash2 size={16} className="text-red-500" />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
              {modules.length === 0 && (
                <tr><td colSpan={6} className="px-4 py-8 text-center text-gray-500">No modules yet</td></tr>
              )}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
