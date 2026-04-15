import { useEffect, useState } from "react";
import { apiFetch } from "../../lib/utils";
import { Card, CardContent } from "../../components/ui/card";

export default function AuditLogPage() {
  const [logs, setLogs] = useState<any[]>([]);
  useEffect(() => { apiFetch<any[]>("/audit").then(setLogs).catch(() => {}); }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Audit Log</h1>
      <Card>
        <CardContent className="p-0">
          <table className="w-full text-sm">
            <thead><tr className="border-b bg-gray-50"><th className="text-left px-4 py-3">User</th><th className="text-left px-4 py-3">Action</th><th className="text-left px-4 py-3">Entity</th><th className="text-left px-4 py-3">Time</th></tr></thead>
            <tbody>
              {logs.map((log) => (
                <tr key={log.id} className="border-b"><td className="px-4 py-3">{log.user?.name || log.userId}</td><td className="px-4 py-3 font-mono text-xs">{log.action}</td><td className="px-4 py-3 text-gray-500">{log.entity}</td><td className="px-4 py-3 text-gray-500 text-xs">{new Date(log.createdAt).toLocaleString()}</td></tr>
              ))}
              {logs.length === 0 && <tr><td colSpan={4} className="px-4 py-8 text-center text-gray-500">No audit logs yet</td></tr>}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
