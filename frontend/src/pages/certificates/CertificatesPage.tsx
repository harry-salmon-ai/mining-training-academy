import { useEffect, useState } from "react";
import { Award } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { Card, CardContent } from "../../components/ui/card";

export default function CertificatesPage() {
  const [certs, setCerts] = useState<any[]>([]);

  useEffect(() => {
    apiFetch<any[]>("/certificates").then(setCerts).catch(() => {});
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">My Certificates</h1>
      {certs.length === 0 ? (
        <Card>
          <CardContent className="p-8 text-center">
            <Award className="mx-auto text-gray-300 mb-3" size={40} />
            <p className="text-gray-500">No certificates earned yet. Complete a module to earn your first certificate.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
          {certs.map((cert) => (
            <Card key={cert.id} className="overflow-hidden">
              <div className="bg-gradient-to-br from-gray-900 to-gray-700 p-6 text-white text-center">
                <Award size={32} className="mx-auto mb-2" />
                <h3 className="font-bold">{cert.module?.title}</h3>
              </div>
              <CardContent className="p-4 text-sm space-y-1">
                <p><span className="text-gray-500">Certificate #:</span> {cert.certificateNo}</p>
                <p><span className="text-gray-500">Issued:</span> {new Date(cert.issuedAt).toLocaleDateString()}</p>
                {cert.score && <p><span className="text-gray-500">Score:</span> {cert.score}%</p>}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
