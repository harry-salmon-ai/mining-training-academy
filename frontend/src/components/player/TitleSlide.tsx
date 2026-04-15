interface TitleSlideProps {
  content: {
    heading?: string;
    subtitle?: string;
    bullets?: string[];
    backgroundImage?: string;
  };
}

export default function TitleSlide({ content }: TitleSlideProps) {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] text-center space-y-6"
      style={content.backgroundImage ? {
        backgroundImage: `url(${content.backgroundImage})`,
        backgroundSize: "cover",
        backgroundPosition: "center",
      } : {}}
    >
      <h1 className="text-4xl md:text-5xl font-bold text-gray-900">{content.heading}</h1>
      {content.subtitle && (
        <p className="text-xl text-gray-500 max-w-2xl">{content.subtitle}</p>
      )}
      {content.bullets && content.bullets.length > 0 && (
        <ul className="text-left space-y-2 mt-6">
          {content.bullets.map((b, i) => (
            <li key={i} className="flex items-start gap-2 text-gray-600">
              <span className="mt-1.5 h-2 w-2 rounded-full bg-primary-500 flex-shrink-0" />
              {b}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
