import { useRef, useState } from "react";
import { Upload, ImageIcon, X } from "lucide-react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";

export default function ImageUploader({
  value,
  onChange,
  label = "Image",
}: {
  value?: string;
  onChange: (url: string) => void;
  label?: string;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState("");

  const handleFile = async (file: File) => {
    setError("");
    setUploading(true);
    try {
      const form = new FormData();
      form.append("file", file);
      const res = await fetch("/api/media/upload", {
        method: "POST",
        body: form,
        credentials: "include",
      });
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.error || "Upload failed");
      }
      const data = await res.json();
      onChange(data.url);
    } catch (e: any) {
      setError(e.message || "Upload failed");
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-1.5">
      <Label>{label}</Label>
      <div className="flex gap-2 items-start">
        <div className="flex-1 space-y-2">
          <Input
            value={value || ""}
            onChange={(e) => onChange(e.target.value)}
            placeholder="/api/media/abc123.jpg or https://..."
            className="text-xs"
          />
          <button
            type="button"
            disabled={uploading}
            onClick={() => inputRef.current?.click()}
            className="flex items-center gap-1.5 text-xs text-primary-600 hover:text-primary-700 disabled:opacity-50"
          >
            <Upload size={12} />
            {uploading ? "Uploading…" : "Upload file"}
          </button>
          {error && <p className="text-xs text-red-500">{error}</p>}
        </div>
        {value ? (
          <div className="relative flex-shrink-0">
            <img
              src={value}
              alt="Preview"
              className="w-20 h-20 object-cover rounded border border-gray-200"
              onError={(e) => (e.currentTarget.style.display = "none")}
            />
            <button
              type="button"
              onClick={() => onChange("")}
              className="absolute -top-1.5 -right-1.5 bg-white border border-gray-300 rounded-full p-0.5 hover:bg-red-50"
            >
              <X size={10} />
            </button>
          </div>
        ) : (
          <div className="w-20 h-20 flex items-center justify-center rounded border-2 border-dashed border-gray-200 text-gray-300 flex-shrink-0">
            <ImageIcon size={20} />
          </div>
        )}
      </div>
      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(e) => {
          const file = e.target.files?.[0];
          if (file) handleFile(file);
          e.target.value = "";
        }}
      />
    </div>
  );
}
