import { useState } from "react";

interface ClassifyScenario {
  event: string;
  answer: string;
  explanation: string;
}

interface CalculationQuestion {
  question: string;
  answer: string;
  formula: string;
}

interface InteractiveSlideProps {
  content: {
    heading?: string;
    description?: string;
    instructions?: string;
    interactionType?: string;
    interactiveType?: string;
    controls?: {
      id: string;
      label: string;
      type: "range" | "toggle" | "select";
      min?: number;
      max?: number;
      step?: number;
      default?: number | boolean | string;
      options?: { label: string; value: string }[];
    }[];
    outputs?: {
      id: string;
      label: string;
      formula?: string;
      unit?: string;
      format?: string;
    }[];
    feedback?: {
      condition: string;
      message: string;
      type: "success" | "warning" | "error";
    }[];
    scenarios?: ClassifyScenario[];
    calculationQuestion?: CalculationQuestion;
  };
}

function ClassifyExercise({
  instructions,
  scenarios,
  calculationQuestion,
}: {
  instructions?: string;
  scenarios: ClassifyScenario[];
  calculationQuestion?: CalculationQuestion;
}) {
  const categories = [...new Set(scenarios.map((s) => s.answer))];
  const [selections, setSelections] = useState<Record<number, string>>({});
  const [revealed, setRevealed] = useState<Record<number, boolean>>({});
  const [calcAnswer, setCalcAnswer] = useState("");
  const [calcRevealed, setCalcRevealed] = useState(false);

  const handleSelect = (index: number, category: string) => {
    if (revealed[index]) return;
    setSelections((prev) => ({ ...prev, [index]: category }));
    setRevealed((prev) => ({ ...prev, [index]: true }));
  };

  const isCorrect = (index: number) => selections[index] === scenarios[index].answer;

  return (
    <div className="space-y-6">
      {instructions && <p className="text-gray-600 italic">{instructions}</p>}

      <div className="space-y-4">
        {scenarios.map((scenario, i) => (
          <div
            key={i}
            className={`border-2 rounded-xl p-5 transition-colors ${
              revealed[i]
                ? isCorrect(i)
                  ? "border-green-300 bg-green-50/50"
                  : "border-red-300 bg-red-50/50"
                : "border-gray-200 bg-white"
            }`}
          >
            <p className="font-medium text-gray-900 mb-3">{scenario.event}</p>

            <div className="flex flex-wrap gap-2 mb-3">
              {categories.map((cat) => (
                <button
                  key={cat}
                  onClick={() => handleSelect(i, cat)}
                  disabled={revealed[i]}
                  className={`text-sm px-4 py-2 rounded-lg border-2 font-medium transition-all ${
                    revealed[i] && cat === scenario.answer
                      ? "border-green-500 bg-green-100 text-green-800"
                      : revealed[i] && selections[i] === cat && cat !== scenario.answer
                        ? "border-red-400 bg-red-100 text-red-800"
                        : selections[i] === cat && !revealed[i]
                          ? "border-blue-400 bg-blue-50 text-blue-800"
                          : "border-gray-200 bg-white text-gray-700 hover:border-gray-300 hover:bg-gray-50"
                  } ${revealed[i] ? "cursor-default" : "cursor-pointer"}`}
                >
                  {cat}
                </button>
              ))}
            </div>

            {revealed[i] && (
              <div
                className={`text-sm rounded-lg p-3 ${
                  isCorrect(i) ? "bg-green-100 text-green-800" : "bg-amber-100 text-amber-800"
                }`}
              >
                <span className="font-semibold">{isCorrect(i) ? "Correct!" : "Not quite."}</span>{" "}
                {scenario.explanation}
              </div>
            )}
          </div>
        ))}
      </div>

      {calculationQuestion && (
        <div className="border-2 border-blue-200 bg-blue-50/50 rounded-xl p-5 space-y-3">
          <h4 className="font-bold text-gray-900">Calculation Challenge</h4>
          <p className="text-gray-700">{calculationQuestion.question}</p>
          <div className="flex gap-3 items-center">
            <input
              type="text"
              value={calcAnswer}
              onChange={(e) => setCalcAnswer(e.target.value)}
              placeholder="Your answer..."
              className="border border-gray-300 rounded-lg px-4 py-2 text-sm flex-1 max-w-xs"
              disabled={calcRevealed}
            />
            <button
              onClick={() => setCalcRevealed(true)}
              disabled={calcRevealed}
              className="px-4 py-2 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-default transition-colors"
            >
              Check
            </button>
          </div>
          {calcRevealed && (
            <div className="bg-white border border-blue-200 rounded-lg p-4 space-y-1">
              <p className="font-semibold text-gray-900">Answer: {calculationQuestion.answer}</p>
              <p className="text-sm text-gray-600">
                <span className="font-medium">Working:</span>{" "}
                <code className="bg-gray-100 px-2 py-0.5 rounded text-sm">{calculationQuestion.formula}</code>
              </p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default function InteractiveSlide({ content }: InteractiveSlideProps) {
  const [values, setValues] = useState<Record<string, number | boolean | string>>(() => {
    const initial: Record<string, number | boolean | string> = {};
    content.controls?.forEach((c) => {
      initial[c.id] = c.default ?? (c.type === "range" ? c.min ?? 0 : c.type === "toggle" ? false : "");
    });
    return initial;
  });

  const effectiveType = content.interactiveType ?? content.interactionType;

  if (effectiveType === "classify" && content.scenarios) {
    return (
      <div className="space-y-6">
        <h2 className="text-2xl font-bold">{content.heading}</h2>
        {content.description && <p className="text-gray-600">{content.description}</p>}
        <ClassifyExercise
          instructions={content.instructions}
          scenarios={content.scenarios}
          calculationQuestion={content.calculationQuestion}
        />
      </div>
    );
  }

  const evaluateFormula = (formula: string) => {
    try {
      const fn = new Function(...Object.keys(values), `return ${formula}`);
      return fn(...Object.values(values));
    } catch {
      return "N/A";
    }
  };

  const feedbackColors = {
    success: "bg-green-50 text-green-800 border-green-200",
    warning: "bg-yellow-50 text-yellow-800 border-yellow-200",
    error: "bg-red-50 text-red-800 border-red-200",
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">{content.heading}</h2>
      {content.description && <p className="text-gray-600">{content.description}</p>}

      <div className="grid md:grid-cols-2 gap-6">
        <div className="space-y-4">
          <h3 className="font-medium text-gray-700">Controls</h3>
          {content.controls?.map((ctrl) => (
            <div key={ctrl.id} className="space-y-1">
              <label className="text-sm font-medium flex justify-between">
                {ctrl.label}
                {ctrl.type === "range" && <span className="text-primary-600">{values[ctrl.id]}</span>}
              </label>
              {ctrl.type === "range" && (
                <input
                  type="range"
                  min={ctrl.min}
                  max={ctrl.max}
                  step={ctrl.step || 1}
                  value={values[ctrl.id] as number}
                  onChange={(e) => setValues({ ...values, [ctrl.id]: Number(e.target.value) })}
                  className="w-full accent-primary-500"
                />
              )}
              {ctrl.type === "toggle" && (
                <button
                  onClick={() => setValues({ ...values, [ctrl.id]: !values[ctrl.id] })}
                  className={`w-12 h-6 rounded-full transition ${values[ctrl.id] ? "bg-primary-500" : "bg-gray-300"}`}
                >
                  <span
                    className={`block w-5 h-5 rounded-full bg-white transition-transform ${
                      values[ctrl.id] ? "translate-x-6" : "translate-x-0.5"
                    }`}
                  />
                </button>
              )}
              {ctrl.type === "select" && (
                <select
                  value={values[ctrl.id] as string}
                  onChange={(e) => setValues({ ...values, [ctrl.id]: e.target.value })}
                  className="border rounded px-3 py-2 text-sm w-full"
                >
                  {ctrl.options?.map((opt) => (
                    <option key={opt.value} value={opt.value}>
                      {opt.label}
                    </option>
                  ))}
                </select>
              )}
            </div>
          ))}
        </div>

        <div className="space-y-4">
          <h3 className="font-medium text-gray-700">Results</h3>
          {content.outputs?.map((out) => {
            const val = out.formula ? evaluateFormula(out.formula) : "—";
            return (
              <div key={out.id} className="bg-gray-50 rounded-lg p-4">
                <p className="text-sm text-gray-500">{out.label}</p>
                <p className="text-2xl font-bold text-gray-900">
                  {typeof val === "number" ? val.toFixed(1) : val}
                  {out.unit && <span className="text-sm font-normal text-gray-500 ml-1">{out.unit}</span>}
                </p>
              </div>
            );
          })}
        </div>
      </div>

      {content.feedback?.map((fb, i) => {
        try {
          const fn = new Function(...Object.keys(values), `return ${fb.condition}`);
          if (fn(...Object.values(values))) {
            return (
              <div key={i} className={`border rounded-lg p-4 ${feedbackColors[fb.type]}`}>
                {fb.message}
              </div>
            );
          }
        } catch {}
        return null;
      })}
    </div>
  );
}
