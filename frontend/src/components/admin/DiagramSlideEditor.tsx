import { Plus, X, Eye, GripVertical } from "lucide-react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Textarea } from "../ui/textarea";
import { Label } from "../ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import ImageUploader from "./ImageUploader";
import DiagramSlide from "../player/DiagramSlide";
import {
  DndContext,
  closestCenter,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
  useSortable,
  arrayMove,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

// ─── Types ────────────────────────────────────────────────────────────────────

const DIAGRAM_TYPES = [
  { value: "kpi-map", label: "KPI / Metric Cards" },
  { value: "layers", label: "Layer Stack (e.g. defence model)" },
  { value: "architecture", label: "Architecture (sensors + algorithm)" },
  { value: "tree", label: "Tree / Hierarchy" },
  { value: "image", label: "Image / Uploaded Diagram" },
  { value: "", label: "Flow (arrow steps)" },
] as const;

const KPI_OWNERS = [
  "Operations",
  "Maintenance",
  "Operations / Dispatch",
  "Planning / Management",
  "Technology",
  "Technology / Operations",
  "Safety",
];

const LAYER_STAGES = [
  "Design",
  "Operate",
  "React",
  "Perception",
  "Comprehension → Projection",
  "Decision → Action",
];

// ─── Sortable drag handle item ────────────────────────────────────────────────

function SortableItem({ id, children }: { id: string; children: React.ReactNode }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id });
  return (
    <div
      ref={setNodeRef}
      style={{
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.4 : 1,
        zIndex: isDragging ? 10 : undefined,
      }}
      className="relative"
    >
      <button
        type="button"
        className="absolute left-1 top-3.5 cursor-grab active:cursor-grabbing touch-none text-gray-300 hover:text-gray-500"
        {...attributes}
        {...listeners}
        tabIndex={-1}
      >
        <GripVertical size={14} />
      </button>
      <div className="pl-6">{children}</div>
    </div>
  );
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

function CalloutEditor({
  callout,
  onChange,
}: {
  callout?: { text: string; description: string };
  onChange: (c: { text: string; description: string } | undefined) => void;
}) {
  if (!callout) {
    return (
      <button
        type="button"
        onClick={() => onChange({ text: "", description: "" })}
        className="text-xs text-gray-500 hover:text-gray-700 flex items-center gap-1"
      >
        <Plus size={11} /> Add callout box
      </button>
    );
  }
  return (
    <div className="border border-yellow-200 bg-yellow-50 rounded-lg p-3 space-y-2">
      <div className="flex items-center justify-between">
        <Label className="text-xs text-yellow-700">Callout Box</Label>
        <button type="button" onClick={() => onChange(undefined)} className="text-yellow-600 hover:text-red-500">
          <X size={13} />
        </button>
      </div>
      <Input
        placeholder="Bold headline"
        value={callout.text}
        onChange={(e) => onChange({ ...callout, text: e.target.value })}
      />
      <Input
        placeholder="Supporting detail"
        value={callout.description}
        onChange={(e) => onChange({ ...callout, description: e.target.value })}
      />
    </div>
  );
}

// ─── KPI Map editor ───────────────────────────────────────────────────────────

function KpiMapEditor({ diagramData, onChange }: { diagramData: any; onChange: (d: any) => void }) {
  const ratios: any[] = diagramData.ratios || [];
  const set = (patch: object) => onChange({ ...diagramData, ...patch });
  const sensors = useSensors(useSensor(PointerSensor));
  const ids = ratios.map((_, i) => String(i));

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (over && active.id !== over.id) {
      const from = ids.indexOf(String(active.id));
      const to = ids.indexOf(String(over.id));
      set({ ratios: arrayMove(ratios, from, to) });
    }
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <Label>KPI Cards</Label>
        <Button
          type="button" variant="ghost" size="sm"
          onClick={() => set({ ratios: [...ratios, { name: "", formula: "", owner: "Operations", description: "" }] })}
        >
          <Plus size={13} className="mr-1" /> Add card
        </Button>
      </div>

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={ids} strategy={verticalListSortingStrategy}>
          {ratios.map((r, i) => (
            <SortableItem key={i} id={String(i)}>
              <div className="border border-gray-200 rounded-lg p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-gray-500">Card {i + 1}</span>
                  <button type="button" onClick={() => set({ ratios: ratios.filter((_, j) => j !== i) })} className="text-gray-400 hover:text-red-500">
                    <X size={14} />
                  </button>
                </div>
                <Input placeholder="KPI name (e.g. Physical Availability)" value={r.name}
                  onChange={(e) => set({ ratios: ratios.map((x, j) => (j === i ? { ...x, name: e.target.value } : x)) })} />
                <Input placeholder="Formula (e.g. (Available Hours / Scheduled Hours) × 100)" value={r.formula} className="font-mono text-sm"
                  onChange={(e) => set({ ratios: ratios.map((x, j) => (j === i ? { ...x, formula: e.target.value } : x)) })} />
                <Textarea placeholder="What this measures and why it matters" rows={2} value={r.description}
                  onChange={(e) => set({ ratios: ratios.map((x, j) => (j === i ? { ...x, description: e.target.value } : x)) })} />
                <div className="space-y-1">
                  <Label className="text-xs text-gray-500">Owner</Label>
                  <Select value={r.owner || "Operations"} onValueChange={(v) => set({ ratios: ratios.map((x, j) => (j === i ? { ...x, owner: v } : x)) })}>
                    <SelectTrigger className="h-8 text-xs"><SelectValue /></SelectTrigger>
                    <SelectContent>{KPI_OWNERS.map((o) => <SelectItem key={o} value={o}>{o}</SelectItem>)}</SelectContent>
                  </Select>
                </div>
              </div>
            </SortableItem>
          ))}
        </SortableContext>
      </DndContext>

      <CalloutEditor callout={diagramData.callout} onChange={(c) => set({ callout: c })} />
    </div>
  );
}

// ─── Layers editor ────────────────────────────────────────────────────────────

function LayersEditor({ diagramData, onChange }: { diagramData: any; onChange: (d: any) => void }) {
  const layers: any[] = diagramData.layers || [];
  const contextFactors: string[] = diagramData.contextFactors || [];
  const set = (patch: object) => onChange({ ...diagramData, ...patch });
  const sensors = useSensors(useSensor(PointerSensor));
  const ids = layers.map((_, i) => String(i));

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (over && active.id !== over.id) {
      const from = ids.indexOf(String(active.id));
      const to = ids.indexOf(String(over.id));
      set({ layers: arrayMove(layers, from, to) });
    }
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <Label>Layers <span className="text-xs text-gray-400 font-normal">(top = highest in stack)</span></Label>
        <Button type="button" variant="ghost" size="sm"
          onClick={() => set({ layers: [...layers, { label: "", description: "", stage: "Operate" }] })}>
          <Plus size={13} className="mr-1" /> Add layer
        </Button>
      </div>

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={ids} strategy={verticalListSortingStrategy}>
          {layers.map((l, i) => (
            <SortableItem key={i} id={String(i)}>
              <div className="border border-gray-200 rounded-lg p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-gray-500">Layer {i + 1}</span>
                  <button type="button" onClick={() => set({ layers: layers.filter((_, j) => j !== i) })} className="text-gray-400 hover:text-red-500">
                    <X size={14} />
                  </button>
                </div>
                <Input placeholder="Layer label" value={l.label}
                  onChange={(e) => set({ layers: layers.map((x, j) => (j === i ? { ...x, label: e.target.value } : x)) })} />
                <Textarea placeholder="What this layer does" rows={2} value={l.description}
                  onChange={(e) => set({ layers: layers.map((x, j) => (j === i ? { ...x, description: e.target.value } : x)) })} />
                <div className="space-y-1">
                  <Label className="text-xs text-gray-500">Stage / Tier</Label>
                  <Select value={l.stage || "Operate"} onValueChange={(v) => set({ layers: layers.map((x, j) => (j === i ? { ...x, stage: v } : x)) })}>
                    <SelectTrigger className="h-8 text-xs"><SelectValue /></SelectTrigger>
                    <SelectContent>{LAYER_STAGES.map((s) => <SelectItem key={s} value={s}>{s}</SelectItem>)}</SelectContent>
                  </Select>
                </div>
              </div>
            </SortableItem>
          ))}
        </SortableContext>
      </DndContext>

      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <Label className="text-xs text-gray-600">Context Factors (optional sidebar pills)</Label>
          <Button type="button" variant="ghost" size="sm" onClick={() => set({ contextFactors: [...contextFactors, ""] })}>
            <Plus size={13} className="mr-1" /> Add
          </Button>
        </div>
        {contextFactors.map((f, i) => (
          <div key={i} className="flex gap-2">
            <Input value={f} placeholder="e.g. Closing speed"
              onChange={(e) => set({ contextFactors: contextFactors.map((x, j) => (j === i ? e.target.value : x)) })} />
            <button type="button" onClick={() => set({ contextFactors: contextFactors.filter((_, j) => j !== i) })} className="text-gray-400 hover:text-red-500">
              <X size={14} />
            </button>
          </div>
        ))}
      </div>

      <CalloutEditor callout={diagramData.callout} onChange={(c) => set({ callout: c })} />
    </div>
  );
}

// ─── Architecture editor ──────────────────────────────────────────────────────

function ArchitectureEditor({ diagramData, onChange }: { diagramData: any; onChange: (d: any) => void }) {
  const sensors: any[] = diagramData.sensors || [];
  const stages: string[] = diagramData.coreAlgorithm?.stages || [];
  const set = (patch: object) => onChange({ ...diagramData, ...patch });
  const dndSensors = useSensors(useSensor(PointerSensor));
  const sensorIds = sensors.map((_, i) => String(i));
  const stageIds = stages.map((_, i) => `stage-${i}`);

  function handleSensorDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (over && active.id !== over.id) {
      const from = sensorIds.indexOf(String(active.id));
      const to = sensorIds.indexOf(String(over.id));
      set({ sensors: arrayMove(sensors, from, to) });
    }
  }

  function handleStageDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (over && active.id !== over.id) {
      const from = stageIds.indexOf(String(active.id));
      const to = stageIds.indexOf(String(over.id));
      set({ coreAlgorithm: { stages: arrayMove(stages, from, to) } });
    }
  }

  return (
    <div className="space-y-4">
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <Label>Sensor Stack (left column)</Label>
          <Button type="button" variant="ghost" size="sm"
            onClick={() => set({ sensors: [...sensors, { name: "", description: "", functions: [] }] })}>
            <Plus size={13} className="mr-1" /> Add sensor
          </Button>
        </div>
        <DndContext sensors={dndSensors} collisionDetection={closestCenter} onDragEnd={handleSensorDragEnd}>
          <SortableContext items={sensorIds} strategy={verticalListSortingStrategy}>
            {sensors.map((s, i) => (
              <SortableItem key={i} id={String(i)}>
                <div className="border border-gray-200 rounded-lg p-3 space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-medium text-gray-500">Sensor {i + 1}</span>
                    <button type="button" onClick={() => set({ sensors: sensors.filter((_, j) => j !== i) })} className="text-gray-400 hover:text-red-500">
                      <X size={14} />
                    </button>
                  </div>
                  <Input placeholder="Sensor name (e.g. 4D Imaging Radar)" value={s.name}
                    onChange={(e) => set({ sensors: sensors.map((x, j) => (j === i ? { ...x, name: e.target.value } : x)) })} />
                  <Textarea placeholder="What this sensor does" rows={2} value={s.description}
                    onChange={(e) => set({ sensors: sensors.map((x, j) => (j === i ? { ...x, description: e.target.value } : x)) })} />
                  <div className="space-y-1">
                    <Label className="text-xs text-gray-500">Function Tags (comma-separated)</Label>
                    <Input placeholder="e.g. All-weather detection, Velocity measurement"
                      value={(s.functions || []).join(", ")}
                      onChange={(e) => set({
                        sensors: sensors.map((x, j) => j === i ? { ...x, functions: e.target.value.split(",").map((f: string) => f.trim()).filter(Boolean) } : x),
                      })} />
                  </div>
                </div>
              </SortableItem>
            ))}
          </SortableContext>
        </DndContext>
      </div>

      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <Label>Core Algorithm Stages (right column)</Label>
          <Button type="button" variant="ghost" size="sm"
            onClick={() => set({ coreAlgorithm: { stages: [...stages, ""] } })}>
            <Plus size={13} className="mr-1" /> Add stage
          </Button>
        </div>
        <DndContext sensors={dndSensors} collisionDetection={closestCenter} onDragEnd={handleStageDragEnd}>
          <SortableContext items={stageIds} strategy={verticalListSortingStrategy}>
            {stages.map((st, i) => (
              <SortableItem key={i} id={`stage-${i}`}>
                <div className="flex gap-2 items-center">
                  <span className="text-xs text-gray-400 w-5 flex-shrink-0">{i + 1}</span>
                  <Input placeholder="Stage name (e.g. Detect)" value={st}
                    onChange={(e) => {
                      const next = stages.map((x, j) => (j === i ? e.target.value : x));
                      set({ coreAlgorithm: { stages: next } });
                    }} />
                  <button type="button" onClick={() => set({ coreAlgorithm: { stages: stages.filter((_, j) => j !== i) } })} className="text-gray-400 hover:text-red-500 flex-shrink-0">
                    <X size={14} />
                  </button>
                </div>
              </SortableItem>
            ))}
          </SortableContext>
        </DndContext>
      </div>
    </div>
  );
}

// ─── Tree / hierarchy editor ──────────────────────────────────────────────────

function TreeEditor({ diagramData, onChange }: { diagramData: any; onChange: (d: any) => void }) {
  const nodes: any[] = diagramData.nodes || [];
  const set = (patch: object) => onChange({ ...diagramData, ...patch });
  const sensors = useSensors(useSensor(PointerSensor));

  // Build indented display order — BFS order with depth info
  const orderedNodes: { node: any; depth: number }[] = [];
  function buildDisplay(parentId: string | undefined, depth: number) {
    const children = nodes.filter((n) =>
      parentId === undefined ? !n.parent : n.parent === parentId
    );
    for (const child of children) {
      orderedNodes.push({ node: child, depth });
      buildDisplay(child.id, depth + 1);
    }
  }
  buildDisplay(undefined, 0);
  // Nodes not reachable from root (e.g. orphans with invalid parent) appended at end
  const reachableIds = new Set(orderedNodes.map((o) => o.node.id));
  const orphans = nodes.filter((n) => !reachableIds.has(n.id));
  for (const n of orphans) orderedNodes.push({ node: n, depth: 0 });

  const sortableIds = orderedNodes.map((o) => o.node.id);

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const fromVisual = sortableIds.indexOf(String(active.id));
    const toVisual = sortableIds.indexOf(String(over.id));
    // Reorder in display order then write back to flat array
    const reordered = arrayMove(orderedNodes, fromVisual, toVisual);
    set({ nodes: reordered.map((o) => o.node) });
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <Label>
          Nodes
          <span className="text-xs text-gray-400 font-normal ml-1">— drag to reorder, set Parent to build hierarchy</span>
        </Label>
        <Button type="button" variant="ghost" size="sm"
          onClick={() => set({ nodes: [...nodes, { id: `node-${Date.now()}`, label: "", parent: "", description: "", highlight: false }] })}>
          <Plus size={13} className="mr-1" /> Add node
        </Button>
      </div>

      {/* Visual tree outline */}
      {orderedNodes.length > 0 && (
        <div className="border border-gray-100 rounded-lg p-3 bg-gray-50 space-y-1">
          <p className="text-xs text-gray-400 mb-2 font-medium">Tree structure</p>
          {orderedNodes.map(({ node, depth }) => (
            <div
              key={node.id}
              className="flex items-center gap-1.5 text-xs text-gray-600"
              style={{ paddingLeft: `${depth * 16}px` }}
            >
              {depth > 0 && <span className="text-gray-300 flex-shrink-0">└─</span>}
              <span
                className={`px-2 py-0.5 rounded text-xs font-medium ${
                  node.highlight
                    ? "bg-amber-100 text-amber-800 border border-amber-200"
                    : "bg-white border border-gray-200 text-gray-700"
                }`}
              >
                {node.label || <span className="italic text-gray-300">{node.id}</span>}
              </span>
            </div>
          ))}
        </div>
      )}

      {/* Editable node list — in display (tree) order */}
      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={sortableIds} strategy={verticalListSortingStrategy}>
          {orderedNodes.map(({ node, depth }, _vi) => {
            const i = nodes.findIndex((n) => n.id === node.id);
            return (
              <SortableItem key={node.id} id={node.id}>
                <div
                  className="border border-gray-200 rounded-lg p-3 space-y-2"
                  style={{ borderLeftWidth: depth > 0 ? 3 : 1, borderLeftColor: depth > 0 ? "#d1d5db" : undefined }}
                >
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-medium text-gray-500 flex items-center gap-1">
                      {depth > 0 && <span className="text-gray-300">{"└─".slice(0, depth)}  </span>}
                      {node.label || <span className="italic">unlabelled</span>}
                    </span>
                    <button type="button" onClick={() => set({ nodes: nodes.filter((_, j) => j !== i) })} className="text-gray-400 hover:text-red-500">
                      <X size={14} />
                    </button>
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    <div className="space-y-1">
                      <Label className="text-xs text-gray-500">ID (unique)</Label>
                      <Input placeholder="e.g. node-1" value={node.id}
                        onChange={(e) => set({ nodes: nodes.map((x, j) => (j === i ? { ...x, id: e.target.value } : x)) })} />
                    </div>
                    <div className="space-y-1">
                      <Label className="text-xs text-gray-500">Parent</Label>
                      <Select
                        value={node.parent || "__none__"}
                        onValueChange={(v) => set({ nodes: nodes.map((x, j) => (j === i ? { ...x, parent: v === "__none__" ? undefined : v } : x)) })}
                      >
                        <SelectTrigger className="h-8 text-xs"><SelectValue /></SelectTrigger>
                        <SelectContent>
                          <SelectItem value="__none__">— root —</SelectItem>
                          {nodes.filter((_, j) => j !== i).map((other) => (
                            <SelectItem key={other.id} value={other.id}>{other.label || other.id}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                  <Input placeholder="Label (shown in node box)" value={node.label}
                    onChange={(e) => set({ nodes: nodes.map((x, j) => (j === i ? { ...x, label: e.target.value } : x)) })} />
                  <Input placeholder="Description (optional, shown below label)" value={node.description || ""}
                    onChange={(e) => set({ nodes: nodes.map((x, j) => (j === i ? { ...x, description: e.target.value } : x)) })} />
                  <label className="flex items-center gap-2 text-xs text-gray-600 cursor-pointer">
                    <input type="checkbox" checked={!!node.highlight}
                      onChange={(e) => set({ nodes: nodes.map((x, j) => (j === i ? { ...x, highlight: e.target.checked } : x)) })}
                      className="rounded border-gray-300" />
                    Highlight (amber)
                  </label>
                </div>
              </SortableItem>
            );
          })}
        </SortableContext>
      </DndContext>

      <CalloutEditor callout={diagramData.callout} onChange={(c) => set({ callout: c })} />
    </div>
  );
}

// ─── Flow / arrow steps editor ────────────────────────────────────────────────

function FlowEditor({ diagramData, onChange }: { diagramData: any; onChange: (d: any) => void }) {
  const steps: any[] = diagramData.steps || [];
  const set = (patch: object) => onChange({ ...diagramData, ...patch });
  const sensors = useSensors(useSensor(PointerSensor));
  const ids = steps.map((_, i) => String(i));

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (over && active.id !== over.id) {
      const from = ids.indexOf(String(active.id));
      const to = ids.indexOf(String(over.id));
      set({ steps: arrayMove(steps, from, to) });
    }
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <Label>Flow Steps <span className="text-xs text-gray-400 font-normal">(left→right with arrows)</span></Label>
        <Button type="button" variant="ghost" size="sm"
          onClick={() => set({ steps: [...steps, { label: "", description: "", highlight: false }] })}>
          <Plus size={13} className="mr-1" /> Add step
        </Button>
      </div>

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={ids} strategy={verticalListSortingStrategy}>
          {steps.map((s, i) => (
            <SortableItem key={i} id={String(i)}>
              <div className="border border-gray-200 rounded-lg p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-gray-500">Step {i + 1}</span>
                  <button type="button" onClick={() => set({ steps: steps.filter((_, j) => j !== i) })} className="text-gray-400 hover:text-red-500">
                    <X size={14} />
                  </button>
                </div>
                <Input placeholder="Step label" value={s.label || ""}
                  onChange={(e) => set({ steps: steps.map((x, j) => (j === i ? { ...x, label: e.target.value } : x)) })} />
                <Input placeholder="Short description" value={s.description || ""}
                  onChange={(e) => set({ steps: steps.map((x, j) => (j === i ? { ...x, description: e.target.value } : x)) })} />
                <label className="flex items-center gap-2 text-xs text-gray-600 cursor-pointer">
                  <input type="checkbox" checked={!!s.highlight}
                    onChange={(e) => set({ steps: steps.map((x, j) => (j === i ? { ...x, highlight: e.target.checked } : x)) })}
                    className="rounded border-gray-300" />
                  Highlight this step
                </label>
              </div>
            </SortableItem>
          ))}
        </SortableContext>
      </DndContext>

      <CalloutEditor callout={diagramData.callout} onChange={(c) => set({ callout: c })} />
    </div>
  );
}

// ─── Main export ──────────────────────────────────────────────────────────────

export default function DiagramSlideEditor({
  content,
  onChange,
}: {
  content: any;
  onChange: (c: any) => void;
}) {
  const set = (key: string, val: any) => onChange({ ...content, [key]: val });
  const diagramType: string = content.diagramType ?? "";
  const diagramData: any = content.diagramData ?? {};

  const formPanel = (
    <div className="space-y-4">
      <div className="space-y-1.5">
        <Label>Heading</Label>
        <Input value={content.heading || ""} onChange={(e) => set("heading", e.target.value)} placeholder="Slide heading" />
      </div>

      <div className="space-y-1.5">
        <Label>Description</Label>
        <Textarea value={content.description || ""} onChange={(e) => set("description", e.target.value)} placeholder="Optional introductory sentence" rows={2} />
      </div>

      <div className="space-y-1.5">
        <Label>Diagram Type</Label>
        <Select value={diagramType || "__flow__"} onValueChange={(v) => set("diagramType", v === "__flow__" ? "" : v)}>
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            {DIAGRAM_TYPES.map((t) => (
              <SelectItem key={t.value || "__flow__"} value={t.value || "__flow__"}>{t.label}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {diagramType === "kpi-map" && <KpiMapEditor diagramData={diagramData} onChange={(d) => set("diagramData", d)} />}
      {diagramType === "layers" && <LayersEditor diagramData={diagramData} onChange={(d) => set("diagramData", d)} />}
      {diagramType === "architecture" && <ArchitectureEditor diagramData={diagramData} onChange={(d) => set("diagramData", d)} />}
      {diagramType === "tree" && <TreeEditor diagramData={diagramData} onChange={(d) => set("diagramData", d)} />}
      {diagramType === "image" && (
        <div className="space-y-3">
          <ImageUploader label="Diagram Image" value={content.imageUrl || ""} onChange={(url) => set("imageUrl", url)} />
          <p className="text-xs text-gray-400">Upload a screenshot or exported diagram. Displayed full-width.</p>
        </div>
      )}
      {diagramType === "" && <FlowEditor diagramData={diagramData} onChange={(d) => set("diagramData", d)} />}
    </div>
  );

  return (
    <div className="flex gap-6 min-h-0">
      {/* Left: form */}
      <div className="flex-1 min-w-0 overflow-y-auto max-h-[calc(90vh-180px)] pr-1">
        {formPanel}
      </div>

      {/* Divider */}
      <div className="w-px bg-gray-100 flex-shrink-0" />

      {/* Right: live preview */}
      <div className="w-[48%] flex-shrink-0 min-w-0">
        <div className="sticky top-0">
          <div className="flex items-center gap-1.5 text-xs font-medium text-gray-400 uppercase tracking-wide mb-2">
            <Eye size={12} />
            Live Preview
          </div>
          <div className="border border-gray-200 rounded-lg bg-white overflow-y-auto max-h-[calc(90vh-180px)]">
            {content.heading || content.diagramData || content.diagramType ? (
              <div className="p-4 text-sm">
                <DiagramSlide content={content} />
              </div>
            ) : (
              <div className="flex items-center justify-center h-48 text-gray-300 text-sm">
                Start filling in the form to see a preview
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
