import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { Clock, BookOpen, User, ArrowLeft, Play, CheckCircle2 } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { useAuth } from "../../lib/auth";
import { Button } from "../../components/ui/button";
import { Badge } from "../../components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";
import { Progress } from "../../components/ui/progress";

interface ModuleDetail {
  id: string;
  title: string;
  slug: string;
  description: string | null;
  level: string;
  duration: number | null;
  status: string;
  category: { id: string; name: string; color: string };
  author: { id: string; name: string };
  tags: { tag: string }[];
  sections: {
    id: string;
    title: string;
    sortOrder: number;
    slides: { id: string; title: string; type: string; sortOrder: number }[];
  }[];
}

export default function ModuleOverviewPage() {
  const { slug } = useParams();
  const navigate = useNavigate();
  useAuth();
  const [module, setModule] = useState<ModuleDetail | null>(null);
  const [enrollment, setEnrollment] = useState<any>(null);
  const [enrolling, setEnrolling] = useState(false);

  useEffect(() => {
    if (!slug) return;
    apiFetch<ModuleDetail>(`/modules/${slug}`).then(setModule).catch(() => navigate("/modules"));
  }, [slug]);

  useEffect(() => {
    if (!module) return;
    apiFetch<any[]>("/enrollments").then((enrollments) => {
      const found = enrollments.find((e: any) => e.moduleId === module.id);
      if (found) setEnrollment(found);
    }).catch(() => {});
  }, [module]);

  const handleEnroll = async () => {
    if (!module) return;
    setEnrolling(true);
    try {
      const data = await apiFetch<any>("/enrollments/self", {
        method: "POST",
        body: JSON.stringify({ moduleId: module.id }),
      });
      setEnrollment(data);
    } catch { }
    setEnrolling(false);
  };

  const totalSlides = module?.sections?.reduce((sum, s) => sum + (s.slides?.length || 0), 0) || 0;
  const firstSlideId = module?.sections?.[0]?.slides?.[0]?.id;

  if (!module) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <Link to="/modules" className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-900">
        <ArrowLeft size={16} /> Back to modules
      </Link>

      <div>
        <div className="flex items-center gap-2 mb-2">
          <Badge variant="secondary">{module.category?.name}</Badge>
          <Badge variant="secondary" className="capitalize">{module.level?.toLowerCase()}</Badge>
        </div>
        <h1 className="text-3xl font-bold">{module.title}</h1>
        {module.description && (
          <p className="text-gray-600 mt-2">{module.description}</p>
        )}
      </div>

      <div className="flex flex-wrap items-center gap-4 text-sm text-gray-500">
        {module.duration && (
          <span className="flex items-center gap-1"><Clock size={16} /> {module.duration} minutes</span>
        )}
        <span className="flex items-center gap-1"><BookOpen size={16} /> {totalSlides} slides</span>
        <span className="flex items-center gap-1"><User size={16} /> {module.author?.name}</span>
      </div>

      {enrollment && (
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm font-medium">Your Progress</span>
              <span className="text-sm font-bold">{Math.round(enrollment.progress || 0)}%</span>
            </div>
            <Progress value={enrollment.progress || 0} />
          </CardContent>
        </Card>
      )}

      <div className="flex gap-3">
        {enrollment ? (
          <Link to={`/modules/${module.slug}/learn${firstSlideId ? `/${firstSlideId}` : ""}`}>
            <Button size="lg">
              <Play size={18} className="mr-2" />
              {enrollment.progress > 0 ? "Continue Learning" : "Start Module"}
            </Button>
          </Link>
        ) : (
          <Button size="lg" onClick={handleEnroll} disabled={enrolling}>
            {enrolling ? "Enrolling..." : "Enroll & Start"}
          </Button>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Course Content</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {module.sections?.map((section, i) => (
              <div key={section.id} className="border rounded-lg p-4">
                <h3 className="font-medium text-sm">
                  Section {i + 1}: {section.title}
                </h3>
                <div className="mt-2 space-y-1">
                  {section.slides?.map((slide) => (
                    <div key={slide.id} className="flex items-center gap-2 text-sm text-gray-500 pl-4">
                      <CheckCircle2 size={14} className="text-gray-300" />
                      {slide.title}
                      <Badge variant="secondary" className="text-xs ml-auto">
                        {slide.type?.toLowerCase().replace("_", " ")}
                      </Badge>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {module.tags && module.tags.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {module.tags.map((t) => (
            <Badge key={t.tag} variant="outline">{t.tag}</Badge>
          ))}
        </div>
      )}
    </div>
  );
}
