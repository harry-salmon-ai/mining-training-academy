import { useEffect, useState } from "react";
import { Bell, Check } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { Button } from "../../components/ui/button";
import { Card, CardContent } from "../../components/ui/card";

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<any[]>([]);

  const load = () => { apiFetch<any[]>("/notifications").then(setNotifications).catch(() => {}); };
  useEffect(load, []);

  const markAllRead = async () => {
    await apiFetch("/notifications/mark-all-read", { method: "POST" });
    load();
  };

  return (
    <div className="space-y-6 max-w-2xl">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Notifications</h1>
        <Button variant="outline" size="sm" onClick={markAllRead}><Check size={14} className="mr-1" /> Mark all read</Button>
      </div>
      {notifications.length === 0 ? (
        <Card><CardContent className="p-8 text-center"><Bell className="mx-auto text-gray-300 mb-3" size={40} /><p className="text-gray-500">No notifications</p></CardContent></Card>
      ) : (
        <div className="space-y-2">
          {notifications.map((n) => (
            <Card key={n.id} className={n.read ? "opacity-60" : ""}>
              <CardContent className="p-4">
                <p className="font-medium text-sm">{n.title}</p>
                <p className="text-sm text-gray-500">{n.message}</p>
                <p className="text-xs text-gray-400 mt-1">{new Date(n.createdAt).toLocaleString()}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
