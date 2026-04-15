import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";

interface ContentSlideProps {
  content: {
    heading?: string;
    introText?: string;
    keyFact?: string;
    items?: { title: string; description: string; icon?: string; imageUrl?: string }[];
  };
}

export default function ContentSlide({ content }: ContentSlideProps) {
  const [expanded, setExpanded] = useState<Set<number>>(new Set());

  const toggle = (index: number) => {
    const next = new Set(expanded);
    if (next.has(index)) next.delete(index);
    else next.add(index);
    setExpanded(next);
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">{content.heading}</h2>
      {content.introText && <p className="text-gray-600 text-lg">{content.introText}</p>}

      {content.items && (
        <div className="space-y-3">
          {content.items.map((item, i) => (
            <div key={i} className="border rounded-lg overflow-hidden">
              <button
                onClick={() => toggle(i)}
                className="w-full flex items-center gap-3 px-4 py-3 text-left hover:bg-gray-50 transition"
              >
                {expanded.has(i) ? <ChevronDown size={18} /> : <ChevronRight size={18} />}
                <span className="font-medium">{item.title}</span>
              </button>
              {expanded.has(i) && (
                <div className="px-4 pb-4 pl-10 space-y-3">
                  <p className="text-gray-600">{item.description}</p>
                  {item.imageUrl && (
                    <img
                      src={item.imageUrl}
                      alt={item.title}
                      className="rounded-lg max-w-full border border-gray-200"
                    />
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {content.keyFact && (
        <div className="bg-primary-50 border border-primary-200 rounded-lg px-4 py-3">
          <p className="text-sm font-semibold text-primary-800">{content.keyFact}</p>
        </div>
      )}
    </div>
  );
}
