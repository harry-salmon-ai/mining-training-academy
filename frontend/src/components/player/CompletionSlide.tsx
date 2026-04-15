import { Award, ArrowRight, CheckCircle2 } from "lucide-react";
import { Link } from "react-router-dom";
import { Button } from "../ui/button";

interface CompletionSlideProps {
  content: {
    heading?: string;
    message?: string;
    keyTakeaways?: string[];
    nextModuleId?: string;
    showCertificate?: boolean;
  };
}

export default function CompletionSlide({ content }: CompletionSlideProps) {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] text-center space-y-6">
      <div className="w-20 h-20 rounded-full bg-green-100 flex items-center justify-center">
        <Award size={40} className="text-green-600" />
      </div>

      <h2 className="text-3xl font-bold">{content.heading || "Module Complete!"}</h2>
      {content.message && <p className="text-gray-600 max-w-lg">{content.message}</p>}

      {content.keyTakeaways && content.keyTakeaways.length > 0 && (
        <div className="bg-gray-50 rounded-lg p-6 max-w-lg w-full text-left">
          <h3 className="font-medium mb-3">Key Takeaways</h3>
          <ul className="space-y-2">
            {content.keyTakeaways.map((t, i) => (
              <li key={i} className="flex items-start gap-2 text-sm text-gray-600">
                <CheckCircle2 size={16} className="text-green-500 mt-0.5 flex-shrink-0" />
                {t}
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="flex gap-3">
        {content.showCertificate && (
          <Link to="/certificates">
            <Button variant="outline">
              <Award size={16} className="mr-2" /> View Certificate
            </Button>
          </Link>
        )}
        <Link to="/modules">
          <Button>
            Browse More Modules <ArrowRight size={16} className="ml-2" />
          </Button>
        </Link>
      </div>
    </div>
  );
}
