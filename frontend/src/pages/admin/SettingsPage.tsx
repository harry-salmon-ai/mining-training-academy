import { useEffect, useState } from "react";
import { apiFetch } from "../../lib/utils";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";

export default function SettingsPage() {
  const [settings, setSettings] = useState<any>({});
  const [saving, setSaving] = useState(false);

  useEffect(() => { apiFetch<any>("/settings").then(setSettings).catch(() => {}); }, []);

  const save = async () => {
    setSaving(true);
    await apiFetch("/settings", { method: "PATCH", body: JSON.stringify(settings) }).catch(() => {});
    setSaving(false);
  };

  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold">Platform Settings</h1>
      <Card>
        <CardHeader><CardTitle>General</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div><label className="block text-sm font-medium mb-1">Platform Name</label><Input value={settings.platformName || ""} onChange={(e) => setSettings({ ...settings, platformName: e.target.value })} /></div>
          <div><label className="block text-sm font-medium mb-1">Tagline</label><Input value={settings.tagline || ""} onChange={(e) => setSettings({ ...settings, tagline: e.target.value })} /></div>
          <div><label className="block text-sm font-medium mb-1">Support Email</label><Input value={settings.supportEmail || ""} onChange={(e) => setSettings({ ...settings, supportEmail: e.target.value })} /></div>
          <div><label className="block text-sm font-medium mb-1">Timezone</label><Input value={settings.timezone || ""} onChange={(e) => setSettings({ ...settings, timezone: e.target.value })} /></div>
          <Button onClick={save} disabled={saving}>{saving ? "Saving..." : "Save Settings"}</Button>
        </CardContent>
      </Card>
    </div>
  );
}
