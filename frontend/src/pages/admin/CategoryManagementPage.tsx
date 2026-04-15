import { useEffect, useState } from "react";
import { Plus, Trash2 } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Card, CardContent } from "../../components/ui/card";

export default function CategoryManagementPage() {
  const [categories, setCategories] = useState<any[]>([]);
  const [name, setName] = useState("");
  const [color, setColor] = useState("#185FA5");

  const load = () => { apiFetch<any[]>("/categories").then(setCategories).catch(() => {}); };
  useEffect(load, []);

  const create = async () => {
    if (!name.trim()) return;
    await apiFetch("/categories", { method: "POST", body: JSON.stringify({ name, color, icon: name[0]?.toUpperCase() }) });
    setName("");
    load();
  };

  const remove = async (id: string) => {
    if (!confirm("Delete category?")) return;
    await apiFetch(`/categories/${id}`, { method: "DELETE" });
    load();
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Categories</h1>
      <div className="flex gap-2">
        <Input placeholder="Category name..." value={name} onChange={(e) => setName(e.target.value)} className="max-w-xs" />
        <input type="color" value={color} onChange={(e) => setColor(e.target.value)} className="w-10 h-9 rounded border cursor-pointer" />
        <Button onClick={create}><Plus size={16} className="mr-1" /> Add</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          <table className="w-full text-sm">
            <thead><tr className="border-b bg-gray-50"><th className="text-left px-4 py-3">Name</th><th className="text-left px-4 py-3">Slug</th><th className="text-left px-4 py-3">Color</th><th className="text-right px-4 py-3">Actions</th></tr></thead>
            <tbody>
              {categories.map((c) => (
                <tr key={c.id} className="border-b"><td className="px-4 py-3 font-medium">{c.name}</td><td className="px-4 py-3 text-gray-500">{c.slug}</td><td className="px-4 py-3"><span className="inline-block w-4 h-4 rounded" style={{ backgroundColor: c.color }} /></td><td className="px-4 py-3 text-right"><Button variant="ghost" size="icon" onClick={() => remove(c.id)}><Trash2 size={16} className="text-red-500" /></Button></td></tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
