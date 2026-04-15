import TitleSlide from "./TitleSlide";
import ContentSlide from "./ContentSlide";
import EquipmentSlide from "./EquipmentSlide";
import DiagramSlide from "./DiagramSlide";
import InteractiveSlide from "./InteractiveSlide";
import ComparisonSlide from "./ComparisonSlide";
import ProcessSlide from "./ProcessSlide";
import QuizSlide from "./QuizSlide";
import VideoSlide from "./VideoSlide";
import CompletionSlide from "./CompletionSlide";

interface SlideRendererProps {
  slide: {
    id: string;
    title: string;
    type: string;
    content: any;
  };
}

export default function SlideRenderer({ slide }: SlideRendererProps) {
  const content = typeof slide.content === "string"
    ? JSON.parse(slide.content)
    : slide.content;

  switch (slide.type) {
    case "TITLE":
      return <TitleSlide content={content} />;
    case "CONTENT":
      return <ContentSlide content={content} />;
    case "EQUIPMENT":
      return <EquipmentSlide content={content} />;
    case "DIAGRAM":
      return <DiagramSlide content={content} />;
    case "INTERACTIVE":
      return <InteractiveSlide content={content} />;
    case "COMPARISON":
      return <ComparisonSlide content={content} />;
    case "PROCESS":
      return <ProcessSlide content={content} />;
    case "QUIZ":
    case "KNOWLEDGE_CHECK":
      return <QuizSlide content={content} />;
    case "VIDEO":
      return <VideoSlide content={content} />;
    case "COMPLETION":
      return <CompletionSlide content={content} />;
    case "IMAGE":
      return (
        <div className="space-y-4">
          <h2 className="text-2xl font-bold">{content.heading || slide.title}</h2>
          {content.description && <p className="text-gray-600">{content.description}</p>}
          {content.imageUrl && <img src={content.imageUrl} alt={content.heading} className="rounded-lg max-w-full" />}
        </div>
      );
    default:
      return <ContentSlide content={content} />;
  }
}
