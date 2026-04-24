import { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { ChevronLeft, ChevronRight, List, X, StickyNote } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { useAuth } from "../../lib/auth";
import { Button } from "../../components/ui/button";
import { Progress } from "../../components/ui/progress";
import SlideRenderer from "../../components/player/SlideRenderer";
import SlideNotes from "../../components/modules/ModuleNotes";

interface SlideData {
  id: string;
  title: string;
  type: string;
  content: any;
  sortOrder: number;
  sectionId: string;
}

interface SectionData {
  id: string;
  title: string;
  sortOrder: number;
  slides: SlideData[];
}

export default function SlidePlayerPage() {
  const { slug } = useParams();
  const navigate = useNavigate();
  const { isAdmin } = useAuth();
  const [sections, setSections] = useState<SectionData[]>([]);
  const [allSlides, setAllSlides] = useState<SlideData[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [moduleTitle, setModuleTitle] = useState("");
  const [showMenu, setShowMenu] = useState(false);
  const [showNotes, setShowNotes] = useState(false);

  useEffect(() => {
    if (!slug) return;
    apiFetch<any>(`/modules/${slug}`).then((mod) => {
      setModuleTitle(mod.title);
      setSections(mod.sections || []);
      const slides = (mod.sections || []).flatMap((s: SectionData) => s.slides || []);
      setAllSlides(slides);
    }).catch(() => navigate("/modules"));
  }, [slug]);

  const currentSlide = allSlides[currentIndex];
  const progress = allSlides.length > 0 ? ((currentIndex + 1) / allSlides.length) * 100 : 0;

  const goNext = useCallback(() => {
    if (currentIndex < allSlides.length - 1) setCurrentIndex(currentIndex + 1);
  }, [currentIndex, allSlides.length]);

  const goPrev = useCallback(() => {
    if (currentIndex > 0) setCurrentIndex(currentIndex - 1);
  }, [currentIndex]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "ArrowRight") goNext();
      if (e.key === "ArrowLeft") goPrev();
      if (e.key === "Escape") setShowMenu(false);
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [goNext, goPrev]);

  // Record progress when slide changes
  useEffect(() => {
    if (!currentSlide) return;
    const timer = setTimeout(() => {
      apiFetch("/progress/slide", {
        method: "POST",
        body: JSON.stringify({ slideId: currentSlide.id, timeSpent: 5 }),
      }).catch(() => {});
    }, 3000);
    return () => clearTimeout(timer);
  }, [currentSlide?.id]);

  if (!currentSlide) {
    return (
      <div className="flex items-center justify-center h-[60vh]">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-white flex flex-col z-50">
      {/* Header */}
      <div className="h-12 border-b flex items-center justify-between px-4 bg-gray-900 text-white">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" className="text-white hover:bg-gray-800" onClick={() => navigate(`/modules/${slug}`)}>
            <X size={18} />
          </Button>
          <span className="text-sm font-medium truncate">{moduleTitle}</span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="text-gray-400">
            {currentIndex + 1} / {allSlides.length}
          </span>
          {isAdmin && (
            <Button variant="ghost" size="icon" className={`text-white hover:bg-gray-800 ${showNotes ? "bg-gray-700" : ""}`} onClick={() => setShowNotes(!showNotes)}>
              <StickyNote size={18} />
            </Button>
          )}
          <Button variant="ghost" size="icon" className="text-white hover:bg-gray-800" onClick={() => setShowMenu(!showMenu)}>
            <List size={18} />
          </Button>
        </div>
      </div>

      {/* Progress bar */}
      <Progress value={progress} className="h-1 rounded-none" />

      {/* Content area */}
      <div className="flex-1 flex overflow-hidden relative">
        {/* Main slide area */}
        <div className="flex-1 overflow-y-auto relative">
          {/* Section menu overlay */}
          {showMenu && (
            <div className="absolute inset-0 z-10 bg-white overflow-y-auto p-6">
              <h2 className="text-lg font-bold mb-4">Contents</h2>
              {sections.map((section, si) => (
                <div key={section.id} className="mb-4">
                  <h3 className="font-medium text-sm text-gray-500 uppercase mb-2">
                    Section {si + 1}: {section.title}
                  </h3>
                  {section.slides?.map((slide) => {
                    const globalIndex = allSlides.findIndex((s) => s.id === slide.id);
                    return (
                      <button
                        key={slide.id}
                        onClick={() => { setCurrentIndex(globalIndex); setShowMenu(false); }}
                        className={`block w-full text-left px-3 py-2 rounded text-sm ${
                          globalIndex === currentIndex ? "bg-primary-100 text-primary-700 font-medium" : "hover:bg-gray-100"
                        }`}
                      >
                        {slide.title}
                      </button>
                    );
                  })}
                </div>
              ))}
            </div>
          )}

          <div className="max-w-4xl mx-auto p-6">
            <SlideRenderer slide={currentSlide} />
          </div>
        </div>

        {/* Notes side panel */}
        {showNotes && isAdmin && (
          <div className="w-80 flex-shrink-0">
            <SlideNotes slideId={currentSlide.id} onClose={() => setShowNotes(false)} />
          </div>
        )}
      </div>

      {/* Navigation footer */}
      <div className="h-14 border-t flex items-center justify-between px-4 bg-white">
        <Button variant="outline" onClick={goPrev} disabled={currentIndex === 0}>
          <ChevronLeft size={18} className="mr-1" /> Previous
        </Button>
        <span className="text-sm text-gray-500">{currentSlide.title}</span>
        <Button onClick={goNext} disabled={currentIndex === allSlides.length - 1}>
          Next <ChevronRight size={18} className="ml-1" />
        </Button>
      </div>
    </div>
  );
}
