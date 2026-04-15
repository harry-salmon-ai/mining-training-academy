import { useState } from "react";
import { useAuth } from "../../lib/auth";
import { apiFetch } from "../../lib/utils";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";

export default function ProfilePage() {
  const { user } = useAuth();
  const [firstName, setFirstName] = useState(user?.firstName || "");
  const [lastName, setLastName] = useState(user?.lastName || "");
  const [department, setDepartment] = useState(user?.department || "");
  const [jobTitle, setJobTitle] = useState(user?.jobTitle || "");
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  const handleSave = async () => {
    setSaving(true);
    try {
      await apiFetch("/profile", {
        method: "PATCH",
        body: JSON.stringify({ firstName, lastName, department, jobTitle }),
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch {}
    setSaving(false);
  };

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Profile Settings</h1>
      <Card>
        <CardHeader>
          <CardTitle>Personal Information</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">First Name</label>
              <Input value={firstName} onChange={(e) => setFirstName(e.target.value)} />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Last Name</label>
              <Input value={lastName} onChange={(e) => setLastName(e.target.value)} />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Email</label>
            <Input value={user?.email || ""} disabled />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Department</label>
            <Input value={department} onChange={(e) => setDepartment(e.target.value)} />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Job Title</label>
            <Input value={jobTitle} onChange={(e) => setJobTitle(e.target.value)} />
          </div>
          <div className="flex items-center gap-3">
            <Button onClick={handleSave} disabled={saving}>
              {saving ? "Saving..." : "Save Changes"}
            </Button>
            {saved && <span className="text-sm text-green-600">Saved!</span>}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
