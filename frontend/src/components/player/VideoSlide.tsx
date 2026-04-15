interface VideoSlideProps {
  content: {
    heading?: string;
    description?: string;
    videoUrl?: string;
    provider?: "youtube" | "vimeo" | "upload" | "url";
    thumbnailUrl?: string;
    timestamps?: { time: number; label: string }[];
  };
}

function getEmbedUrl(url: string, provider?: string): string {
  if (provider === "youtube" || url.includes("youtube.com") || url.includes("youtu.be")) {
    const id = url.match(/(?:v=|\/)([\w-]{11})/)?.[1];
    return id ? `https://www.youtube.com/embed/${id}` : url;
  }
  if (provider === "vimeo" || url.includes("vimeo.com")) {
    const id = url.match(/vimeo\.com\/(\d+)/)?.[1];
    return id ? `https://player.vimeo.com/video/${id}` : url;
  }
  return url;
}

export default function VideoSlide({ content }: VideoSlideProps) {
  const embedUrl = content.videoUrl ? getEmbedUrl(content.videoUrl, content.provider) : "";

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">{content.heading}</h2>
      {content.description && <p className="text-gray-600">{content.description}</p>}

      {content.videoUrl && (
        <div className="aspect-video bg-black rounded-lg overflow-hidden">
          {content.provider === "upload" || content.provider === "url" ? (
            <video src={content.videoUrl} controls className="w-full h-full" poster={content.thumbnailUrl} />
          ) : (
            <iframe
              src={embedUrl}
              className="w-full h-full"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowFullScreen
            />
          )}
        </div>
      )}

      {content.timestamps && content.timestamps.length > 0 && (
        <div>
          <h3 className="font-medium text-sm text-gray-700 mb-2">Timestamps</h3>
          <div className="space-y-1">
            {content.timestamps.map((ts, i) => (
              <div key={i} className="text-sm flex gap-2">
                <span className="font-mono text-primary-600">
                  {Math.floor(ts.time / 60)}:{String(ts.time % 60).padStart(2, "0")}
                </span>
                <span className="text-gray-600">{ts.label}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
