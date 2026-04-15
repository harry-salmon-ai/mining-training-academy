import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Plus, GripVertical, Trash2, Edit2, Check, X } from "lucide-react";
import DiagramSlideEditor from "../../components/admin/DiagramSlideEditor";
import ImageUploaderShared from "../../components/admin/ImageUploader";
import { apiFetch } from "../../lib/utils";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Textarea } from "../../components/ui/textarea";
import { Label } from "../../components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";
import { Badge } from "../../components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "../../components/ui/dialog";

const SLIDE_TYPES = [
  "TITLE", "CONTENT", "DIAGRAM", "INTERACTIVE", "EQUIPMENT",
  "COMPARISON", "PROCESS", "QUIZ", "VIDEO", "IMAGE",
  "KNOWLEDGE_CHECK", "COMPLETION",
] as const;

const LEVELS = ["FOUNDATION", "INTERMEDIATE", "ADVANCED", "SPECIALIST"] as const;

// --- Reusable string list editor ---

function StringListEditor({
  label,
  items,
  onChange,
}: {
  label: string;
  items: string[];
  onChange: (items: string[]) => void;
}) {
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between">
        <Label>{label}</Label>
        <Button type="button" variant="ghost" size="sm" onClick={() => onChange([...items, ""])}>
          <Plus size={13} className="mr-1" /> Add
        </Button>
      </div>
      {items.map((item, i) => (
        <div key={i} className="flex gap-2">
          <Input
            value={item}
            onChange={(e) => onChange(items.map((x, j) => (j === i ? e.target.value : x)))}
            placeholder={`${label} ${i + 1}`}
          />
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={() => onChange(items.filter((_, j) => j !== i))}
          >
            <X size={14} />
          </Button>
        </div>
      ))}
    </div>
  );
}

// Use shared ImageUploader
const ImageUploader = ImageUploaderShared;

// --- Per-type slide content editor ---

function SlideContentEditor({
  type,
  content,
  onChange,
}: {
  type: string;
  content: any;
  onChange: (c: any) => void;
}) {
  const set = (key: string, val: any) => onChange({ ...content, [key]: val });

  const headingField = (
    <div className="space-y-1.5">
      <Label>Heading</Label>
      <Input
        value={content.heading || ""}
        onChange={(e) => set("heading", e.target.value)}
        placeholder="Slide heading"
      />
    </div>
  );

  const descriptionField = (
    <div className="space-y-1.5">
      <Label>Description</Label>
      <Textarea
        value={content.description || ""}
        onChange={(e) => set("description", e.target.value)}
        placeholder="Body text or description"
        rows={3}
      />
    </div>
  );

  const updateListItem = (listKey: string, idx: number, patch: object) =>
    set(listKey, (content[listKey] || []).map((x: any, j: number) => (j === idx ? { ...x, ...patch } : x)));

  const removeListItem = (listKey: string, idx: number) =>
    set(listKey, (content[listKey] || []).filter((_: any, j: number) => j !== idx));

  switch (type) {
    case "TITLE":
      return (
        <div className="space-y-4">
          {headingField}
          <div className="space-y-1.5">
            <Label>Subtitle</Label>
            <Input
              value={content.subtitle || ""}
              onChange={(e) => set("subtitle", e.target.value)}
              placeholder="Optional subtitle"
            />
          </div>
          <StringListEditor
            label="Key Points"
            items={content.bullets || []}
            onChange={(v) => set("bullets", v)}
          />
        </div>
      );

    case "CONTENT":
      return (
        <div className="space-y-4">
          {headingField}
          <div className="space-y-1.5">
            <Label>Intro Text</Label>
            <Textarea
              value={content.introText || ""}
              onChange={(e) => set("introText", e.target.value)}
              rows={2}
            />
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label>Content Items</Label>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => set("items", [...(content.items || []), { title: "", description: "" }])}
              >
                <Plus size={13} className="mr-1" /> Add
              </Button>
            </div>
            {(content.items || []).map((item: any, i: number) => (
              <div key={i} className="border border-gray-200 rounded-md p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs text-gray-500 font-medium">Item {i + 1}</span>
                  <Button type="button" variant="ghost" size="icon" onClick={() => removeListItem("items", i)}>
                    <X size={13} />
                  </Button>
                </div>
                <Input
                  placeholder="Title"
                  value={item.title || ""}
                  onChange={(e) => updateListItem("items", i, { title: e.target.value })}
                />
                <Textarea
                  placeholder="Description"
                  value={item.description || ""}
                  rows={2}
                  onChange={(e) => updateListItem("items", i, { description: e.target.value })}
                />
                <ImageUploader
                  label="Item Image (optional)"
                  value={item.imageUrl || ""}
                  onChange={(url) => updateListItem("items", i, { imageUrl: url })}
                />
              </div>
            ))}
          </div>
          <div className="space-y-1.5">
            <Label>Key Fact (optional callout)</Label>
            <Input
              value={content.keyFact || ""}
              onChange={(e) => set("keyFact", e.target.value)}
              placeholder="Highlighted fact shown at bottom of slide"
            />
          </div>
        </div>
      );

    case "VIDEO":
      return (
        <div className="space-y-4">
          {headingField}
          {descriptionField}
          <div className="space-y-1.5">
            <Label>Provider</Label>
            <Select
              value={content.provider || "url"}
              onValueChange={(v) => set("provider", v)}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {["youtube", "vimeo", "url", "upload"].map((p) => (
                  <SelectItem key={p} value={p}>
                    {p.charAt(0).toUpperCase() + p.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-1.5">
            <Label>Video URL</Label>
            <Input
              value={content.videoUrl || ""}
              onChange={(e) => set("videoUrl", e.target.value)}
              placeholder="https://..."
            />
          </div>
          <div className="space-y-1.5">
            <Label>Thumbnail URL</Label>
            <Input
              value={content.thumbnailUrl || ""}
              onChange={(e) => set("thumbnailUrl", e.target.value)}
              placeholder="https://..."
            />
          </div>
        </div>
      );

    case "QUIZ":
    case "KNOWLEDGE_CHECK":
      return (
        <div className="space-y-4">
          {headingField}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label>Questions</Label>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() =>
                  set("inlineQuestions", [
                    ...(content.inlineQuestions || []),
                    { question: "", options: ["", "", "", ""], correctIndex: 0, explanation: "" },
                  ])
                }
              >
                <Plus size={13} className="mr-1" /> Add Question
              </Button>
            </div>
            {(content.inlineQuestions || []).map((q: any, i: number) => (
              <div key={i} className="border border-gray-200 rounded-md p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-gray-500">Q{i + 1}</span>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => removeListItem("inlineQuestions", i)}
                  >
                    <X size={13} />
                  </Button>
                </div>
                <Textarea
                  placeholder="Question text"
                  value={q.question || ""}
                  rows={2}
                  onChange={(e) => updateListItem("inlineQuestions", i, { question: e.target.value })}
                />
                <div className="grid grid-cols-2 gap-2">
                  {(q.options || ["", "", "", ""]).map((opt: string, oi: number) => (
                    <div key={oi} className="flex items-center gap-1.5">
                      <input
                        type="radio"
                        name={`correct-${i}`}
                        checked={q.correctIndex === oi}
                        onChange={() => updateListItem("inlineQuestions", i, { correctIndex: oi })}
                        className="flex-shrink-0"
                      />
                      <Input
                        placeholder={`Option ${oi + 1}`}
                        value={opt}
                        className="h-7 text-xs"
                        onChange={(e) => {
                          const newOpts = [...(q.options || [])];
                          newOpts[oi] = e.target.value;
                          updateListItem("inlineQuestions", i, { options: newOpts });
                        }}
                      />
                    </div>
                  ))}
                </div>
                <div className="space-y-1">
                  <Label className="text-xs text-gray-500">Explanation (shown after answer)</Label>
                  <Input
                    placeholder="Why is this the correct answer?"
                    value={q.explanation || ""}
                    onChange={(e) => updateListItem("inlineQuestions", i, { explanation: e.target.value })}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      );

    case "EQUIPMENT":
      return (
        <div className="space-y-4">
          {headingField}
          {descriptionField}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label>Specifications</Label>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => set("specifications", [...(content.specifications || []), { label: "", value: "" }])}
              >
                <Plus size={13} className="mr-1" /> Add
              </Button>
            </div>
            {(content.specifications || []).map((spec: any, i: number) => (
              <div key={i} className="flex gap-2">
                <Input
                  placeholder="Label"
                  value={spec.label || ""}
                  onChange={(e) => updateListItem("specifications", i, { label: e.target.value })}
                />
                <Input
                  placeholder="Value"
                  value={spec.value || ""}
                  onChange={(e) => updateListItem("specifications", i, { value: e.target.value })}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => removeListItem("specifications", i)}
                >
                  <X size={13} />
                </Button>
              </div>
            ))}
          </div>
          <StringListEditor
            label="Advantages"
            items={content.advantages || []}
            onChange={(v) => set("advantages", v)}
          />
          <StringListEditor
            label="Limitations"
            items={content.limitations || []}
            onChange={(v) => set("limitations", v)}
          />
        </div>
      );

    case "PROCESS":
      return (
        <div className="space-y-4">
          {headingField}
          {descriptionField}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label>Steps</Label>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() =>
                  set("steps", [...(content.steps || []), { title: "", description: "", duration: "" }])
                }
              >
                <Plus size={13} className="mr-1" /> Add Step
              </Button>
            </div>
            {(content.steps || []).map((step: any, i: number) => (
              <div key={i} className="border border-gray-200 rounded-md p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-gray-500">Step {i + 1}</span>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => removeListItem("steps", i)}
                  >
                    <X size={13} />
                  </Button>
                </div>
                <Input
                  placeholder="Title"
                  value={step.title || ""}
                  onChange={(e) => updateListItem("steps", i, { title: e.target.value })}
                />
                <Textarea
                  placeholder="Description"
                  value={step.description || ""}
                  rows={2}
                  onChange={(e) => updateListItem("steps", i, { description: e.target.value })}
                />
                <Input
                  placeholder="Duration (e.g. 5 min)"
                  value={step.duration || ""}
                  onChange={(e) => updateListItem("steps", i, { duration: e.target.value })}
                />
              </div>
            ))}
          </div>
        </div>
      );

    case "COMPARISON":
      return (
        <div className="space-y-4">
          {headingField}
          {descriptionField}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label>Comparison Items</Label>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() =>
                  set("items", [
                    ...(content.items || []),
                    { name: "", description: "", pros: [], cons: [] },
                  ])
                }
              >
                <Plus size={13} className="mr-1" /> Add Item
              </Button>
            </div>
            {(content.items || []).map((item: any, i: number) => (
              <div key={i} className="border border-gray-200 rounded-md p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-gray-500">Item {i + 1}</span>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => removeListItem("items", i)}
                  >
                    <X size={13} />
                  </Button>
                </div>
                <Input
                  placeholder="Name"
                  value={item.name || ""}
                  onChange={(e) => updateListItem("items", i, { name: e.target.value })}
                />
                <Textarea
                  placeholder="Description"
                  value={item.description || ""}
                  rows={2}
                  onChange={(e) => updateListItem("items", i, { description: e.target.value })}
                />
                <div className="grid grid-cols-2 gap-2">
                  <div className="space-y-1">
                    <Label className="text-xs text-gray-500">Pros (one per line)</Label>
                    <Textarea
                      rows={3}
                      value={(item.pros || []).join("\n")}
                      onChange={(e) =>
                        updateListItem("items", i, {
                          pros: e.target.value.split("\n").filter(Boolean),
                        })
                      }
                    />
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs text-gray-500">Cons (one per line)</Label>
                    <Textarea
                      rows={3}
                      value={(item.cons || []).join("\n")}
                      onChange={(e) =>
                        updateListItem("items", i, {
                          cons: e.target.value.split("\n").filter(Boolean),
                        })
                      }
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      );

    case "DIAGRAM":
      return <DiagramSlideEditor content={content} onChange={onChange} />;

    case "COMPLETION":
      return (
        <div className="space-y-4">
          {headingField}
          <div className="space-y-1.5">
            <Label>Completion Message</Label>
            <Textarea
              value={content.message || ""}
              onChange={(e) => set("message", e.target.value)}
              rows={3}
            />
          </div>
          <StringListEditor
            label="Key Takeaways"
            items={content.keyTakeaways || []}
            onChange={(v) => set("keyTakeaways", v)}
          />
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="showCert"
              checked={!!content.showCertificate}
              onChange={(e) => set("showCertificate", e.target.checked)}
              className="rounded border-gray-300"
            />
            <Label htmlFor="showCert">Show certificate option</Label>
          </div>
        </div>
      );

    case "IMAGE":
      return (
        <div className="space-y-4">
          {headingField}
          {descriptionField}
          <ImageUploader
            label="Slide Image"
            value={content.imageUrl || ""}
            onChange={(url) => set("imageUrl", url)}
          />
        </div>
      );

    default:
      return (
        <div className="space-y-4">
          {headingField}
          {descriptionField}
        </div>
      );
  }
}

// --- Slide editor dialog ---

interface SlideEditorDialogProps {
  open: boolean;
  onClose: () => void;
  sectionId: string;
  slide: any | null;
  onSaved: () => void;
}

function SlideEditorDialog({ open, onClose, sectionId, slide, onSaved }: SlideEditorDialogProps) {
  const isNew = !slide;
  const [title, setTitle] = useState("");
  const [type, setType] = useState("CONTENT");
  const [content, setContent] = useState<any>({});
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (open) {
      setTitle(slide?.title || "");
      setType(slide?.type || "CONTENT");
      const c = slide?.content;
      setContent(c ? (typeof c === "string" ? JSON.parse(c) : c) : {});
      setSaving(false);
    }
  }, [open, slide]);

  const handleSave = async () => {
    if (!title.trim()) return;
    setSaving(true);
    try {
      if (isNew) {
        await apiFetch(`/sections/${sectionId}/slides`, {
          method: "POST",
          body: JSON.stringify({ title: title.trim(), type, content, sortOrder: 999 }),
        });
      } else {
        await apiFetch(`/slides/${slide.id}`, {
          method: "PATCH",
          body: JSON.stringify({ title: title.trim(), content }),
        });
      }
      onSaved();
      onClose();
    } catch {
      setSaving(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(o) => !o && onClose()}>
      <DialogContent className={`${type === "DIAGRAM" ? "max-w-5xl" : "max-w-2xl"} max-h-[90vh] overflow-y-auto`}>
        <DialogHeader>
          <DialogTitle>{isNew ? "Add Slide" : "Edit Slide"}</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label>Slide Title</Label>
              <Input
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="e.g. Introduction to Blast Patterns"
                autoFocus
              />
            </div>
            {isNew ? (
              <div className="space-y-1.5">
                <Label>Type</Label>
                <Select value={type} onValueChange={setType}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {SLIDE_TYPES.map((t) => (
                      <SelectItem key={t} value={t}>
                        {t.charAt(0) + t.slice(1).toLowerCase().replace(/_/g, " ")}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            ) : (
              <div className="space-y-1.5">
                <Label>Type</Label>
                <div className="flex items-center h-9">
                  <Badge variant="secondary" className="text-xs">
                    {type.toLowerCase()}
                  </Badge>
                </div>
              </div>
            )}
          </div>

          <div className="border-t border-gray-100 pt-4">
            <SlideContentEditor type={type} content={content} onChange={setContent} />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={saving || !title.trim()}>
            {saving ? "Saving..." : isNew ? "Add Slide" : "Save Changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// --- Main builder page ---

export default function ModuleBuilderPage() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [module, setModule] = useState<any>(null);
  const [sections, setSections] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [meta, setMeta] = useState<any>({});
  const [metaSaving, setMetaSaving] = useState(false);
  const [slideEditor, setSlideEditor] = useState<{ sectionId: string; slide: any | null } | null>(null);
  const [dragSectionIdx, setDragSectionIdx] = useState<number | null>(null);
  const [dragOverSectionIdx, setDragOverSectionIdx] = useState<number | null>(null);
  const [dragSlide, setDragSlide] = useState<{ sectionIdx: number; slideIdx: number } | null>(null);
  const [dragOverSlide, setDragOverSlide] = useState<{ sectionIdx: number; slideIdx: number } | null>(null);
  const [editSectionId, setEditSectionId] = useState<string | null>(null);
  const [editSectionTitle, setEditSectionTitle] = useState("");
  const [newSectionTitle, setNewSectionTitle] = useState("");

  const load = () => {
    if (!id) return;
    apiFetch<any>(`/modules/${id}`)
      .then((m) => {
        setModule(m);
        setSections(m.sections || []);
        setMeta({
          title: m.title || "",
          description: m.description || "",
          categoryId: m.category?.id || "",
          level: m.level || "FOUNDATION",
          duration: m.duration?.toString() || "",
        });
      })
      .catch(() => navigate("/admin/modules"));
  };

  useEffect(() => {
    load();
    apiFetch<any[]>("/categories?active=true").then(setCategories).catch(() => {});
  }, [id]);

  const saveMeta = async () => {
    setMetaSaving(true);
    try {
      await apiFetch(`/modules/${id}`, {
        method: "PATCH",
        body: JSON.stringify({
          title: meta.title,
          description: meta.description,
          category_id: meta.categoryId || undefined,
          level: meta.level,
          duration: meta.duration ? parseInt(meta.duration) : null,
        }),
      });
      load();
    } finally {
      setMetaSaving(false);
    }
  };

  const addSection = async () => {
    if (!newSectionTitle.trim()) return;
    await apiFetch(`/modules/${id}/sections`, {
      method: "POST",
      body: JSON.stringify({ title: newSectionTitle.trim(), sortOrder: sections.length }),
    });
    setNewSectionTitle("");
    load();
  };

  const saveSection = async (sectionId: string) => {
    if (!editSectionTitle.trim()) return;
    await apiFetch(`/sections/${sectionId}`, {
      method: "PATCH",
      body: JSON.stringify({ title: editSectionTitle.trim() }),
    });
    setEditSectionId(null);
    load();
  };

  const deleteSection = async (sectionId: string) => {
    if (!confirm("Delete this section and all its slides?")) return;
    await apiFetch(`/sections/${sectionId}`, { method: "DELETE" });
    load();
  };

  const deleteSlide = async (slideId: string) => {
    if (!confirm("Delete this slide?")) return;
    await apiFetch(`/slides/${slideId}`, { method: "DELETE" });
    load();
  };

  // Section drag-and-drop
  const onSectionDragStart = (e: React.DragEvent, idx: number) => {
    setDragSectionIdx(idx);
    e.dataTransfer.effectAllowed = "move";
  };

  const onSectionDragOver = (e: React.DragEvent, idx: number) => {
    e.preventDefault();
    setDragOverSectionIdx(idx);
  };

  const onSectionDrop = async (e: React.DragEvent, targetIdx: number) => {
    e.preventDefault();
    setDragOverSectionIdx(null);
    if (dragSectionIdx === null || dragSectionIdx === targetIdx) {
      setDragSectionIdx(null);
      return;
    }
    const next = [...sections];
    const [moved] = next.splice(dragSectionIdx, 1);
    next.splice(targetIdx, 0, moved);
    setSections(next);
    setDragSectionIdx(null);
    await Promise.all(
      next.map((s, i) =>
        apiFetch(`/sections/${s.id}`, { method: "PATCH", body: JSON.stringify({ sort_order: i }) })
      )
    );
  };

  // Slide drag-and-drop (within same section)
  const onSlideDragStart = (e: React.DragEvent, sectionIdx: number, slideIdx: number) => {
    e.stopPropagation();
    setDragSlide({ sectionIdx, slideIdx });
    e.dataTransfer.effectAllowed = "move";
  };

  const onSlideDragOver = (e: React.DragEvent, sectionIdx: number, slideIdx: number) => {
    e.preventDefault();
    e.stopPropagation();
    setDragOverSlide({ sectionIdx, slideIdx });
  };

  const onSlideDrop = async (e: React.DragEvent, targetSectionIdx: number, targetSlideIdx: number) => {
    e.preventDefault();
    e.stopPropagation();
    setDragOverSlide(null);
    if (!dragSlide) return;
    const { sectionIdx, slideIdx } = dragSlide;
    setDragSlide(null);
    if (sectionIdx !== targetSectionIdx || slideIdx === targetSlideIdx) return;
    const next = [...sections];
    const slides = [...(next[sectionIdx].slides || [])];
    const [moved] = slides.splice(slideIdx, 1);
    slides.splice(targetSlideIdx, 0, moved);
    next[sectionIdx] = { ...next[sectionIdx], slides };
    setSections(next);
    await Promise.all(
      slides.map((s: any, i: number) =>
        apiFetch(`/slides/${s.id}`, { method: "PATCH", body: JSON.stringify({ sort_order: i }) })
      )
    );
  };

  if (!module) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Module Builder</h1>
          <p className="text-sm text-gray-500">{module.title}</p>
        </div>
        <Button variant="outline" onClick={() => navigate("/admin/modules")}>
          Back to Modules
        </Button>
      </div>

      {/* Module metadata */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Module Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5 col-span-2">
              <Label>Title</Label>
              <Input
                value={meta.title || ""}
                onChange={(e) => setMeta({ ...meta, title: e.target.value })}
              />
            </div>
            <div className="space-y-1.5 col-span-2">
              <Label>Description</Label>
              <Textarea
                value={meta.description || ""}
                onChange={(e) => setMeta({ ...meta, description: e.target.value })}
                rows={2}
              />
            </div>
            <div className="space-y-1.5">
              <Label>Category</Label>
              <Select
                value={meta.categoryId || ""}
                onValueChange={(v) => setMeta({ ...meta, categoryId: v })}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select category" />
                </SelectTrigger>
                <SelectContent>
                  {categories.map((c) => (
                    <SelectItem key={c.id} value={c.id}>
                      {c.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1.5">
              <Label>Level</Label>
              <Select
                value={meta.level || "FOUNDATION"}
                onValueChange={(v) => setMeta({ ...meta, level: v })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {LEVELS.map((l) => (
                    <SelectItem key={l} value={l}>
                      {l.charAt(0) + l.slice(1).toLowerCase()}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1.5">
              <Label>Duration (minutes)</Label>
              <Input
                type="number"
                min="1"
                value={meta.duration || ""}
                onChange={(e) => setMeta({ ...meta, duration: e.target.value })}
                className="max-w-[160px]"
              />
            </div>
          </div>
          <div className="flex justify-end">
            <Button onClick={saveMeta} disabled={metaSaving}>
              {metaSaving ? "Saving..." : "Save Details"}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Sections */}
      <div className="space-y-3">
        {sections.map((section, si) => (
          <Card
            key={section.id}
            draggable
            onDragStart={(e) => onSectionDragStart(e, si)}
            onDragOver={(e) => onSectionDragOver(e, si)}
            onDrop={(e) => onSectionDrop(e, si)}
            onDragEnd={() => {
              setDragSectionIdx(null);
              setDragOverSectionIdx(null);
            }}
            className={[
              "transition-opacity",
              dragSectionIdx === si ? "opacity-40" : "",
              dragOverSectionIdx === si && dragSectionIdx !== si
                ? "ring-2 ring-gray-400 ring-offset-1"
                : "",
            ]
              .filter(Boolean)
              .join(" ")}
          >
            <CardHeader className="flex flex-row items-center justify-between py-3">
              <div className="flex items-center gap-2 flex-1 min-w-0">
                <GripVertical
                  size={16}
                  className="text-gray-400 cursor-grab flex-shrink-0"
                />
                {editSectionId === section.id ? (
                  <div className="flex items-center gap-2 flex-1">
                    <Input
                      value={editSectionTitle}
                      onChange={(e) => setEditSectionTitle(e.target.value)}
                      onKeyDown={(e) => {
                        if (e.key === "Enter") saveSection(section.id);
                        if (e.key === "Escape") setEditSectionId(null);
                      }}
                      className="h-7 text-sm"
                      autoFocus
                    />
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => saveSection(section.id)}
                    >
                      <Check size={14} />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setEditSectionId(null)}
                    >
                      <X size={14} />
                    </Button>
                  </div>
                ) : (
                  <button
                    className="text-sm font-semibold text-left truncate hover:text-gray-600 transition-colors"
                    onClick={() => {
                      setEditSectionId(section.id);
                      setEditSectionTitle(section.title);
                    }}
                  >
                    Section {si + 1}: {section.title}
                  </button>
                )}
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => deleteSection(section.id)}
                className="flex-shrink-0"
              >
                <Trash2 size={16} className="text-red-500" />
              </Button>
            </CardHeader>
            <CardContent className="space-y-2 pt-0">
              {section.slides?.map((slide: any, sli: number) => {
                const isDragging =
                  dragSlide?.sectionIdx === si && dragSlide?.slideIdx === sli;
                const isDragOver =
                  dragOverSlide?.sectionIdx === si &&
                  dragOverSlide?.slideIdx === sli &&
                  !isDragging;
                return (
                  <div
                    key={slide.id}
                    draggable
                    onDragStart={(e) => onSlideDragStart(e, si, sli)}
                    onDragOver={(e) => onSlideDragOver(e, si, sli)}
                    onDrop={(e) => onSlideDrop(e, si, sli)}
                    onDragEnd={() => {
                      setDragSlide(null);
                      setDragOverSlide(null);
                    }}
                    className={[
                      "flex items-center justify-between px-3 py-2 bg-gray-50 rounded border text-sm transition-opacity",
                      isDragging ? "opacity-40" : "",
                      isDragOver ? "ring-2 ring-gray-400 ring-offset-1" : "",
                    ]
                      .filter(Boolean)
                      .join(" ")}
                  >
                    <div className="flex items-center gap-2 min-w-0">
                      <GripVertical
                        size={14}
                        className="text-gray-400 cursor-grab flex-shrink-0"
                      />
                      <span className="truncate">{slide.title}</span>
                      <Badge variant="secondary" className="text-xs flex-shrink-0">
                        {slide.type?.toLowerCase()}
                      </Badge>
                    </div>
                    <div className="flex items-center gap-1 flex-shrink-0">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() =>
                          setSlideEditor({ sectionId: section.id, slide })
                        }
                      >
                        <Edit2 size={14} />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => deleteSlide(slide.id)}
                      >
                        <Trash2 size={14} className="text-red-400" />
                      </Button>
                    </div>
                  </div>
                );
              })}
              <Button
                variant="outline"
                size="sm"
                onClick={() =>
                  setSlideEditor({ sectionId: section.id, slide: null })
                }
              >
                <Plus size={14} className="mr-1" /> Add Slide
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Add section */}
      <div className="flex items-center gap-2">
        <Input
          placeholder="New section title..."
          value={newSectionTitle}
          onChange={(e) => setNewSectionTitle(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && addSection()}
          className="max-w-xs"
        />
        <Button onClick={addSection}>
          <Plus size={16} className="mr-1" /> Add Section
        </Button>
      </div>

      {/* Slide editor dialog */}
      {slideEditor && (
        <SlideEditorDialog
          open={!!slideEditor}
          onClose={() => setSlideEditor(null)}
          sectionId={slideEditor.sectionId}
          slide={slideEditor.slide}
          onSaved={load}
        />
      )}
    </div>
  );
}
