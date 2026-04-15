interface ComparisonSlideProps {
  content: {
    heading?: string;
    description?: string;
    items?: {
      name: string;
      description?: string;
      image?: string;
      pros?: string[];
      cons?: string[];
      specs?: { label: string; value: string }[];
    }[];
  };
}

export default function ComparisonSlide({ content }: ComparisonSlideProps) {
  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">{content.heading}</h2>
      {content.description && <p className="text-gray-600">{content.description}</p>}

      <div className={`grid gap-4 ${content.items && content.items.length === 2 ? "md:grid-cols-2" : "md:grid-cols-3"}`}>
        {content.items?.map((item, i) => (
          <div key={i} className="border rounded-lg overflow-hidden">
            {item.image && (
              <div className="h-32 bg-gray-100">
                <img src={item.image} alt={item.name} className="w-full h-full object-cover" />
              </div>
            )}
            <div className="p-4 space-y-3">
              <h3 className="font-bold text-lg">{item.name}</h3>
              {item.description && <p className="text-sm text-gray-600">{item.description}</p>}

              {item.specs && item.specs.length > 0 && (
                <table className="w-full text-xs">
                  <tbody>
                    {item.specs.map((s, j) => (
                      <tr key={j} className={j % 2 === 0 ? "bg-gray-50" : ""}>
                        <td className="px-2 py-1 font-medium">{s.label}</td>
                        <td className="px-2 py-1">{s.value}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}

              {item.pros && item.pros.length > 0 && (
                <div>
                  <p className="text-xs font-medium text-green-700 mb-1">Advantages</p>
                  <ul className="space-y-1">
                    {item.pros.map((p, j) => (
                      <li key={j} className="text-xs flex items-start gap-1">
                        <span className="text-green-500">+</span> {p}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
              {item.cons && item.cons.length > 0 && (
                <div>
                  <p className="text-xs font-medium text-red-700 mb-1">Limitations</p>
                  <ul className="space-y-1">
                    {item.cons.map((c, j) => (
                      <li key={j} className="text-xs flex items-start gap-1">
                        <span className="text-red-500">-</span> {c}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
