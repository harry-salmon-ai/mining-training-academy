import { useState } from "react";
import { cn } from "../../lib/utils";

interface EquipmentSlideProps {
  content: {
    heading?: string;
    description?: string;
    image?: string;
    specifications?: { label: string; value: string }[];
    advantages?: string[];
    limitations?: string[];
  };
}

export default function EquipmentSlide({ content }: EquipmentSlideProps) {
  const [tab, setTab] = useState<"specs" | "pros" | "cons">("specs");

  const tabs = [
    { key: "specs" as const, label: "Specifications" },
    { key: "pros" as const, label: "Advantages" },
    { key: "cons" as const, label: "Limitations" },
  ];

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">{content.heading}</h2>
      {content.description && <p className="text-gray-600">{content.description}</p>}

      {content.image && (
        <img src={content.image} alt={content.heading} className="rounded-lg max-w-full mx-auto" />
      )}

      <div className="border rounded-lg overflow-hidden">
        <div className="flex border-b">
          {tabs.map((t) => (
            <button
              key={t.key}
              onClick={() => setTab(t.key)}
              className={cn(
                "flex-1 px-4 py-2 text-sm font-medium transition",
                tab === t.key ? "bg-gray-900 text-white" : "bg-gray-50 text-gray-600 hover:bg-gray-100"
              )}
            >
              {t.label}
            </button>
          ))}
        </div>
        <div className="p-4">
          {tab === "specs" && content.specifications && (
            <table className="w-full text-sm">
              <tbody>
                {content.specifications.map((spec, i) => (
                  <tr key={i} className={i % 2 === 0 ? "bg-gray-50" : ""}>
                    <td className="px-3 py-2 font-medium text-gray-700">{spec.label}</td>
                    <td className="px-3 py-2 text-gray-600">{spec.value}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          {tab === "pros" && content.advantages && (
            <ul className="space-y-2">
              {content.advantages.map((a, i) => (
                <li key={i} className="flex items-start gap-2 text-sm">
                  <span className="text-green-500 mt-0.5">+</span> {a}
                </li>
              ))}
            </ul>
          )}
          {tab === "cons" && content.limitations && (
            <ul className="space-y-2">
              {content.limitations.map((l, i) => (
                <li key={i} className="flex items-start gap-2 text-sm">
                  <span className="text-red-500 mt-0.5">-</span> {l}
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
}
