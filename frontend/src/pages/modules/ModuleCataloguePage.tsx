import { useEffect, useState } from "react";
import { Search, Filter } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { Input } from "../../components/ui/input";
import { Badge } from "../../components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";
import { Link } from "react-router-dom";

interface Module {
  id: string;
  title: string;
  slug: string;
  description: string | null;
  thumbnail: string | null;
  status: string;
  level: string;
  duration: number | null;
  category: { id: string; name: string; color: string };
  tags: { tag: string }[];
  author: { name: string };
}

interface Category {
  id: string;
  name: string;
  slug: string;
  color: string;
}

const levelColors: Record<string, string> = {
  FOUNDATION: "bg-green-100 text-green-800",
  INTERMEDIATE: "bg-yellow-100 text-yellow-800",
  ADVANCED: "bg-orange-100 text-orange-800",
  SPECIALIST: "bg-red-100 text-red-800",
};

export default function ModuleCataloguePage() {
  const [modules, setModules] = useState<Module[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [search, setSearch] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("");
  const [selectedLevel, setSelectedLevel] = useState("");

  useEffect(() => {
    apiFetch<Category[]>("/categories").then(setCategories).catch(() => {});
  }, []);

  useEffect(() => {
    const params = new URLSearchParams();
    if (search) params.set("search", search);
    if (selectedCategory) params.set("category", selectedCategory);
    if (selectedLevel) params.set("level", selectedLevel);

    apiFetch<{ modules: Module[] }>(`/modules?${params.toString()}`)
      .then((data) => setModules(data.modules || []))
      .catch(() => {});
  }, [search, selectedCategory, selectedLevel]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Module Catalogue</h1>
        <p className="text-gray-500 text-sm mt-1">Browse and enroll in training modules</p>
      </div>

      <div className="flex flex-col sm:flex-row gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" size={16} />
          <Input
            placeholder="Search modules..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9"
          />
        </div>
        <select
          value={selectedCategory}
          onChange={(e) => setSelectedCategory(e.target.value)}
          className="border border-gray-300 rounded-md px-3 py-2 text-sm bg-white"
        >
          <option value="">All Categories</option>
          {categories.map((cat) => (
            <option key={cat.id} value={cat.id}>{cat.name}</option>
          ))}
        </select>
        <select
          value={selectedLevel}
          onChange={(e) => setSelectedLevel(e.target.value)}
          className="border border-gray-300 rounded-md px-3 py-2 text-sm bg-white"
        >
          <option value="">All Levels</option>
          <option value="FOUNDATION">Foundation</option>
          <option value="INTERMEDIATE">Intermediate</option>
          <option value="ADVANCED">Advanced</option>
          <option value="SPECIALIST">Specialist</option>
        </select>
      </div>

      {modules.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {modules.map((mod) => (
            <Link key={mod.id} to={`/modules/${mod.slug}`}>
              <Card className="h-full hover:shadow-md transition-shadow cursor-pointer">
                {mod.thumbnail && (
                  <div className="h-40 bg-gray-100 rounded-t-xl overflow-hidden">
                    <img src={mod.thumbnail} alt={mod.title} className="w-full h-full object-cover" />
                  </div>
                )}
                <CardHeader className={mod.thumbnail ? "pt-3" : ""}>
                  <div className="flex items-center gap-2 mb-1">
                    <Badge variant="secondary" className="text-xs">
                      {mod.category?.name}
                    </Badge>
                    <Badge className={`text-xs ${levelColors[mod.level] || ""}`} variant="secondary">
                      {mod.level?.toLowerCase()}
                    </Badge>
                  </div>
                  <CardTitle className="text-base">{mod.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  {mod.description && (
                    <p className="text-sm text-gray-500 line-clamp-2 mb-3">{mod.description}</p>
                  )}
                  <div className="flex items-center justify-between text-xs text-gray-400">
                    {mod.duration && <span>{mod.duration} min</span>}
                    <span>By {mod.author?.name}</span>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      ) : (
        <Card>
          <CardContent className="p-8 text-center">
            <Filter className="mx-auto text-gray-300 mb-3" size={40} />
            <p className="text-gray-500">No modules found</p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
