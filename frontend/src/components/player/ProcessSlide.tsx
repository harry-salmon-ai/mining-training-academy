import { AlertTriangle } from "lucide-react";

interface ProcessSlideProps {
  content: {
    heading?: string;
    description?: string;
    steps?: {
      title?: string;
      label?: string;
      description: string;
      icon?: string;
      image?: string;
      duration?: string;
      tier?: string;
      warnings?: string[];
    }[];
  };
}

const tierColors: Record<string, string> = {
  Design: "bg-blue-100 text-blue-800 border-blue-200",
  Operate: "bg-green-100 text-green-800 border-green-200",
  React: "bg-red-100 text-red-800 border-red-200",
};

export default function ProcessSlide({ content }: ProcessSlideProps) {
  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">{content.heading}</h2>
      {content.description && <p className="text-gray-600">{content.description}</p>}

      <div className="space-y-4">
        {content.steps?.map((step, i) => {
          const name = step.title || step.label || `Step ${i + 1}`;
          return (
            <div key={i} className="flex gap-4">
              <div className="flex flex-col items-center">
                <div className="w-8 h-8 rounded-full bg-primary-500 text-white flex items-center justify-center text-sm font-bold flex-shrink-0">
                  {i + 1}
                </div>
                {i < (content.steps?.length ?? 0) - 1 && (
                  <div className="w-0.5 flex-1 bg-gray-200 mt-2" />
                )}
              </div>
              <div className="pb-6 flex-1">
                <div className="flex items-center gap-2 flex-wrap">
                  <h3 className="font-bold">{name}</h3>
                  {step.tier && (
                    <span
                      className={`text-xs font-medium px-2 py-0.5 rounded-full border ${
                        tierColors[step.tier] ?? "bg-gray-100 text-gray-700 border-gray-200"
                      }`}
                    >
                      {step.tier}
                    </span>
                  )}
                  {step.duration && (
                    <span className="text-xs bg-gray-100 px-2 py-0.5 rounded text-gray-500">{step.duration}</span>
                  )}
                </div>
                <p className="text-gray-600 text-sm mt-1">{step.description}</p>
                {step.image && (
                  <img src={step.image} alt={name} className="mt-2 rounded max-w-sm" />
                )}
                {step.warnings && step.warnings.length > 0 && (
                  <div className="mt-2 space-y-1">
                    {step.warnings.map((w, j) => (
                      <div
                        key={j}
                        className="flex items-start gap-2 text-xs bg-yellow-50 text-yellow-800 border border-yellow-200 rounded p-2"
                      >
                        <AlertTriangle size={14} className="mt-0.5 flex-shrink-0" /> {w}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
