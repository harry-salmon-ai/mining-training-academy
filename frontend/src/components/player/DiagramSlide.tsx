import { useMemo } from "react";

interface TreeNode {
  id: string;
  label: string;
  parent?: string;
  description?: string;
  highlight?: boolean;
}

interface DiagramSlideProps {
  content: {
    heading?: string;
    description?: string;
    diagramType?: string;
    svgContent?: string;
    diagramData?: {
      steps?: { label: string; description: string; highlight?: boolean }[];
      callout?: { text: string; description: string };
      // tree type
      nodes?: TreeNode[];
      // kpi-map type
      ratios?: { name: string; formula: string; owner: string; description: string }[];
      // layers type
      layers?: { label: string; description: string; stage: string }[];
      contextFactors?: string[];
      // architecture type
      sensors?: { name: string; description: string; functions: string[] }[];
      coreAlgorithm?: { stages: string[] };
    };
  };
}

function TreeDiagram({ nodes, callout }: { nodes: TreeNode[]; callout?: { text: string; description: string } }) {
  const levels = useMemo(() => {
    const childrenOf = new Map<string | undefined, TreeNode[]>();
    for (const n of nodes) {
      const key = n.parent ?? "__root__";
      if (!childrenOf.has(key)) childrenOf.set(key, []);
      childrenOf.get(key)!.push(n);
    }

    const result: TreeNode[][] = [];
    let current = childrenOf.get("__root__") ?? nodes.filter((n) => !n.parent);
    while (current.length > 0) {
      result.push(current);
      const next: TreeNode[] = [];
      for (const n of current) {
        const kids = childrenOf.get(n.id);
        if (kids) next.push(...kids);
      }
      current = next;
    }
    return result;
  }, [nodes]);

  return (
    <div className="space-y-2 py-4">
      {levels.map((level, li) => (
        <div key={li}>
          <div className="flex flex-wrap justify-center gap-3">
            {level.map((node) => (
              <div
                key={node.id}
                className={`px-4 py-3 rounded-lg border-2 text-center min-w-[140px] max-w-[200px] transition-all ${
                  node.highlight
                    ? "border-amber-400 bg-amber-50 text-amber-900 shadow-md"
                    : "border-gray-200 bg-white text-gray-700"
                }`}
              >
                <p className="font-semibold text-sm">{node.label}</p>
                {node.description && (
                  <p className="text-xs text-gray-500 mt-1 leading-tight">{node.description}</p>
                )}
              </div>
            ))}
          </div>
          {li < levels.length - 1 && (
            <div className="flex justify-center py-1">
              <svg width="24" height="20" viewBox="0 0 24 20" className="text-gray-300">
                <path d="M12 0 L12 14 M6 10 L12 16 L18 10" stroke="currentColor" strokeWidth="2" fill="none" />
              </svg>
            </div>
          )}
        </div>
      ))}

      {callout && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 text-center mt-4">
          <p className="font-bold text-yellow-800">{callout.text}</p>
          <p className="text-sm text-yellow-700 mt-1">{callout.description}</p>
        </div>
      )}
    </div>
  );
}

function KpiMapDiagram({
  ratios,
  callout,
}: {
  ratios: { name: string; formula: string; owner: string; description: string }[];
  callout?: { text: string; description: string };
}) {
  const ownerColors: Record<string, string> = {
    Maintenance: "bg-blue-100 text-blue-800",
    Operations: "bg-green-100 text-green-800",
    "Operations / Dispatch": "bg-green-100 text-green-800",
    "Planning / Management": "bg-purple-100 text-purple-800",
  };

  return (
    <div className="space-y-4 py-4">
      <div className="grid sm:grid-cols-2 gap-4">
        {ratios.map((r, i) => (
          <div key={i} className="bg-white border border-gray-200 rounded-xl p-5 space-y-3 shadow-sm">
            <h4 className="font-bold text-gray-900">{r.name}</h4>
            <code className="block bg-gray-50 border border-gray-200 rounded-md px-3 py-2 text-sm font-mono text-gray-800">
              {r.formula}
            </code>
            <p className="text-sm text-gray-600">{r.description}</p>
            <span
              className={`inline-block text-xs font-medium px-2.5 py-1 rounded-full ${
                ownerColors[r.owner] ?? "bg-gray-100 text-gray-700"
              }`}
            >
              {r.owner}
            </span>
          </div>
        ))}
      </div>

      {callout && (
        <div className="bg-amber-50 border border-amber-200 rounded-lg p-5 text-center">
          <p className="font-bold text-amber-900 text-lg">{callout.text}</p>
          <p className="text-sm text-amber-700 mt-1">{callout.description}</p>
        </div>
      )}
    </div>
  );
}

function LayersDiagram({
  layers,
  contextFactors,
  callout,
}: {
  layers: { label: string; description: string; stage: string }[];
  contextFactors?: string[];
  callout?: { text: string; description: string };
}) {
  const stageColors: Record<string, string> = {
    Perception: "border-blue-400 bg-blue-50",
    "Comprehension → Projection": "border-amber-400 bg-amber-50",
    "Decision → Action": "border-red-400 bg-red-50",
  };

  return (
    <div className="py-4 space-y-4">
      <div className="flex gap-6">
        <div className="flex-1 space-y-3">
          {[...layers].reverse().map((layer, i) => (
            <div
              key={i}
              className={`border-2 rounded-lg p-4 ${stageColors[layer.stage] ?? "border-gray-200 bg-gray-50"}`}
            >
              <div className="flex items-center justify-between gap-4 flex-wrap">
                <div>
                  <h4 className="font-bold text-gray-900">{layer.label}</h4>
                  <p className="text-sm text-gray-600 mt-1">{layer.description}</p>
                </div>
                <span className="text-xs font-medium bg-white/70 border border-gray-200 px-2.5 py-1 rounded-full whitespace-nowrap">
                  {layer.stage}
                </span>
              </div>
            </div>
          ))}
        </div>

        {contextFactors && contextFactors.length > 0 && (
          <div className="flex flex-col justify-center gap-2 min-w-[160px]">
            <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide">Context</p>
            {contextFactors.map((f, i) => (
              <span
                key={i}
                className="text-xs bg-gray-100 border border-gray-200 text-gray-700 px-3 py-1.5 rounded-full text-center"
              >
                {f}
              </span>
            ))}
          </div>
        )}
      </div>

      {callout && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 text-center">
          <p className="font-bold text-yellow-800">{callout.text}</p>
          <p className="text-sm text-yellow-700 mt-1">{callout.description}</p>
        </div>
      )}
    </div>
  );
}

function ArchitectureDiagram({
  sensors,
  coreAlgorithm,
}: {
  sensors: { name: string; description: string; functions: string[] }[];
  coreAlgorithm?: { stages: string[] };
}) {
  return (
    <div className="py-4">
      <div className="flex flex-col lg:flex-row gap-6 items-stretch">
        {/* Sensors */}
        <div className="flex-1 space-y-3">
          <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Sensor Stack</h4>
          {sensors.map((s, i) => (
            <div key={i} className="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
              <h5 className="font-bold text-gray-900">{s.name}</h5>
              <p className="text-sm text-gray-600 mt-1">{s.description}</p>
              <div className="flex flex-wrap gap-1.5 mt-2">
                {s.functions.map((fn, j) => (
                  <span key={j} className="text-xs bg-blue-50 text-blue-700 border border-blue-200 px-2 py-0.5 rounded">
                    {fn}
                  </span>
                ))}
              </div>
            </div>
          ))}
        </div>

        {/* Arrow connector */}
        <div className="flex items-center justify-center lg:flex-col">
          <div className="hidden lg:flex flex-col items-center gap-1">
            <div className="w-0.5 h-8 bg-gray-300" />
            <svg width="20" height="20" viewBox="0 0 20 20" className="text-gray-400">
              <path d="M4 4 L10 16 L16 4" stroke="currentColor" strokeWidth="2" fill="none" />
            </svg>
          </div>
          <div className="lg:hidden flex items-center gap-1">
            <div className="h-0.5 w-8 bg-gray-300" />
            <span className="text-gray-400 text-xl">&rarr;</span>
          </div>
        </div>

        {/* Core Algorithm */}
        {coreAlgorithm && (
          <div className="lg:w-64">
            <h4 className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">Core Algorithm</h4>
            <div className="bg-gradient-to-b from-gray-50 to-gray-100 border-2 border-gray-300 rounded-xl p-5 space-y-3">
              {coreAlgorithm.stages.map((stage, i) => (
                <div key={i}>
                  <div className="bg-white border border-gray-200 rounded-lg px-4 py-3 text-center shadow-sm">
                    <p className="font-semibold text-sm text-gray-800">{stage}</p>
                  </div>
                  {i < coreAlgorithm.stages.length - 1 && (
                    <div className="flex justify-center py-1">
                      <svg width="16" height="16" viewBox="0 0 16 16" className="text-gray-300">
                        <path d="M8 0 L8 10 M4 7 L8 12 L12 7" stroke="currentColor" strokeWidth="1.5" fill="none" />
                      </svg>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default function DiagramSlide({ content }: DiagramSlideProps) {
  const { diagramType, diagramData } = content;

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">{content.heading}</h2>
      {content.description && <p className="text-gray-600">{content.description}</p>}

      {diagramType === "image" && (content as any).imageUrl && (
        <div className="flex justify-center">
          <img
            src={(content as any).imageUrl}
            alt={content.heading}
            className="rounded-lg max-w-full border border-gray-200 shadow-sm"
          />
        </div>
      )}

      {content.svgContent && (
        <div className="flex justify-center" dangerouslySetInnerHTML={{ __html: content.svgContent }} />
      )}

      {diagramType === "tree" && diagramData?.nodes && (
        <TreeDiagram nodes={diagramData.nodes} callout={diagramData.callout} />
      )}

      {diagramType === "kpi-map" && diagramData?.ratios && (
        <KpiMapDiagram ratios={diagramData.ratios} callout={diagramData.callout} />
      )}

      {diagramType === "layers" && diagramData?.layers && (
        <LayersDiagram
          layers={diagramData.layers}
          contextFactors={diagramData.contextFactors}
          callout={diagramData.callout}
        />
      )}

      {diagramType === "architecture" && diagramData?.sensors && (
        <ArchitectureDiagram sensors={diagramData.sensors} coreAlgorithm={diagramData.coreAlgorithm} />
      )}

      {/* Fallback: flat steps (original behavior for process-style diagrams) */}
      {(!diagramType || diagramType === "process") && diagramData?.steps && (
        <div className="flex flex-wrap items-center justify-center gap-2 py-6">
          {diagramData.steps.map((step, i) => (
            <div key={i} className="flex items-center gap-2">
              <div
                className={`px-4 py-3 rounded-lg border-2 text-center min-w-[120px] ${
                  step.highlight
                    ? "border-primary-500 bg-primary-50 text-primary-900"
                    : "border-gray-200 bg-white text-gray-700"
                }`}
              >
                <p className="font-medium text-sm">{step.label}</p>
                <p className="text-xs text-gray-500 mt-1">{step.description}</p>
              </div>
              {i < diagramData.steps!.length - 1 && <span className="text-gray-300 text-2xl">&rarr;</span>}
            </div>
          ))}
        </div>
      )}

      {/* Callout for non-typed diagrams */}
      {(!diagramType || diagramType === "process") && diagramData?.callout && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 text-center">
          <p className="font-bold text-yellow-800">{diagramData.callout.text}</p>
          <p className="text-sm text-yellow-700 mt-1">{diagramData.callout.description}</p>
        </div>
      )}
    </div>
  );
}
