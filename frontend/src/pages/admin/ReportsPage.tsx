import { Card, CardContent } from "../../components/ui/card";
import { BarChart3 } from "lucide-react";

export default function ReportsPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Reports</h1>
      <Card>
        <CardContent className="p-8 text-center">
          <BarChart3 className="mx-auto text-gray-300 mb-3" size={40} />
          <p className="text-gray-500">Detailed reports and analytics coming soon.</p>
          <p className="text-xs text-gray-400 mt-1">View summary stats on the Admin Dashboard.</p>
        </CardContent>
      </Card>
    </div>
  );
}
