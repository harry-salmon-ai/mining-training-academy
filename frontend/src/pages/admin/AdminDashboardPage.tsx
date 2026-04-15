import { useEffect, useState } from "react";
import { Users, BookOpen, GraduationCap, TrendingUp } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";

interface AdminStats {
  totalUsers: number;
  activeUsers: number;
  totalModules: number;
  publishedModules: number;
  totalEnrollments: number;
  completedEnrollments: number;
  completionRate: number;
  avgScore: number;
  recentActivity: any[];
}

export default function AdminDashboardPage() {
  const [stats, setStats] = useState<AdminStats | null>(null);

  useEffect(() => {
    apiFetch<AdminStats>("/reports/dashboard").then(setStats).catch(() => {});
  }, []);

  const cards = [
    { label: "Total Users", value: stats?.totalUsers ?? 0, sub: `${stats?.activeUsers ?? 0} active`, icon: Users, color: "bg-primary-100 text-primary-700" },
    { label: "Published Modules", value: stats?.publishedModules ?? 0, sub: `${stats?.totalModules ?? 0} total`, icon: BookOpen, color: "bg-green-100 text-green-700" },
    { label: "Completion Rate", value: `${stats?.completionRate ?? 0}%`, sub: `${stats?.completedEnrollments ?? 0} completed`, icon: GraduationCap, color: "bg-yellow-100 text-yellow-700" },
    { label: "Avg Quiz Score", value: `${stats?.avgScore ?? 0}%`, sub: `${stats?.totalEnrollments ?? 0} enrollments`, icon: TrendingUp, color: "bg-purple-100 text-purple-700" },
  ];

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Admin Dashboard</h1>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map((card) => (
          <Card key={card.label}>
            <CardContent className="p-4 flex items-center gap-4">
              <div className={`p-3 rounded-lg ${card.color}`}>
                <card.icon size={20} />
              </div>
              <div>
                <p className="text-2xl font-bold">{card.value}</p>
                <p className="text-sm text-gray-500">{card.label}</p>
                <p className="text-xs text-gray-400">{card.sub}</p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
        </CardHeader>
        <CardContent>
          {stats?.recentActivity && stats.recentActivity.length > 0 ? (
            <div className="space-y-2">
              {stats.recentActivity.map((log: any) => (
                <div key={log.id} className="flex items-center justify-between text-sm py-2 border-b last:border-0">
                  <div>
                    <span className="font-medium">{log.user?.name || "User"}</span>
                    <span className="text-gray-500 ml-2">{log.action}</span>
                  </div>
                  <span className="text-xs text-gray-400">
                    {new Date(log.createdAt).toLocaleString()}
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-sm">No recent activity</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
