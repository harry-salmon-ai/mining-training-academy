import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { BookOpen, GraduationCap, Award, TrendingUp } from "lucide-react";
import { useAuth } from "../../lib/auth";
import { apiFetch } from "../../lib/utils";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { Progress } from "../../components/ui/progress";

interface DashboardData {
  enrollments: number;
  completedModules: number;
  inProgress: any[];
  certificates: number;
  avgScore: number;
}

export default function DashboardPage() {
  const { user } = useAuth();
  const [data, setData] = useState<DashboardData | null>(null);

  useEffect(() => {
    apiFetch<DashboardData>("/reports/dashboard").then(setData).catch(() => {});
  }, []);

  const stats = [
    { label: "Enrolled Modules", value: data?.enrollments ?? 0, icon: BookOpen, color: "bg-primary-100 text-primary-700" },
    { label: "Completed", value: data?.completedModules ?? 0, icon: GraduationCap, color: "bg-green-100 text-green-700" },
    { label: "Certificates", value: data?.certificates ?? 0, icon: Award, color: "bg-yellow-100 text-yellow-700" },
    { label: "Avg Score", value: `${data?.avgScore ?? 0}%`, icon: TrendingUp, color: "bg-purple-100 text-purple-700" },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">
          Welcome back, {user?.firstName || user?.name || "Learner"}
        </h1>
        <p className="text-gray-500 text-sm mt-1">
          Track your progress and continue learning
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <Card key={stat.label}>
            <CardContent className="p-4 flex items-center gap-4">
              <div className={`p-3 rounded-lg ${stat.color}`}>
                <stat.icon size={20} />
              </div>
              <div>
                <p className="text-2xl font-bold">{stat.value}</p>
                <p className="text-sm text-gray-500">{stat.label}</p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">In Progress</h2>
          <Link to="/my-learning">
            <Button variant="ghost" size="sm">View all</Button>
          </Link>
        </div>

        {data?.inProgress && data.inProgress.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {data.inProgress.map((enrollment: any) => (
              <Card key={enrollment.id}>
                <CardHeader className="pb-2">
                  <CardTitle className="text-base">
                    {enrollment.module?.title}
                  </CardTitle>
                  <p className="text-xs text-gray-500">
                    {enrollment.module?.category?.name}
                  </p>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-500">Progress</span>
                      <span className="font-medium">{Math.round(enrollment.progress)}%</span>
                    </div>
                    <Progress value={enrollment.progress} />
                    <Link to={`/modules/${enrollment.module?.slug || enrollment.moduleId}`}>
                      <Button size="sm" className="w-full mt-2">Continue</Button>
                    </Link>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : (
          <Card>
            <CardContent className="p-8 text-center">
              <BookOpen className="mx-auto text-gray-300 mb-3" size={40} />
              <p className="text-gray-500">No modules in progress</p>
              <Link to="/modules">
                <Button variant="outline" className="mt-3">Browse Modules</Button>
              </Link>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
