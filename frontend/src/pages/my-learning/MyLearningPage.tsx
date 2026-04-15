import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { apiFetch } from "../../lib/utils";
import { Card, CardContent } from "../../components/ui/card";
import { Badge } from "../../components/ui/badge";
import { Progress } from "../../components/ui/progress";
import { Button } from "../../components/ui/button";

const statusColors: Record<string, string> = {
  NOT_STARTED: "secondary",
  IN_PROGRESS: "warning",
  COMPLETED: "success",
  OVERDUE: "destructive",
};

export default function MyLearningPage() {
  const [enrollments, setEnrollments] = useState<any[]>([]);

  useEffect(() => {
    apiFetch<any[]>("/enrollments").then(setEnrollments).catch(() => {});
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">My Learning</h1>

      {enrollments.length === 0 ? (
        <Card>
          <CardContent className="p-8 text-center">
            <p className="text-gray-500">You haven't enrolled in any modules yet.</p>
            <Link to="/modules"><Button variant="outline" className="mt-3">Browse Modules</Button></Link>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {enrollments.map((e) => (
            <Card key={e.id}>
              <CardContent className="p-4 flex items-center gap-4">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="font-medium">{e.module?.title || "Module"}</h3>
                    <Badge variant={statusColors[e.status] as any || "secondary"} className="text-xs">
                      {e.status?.replace("_", " ")}
                    </Badge>
                  </div>
                  <p className="text-xs text-gray-500">{e.module?.category?.name}</p>
                  <div className="mt-2 flex items-center gap-3">
                    <Progress value={e.progress} className="flex-1" />
                    <span className="text-sm font-medium w-12 text-right">{Math.round(e.progress)}%</span>
                  </div>
                </div>
                <Link to={`/modules/${e.module?.slug || e.moduleId}`}>
                  <Button size="sm">{e.progress > 0 ? "Continue" : "Start"}</Button>
                </Link>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
