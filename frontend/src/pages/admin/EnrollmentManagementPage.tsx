import { useEffect, useState } from "react";
import { apiFetch } from "../../lib/utils";
import { Badge } from "../../components/ui/badge";
import { Card, CardContent } from "../../components/ui/card";
import { Progress } from "../../components/ui/progress";

export default function EnrollmentManagementPage() {
  const [enrollments, setEnrollments] = useState<any[]>([]);

  useEffect(() => {
    apiFetch<any[]>("/enrollments").then(setEnrollments).catch(() => {});
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Enrollments</h1>
      <Card>
        <CardContent className="p-0">
          <table className="w-full text-sm">
            <thead><tr className="border-b bg-gray-50"><th className="text-left px-4 py-3">User</th><th className="text-left px-4 py-3">Module</th><th className="text-left px-4 py-3">Status</th><th className="text-left px-4 py-3">Progress</th><th className="text-left px-4 py-3">Enrolled</th></tr></thead>
            <tbody>
              {enrollments.map((e) => (
                <tr key={e.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3">{e.user?.name || e.userId}</td>
                  <td className="px-4 py-3">{e.module?.title || e.moduleId}</td>
                  <td className="px-4 py-3"><Badge variant="secondary" className="text-xs">{e.status?.replace("_", " ")}</Badge></td>
                  <td className="px-4 py-3 w-32"><Progress value={e.progress} /></td>
                  <td className="px-4 py-3 text-gray-500 text-xs">{new Date(e.enrolledAt).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
