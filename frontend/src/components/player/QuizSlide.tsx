import { useState } from "react";
import { CheckCircle2, XCircle } from "lucide-react";
import { Button } from "../ui/button";

interface QuizSlideProps {
  content: {
    heading?: string;
    quizId?: string;
    inlineQuestions?: {
      question: string;
      options: string[];
      correctIndex: number;
      explanation?: string;
    }[];
  };
}

export default function QuizSlide({ content }: QuizSlideProps) {
  const [answers, setAnswers] = useState<Record<number, number>>({});
  const [submitted, setSubmitted] = useState(false);

  const questions = content.inlineQuestions || [];

  const handleSubmit = () => {
    setSubmitted(true);
  };

  const correctCount = questions.filter((q, i) => answers[i] === q.correctIndex).length;

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">{content.heading || "Knowledge Check"}</h2>

      {questions.map((q, qi) => {
        const userAnswer = answers[qi];
        const isCorrect = userAnswer === q.correctIndex;

        return (
          <div key={qi} className="border rounded-lg p-4 space-y-3">
            <p className="font-medium">{q.question}</p>
            <div className="space-y-2">
              {q.options.map((opt, oi) => {
                let optClass = "border rounded-md px-4 py-2 text-sm cursor-pointer transition";
                if (submitted) {
                  if (oi === q.correctIndex) optClass += " bg-green-50 border-green-500 text-green-800";
                  else if (oi === userAnswer) optClass += " bg-red-50 border-red-500 text-red-800";
                  else optClass += " opacity-50";
                } else if (userAnswer === oi) {
                  optClass += " bg-primary-50 border-primary-500";
                } else {
                  optClass += " hover:bg-gray-50";
                }

                return (
                  <button
                    key={oi}
                    onClick={() => !submitted && setAnswers({ ...answers, [qi]: oi })}
                    className={`block w-full text-left ${optClass}`}
                    disabled={submitted}
                  >
                    <span className="flex items-center gap-2">
                      {submitted && oi === q.correctIndex && <CheckCircle2 size={16} className="text-green-500" />}
                      {submitted && oi === userAnswer && oi !== q.correctIndex && <XCircle size={16} className="text-red-500" />}
                      {opt}
                    </span>
                  </button>
                );
              })}
            </div>
            {submitted && q.explanation && (
              <div className={`text-sm p-3 rounded ${isCorrect ? "bg-green-50 text-green-700" : "bg-yellow-50 text-yellow-700"}`}>
                {q.explanation}
              </div>
            )}
          </div>
        );
      })}

      {questions.length > 0 && !submitted && (
        <Button onClick={handleSubmit} disabled={Object.keys(answers).length < questions.length}>
          Check Answers
        </Button>
      )}

      {submitted && (
        <div className="text-center p-4 bg-gray-50 rounded-lg">
          <p className="text-lg font-bold">
            Score: {correctCount} / {questions.length}
          </p>
          <p className="text-sm text-gray-500 mt-1">
            {correctCount === questions.length ? "Perfect score!" : "Review the explanations above."}
          </p>
        </div>
      )}
    </div>
  );
}
