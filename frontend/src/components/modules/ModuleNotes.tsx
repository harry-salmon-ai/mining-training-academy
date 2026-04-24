import { useEffect, useState } from "react";
import { Pencil, Trash2, MessageSquarePlus, X } from "lucide-react";
import { apiFetch } from "../../lib/utils";
import { useAuth } from "../../lib/auth";
import { Button } from "../ui/button";
import { Textarea } from "../ui/textarea";

interface Note {
  id: string;
  userId: string;
  content: string;
  createdAt: string;
  updatedAt: string;
  user: { id: string; name: string; firstName?: string; lastName?: string };
}

interface SlideNotesProps {
  slideId: string;
  onClose: () => void;
}

export default function SlideNotes({ slideId, onClose }: SlideNotesProps) {
  const { user } = useAuth();
  const [notes, setNotes] = useState<Note[]>([]);
  const [newContent, setNewContent] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editContent, setEditContent] = useState("");
  const [adding, setAdding] = useState(false);
  const [saving, setSaving] = useState(false);

  const fetchNotes = () => {
    apiFetch<Note[]>(`/slides/${slideId}/notes`).then(setNotes).catch(() => {});
  };

  useEffect(() => {
    fetchNotes();
    setAdding(false);
    setEditingId(null);
  }, [slideId]);

  const handleAdd = async () => {
    if (!newContent.trim()) return;
    setSaving(true);
    try {
      await apiFetch(`/slides/${slideId}/notes`, {
        method: "POST",
        body: JSON.stringify({ content: newContent.trim() }),
      });
      setNewContent("");
      setAdding(false);
      fetchNotes();
    } catch {}
    setSaving(false);
  };

  const handleUpdate = async (id: string) => {
    if (!editContent.trim()) return;
    setSaving(true);
    try {
      await apiFetch(`/notes/${id}`, {
        method: "PATCH",
        body: JSON.stringify({ content: editContent.trim() }),
      });
      setEditingId(null);
      fetchNotes();
    } catch {}
    setSaving(false);
  };

  const handleDelete = async (id: string) => {
    try {
      await apiFetch(`/notes/${id}`, { method: "DELETE" });
      fetchNotes();
    } catch {}
  };

  const canEdit = (note: Note) =>
    user?.id === note.userId || user?.role === "SUPERADMIN";

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString("en-AU", {
      day: "numeric",
      month: "short",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const displayName = (n: Note) =>
    n.user?.name ||
    [n.user?.firstName, n.user?.lastName].filter(Boolean).join(" ") ||
    "Unknown";

  return (
    <div className="h-full flex flex-col bg-gray-50 border-l">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b bg-white">
        <h3 className="text-sm font-semibold">
          SME Notes{" "}
          {notes.length > 0 && (
            <span className="text-xs font-normal text-gray-400">({notes.length})</span>
          )}
        </h3>
        <Button variant="ghost" size="icon" className="h-7 w-7" onClick={onClose}>
          <X size={16} />
        </Button>
      </div>

      {/* Notes list */}
      <div className="flex-1 overflow-y-auto p-3 space-y-3">
        {notes.length === 0 && !adding && (
          <p className="text-xs text-gray-400 text-center py-6">
            No notes on this slide yet.
          </p>
        )}

        {notes.map((note) => (
          <div key={note.id} className="bg-white border rounded-lg p-3 space-y-1.5">
            {editingId === note.id ? (
              <div className="space-y-2">
                <Textarea
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  className="min-h-[60px] text-sm"
                  autoFocus
                />
                <div className="flex gap-1.5 justify-end">
                  <Button variant="ghost" size="sm" className="h-7 text-xs" onClick={() => setEditingId(null)}>
                    Cancel
                  </Button>
                  <Button size="sm" className="h-7 text-xs" onClick={() => handleUpdate(note.id)} disabled={saving || !editContent.trim()}>
                    Update
                  </Button>
                </div>
              </div>
            ) : (
              <>
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-gray-600">{displayName(note)}</span>
                  {canEdit(note) && (
                    <div className="flex gap-0.5">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0"
                        onClick={() => {
                          setEditingId(note.id);
                          setEditContent(note.content);
                        }}
                      >
                        <Pencil size={12} />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0 text-red-500 hover:text-red-700"
                        onClick={() => handleDelete(note.id)}
                      >
                        <Trash2 size={12} />
                      </Button>
                    </div>
                  )}
                </div>
                <p className="text-sm text-gray-700 whitespace-pre-wrap">{note.content}</p>
                <span className="text-[10px] text-gray-400">
                  {formatDate(note.createdAt)}
                  {note.updatedAt !== note.createdAt && " (edited)"}
                </span>
              </>
            )}
          </div>
        ))}

        {adding && (
          <div className="bg-white border rounded-lg p-3 space-y-2">
            <Textarea
              value={newContent}
              onChange={(e) => setNewContent(e.target.value)}
              placeholder="Add context or suggestions for this slide..."
              className="min-h-[60px] text-sm"
              autoFocus
            />
            <div className="flex gap-1.5 justify-end">
              <Button variant="ghost" size="sm" className="h-7 text-xs" onClick={() => { setAdding(false); setNewContent(""); }}>
                Cancel
              </Button>
              <Button size="sm" className="h-7 text-xs" onClick={handleAdd} disabled={saving || !newContent.trim()}>
                {saving ? "Saving..." : "Save"}
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Add button footer */}
      {!adding && (
        <div className="p-3 border-t bg-white">
          <Button variant="outline" size="sm" className="w-full text-xs" onClick={() => setAdding(true)}>
            <MessageSquarePlus size={14} className="mr-1.5" />
            Add Note
          </Button>
        </div>
      )}
    </div>
  );
}
