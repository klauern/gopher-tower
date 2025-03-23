export default function Loading() {
  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="bg-white shadow-sm rounded-lg p-6 animate-pulse">
        <div className="flex justify-between items-start mb-6">
          <div>
            <div className="h-8 w-64 bg-gray-200 rounded"></div>
            <div className="mt-2 h-4 w-96 bg-gray-200 rounded"></div>
          </div>
          <div className="h-6 w-24 bg-gray-200 rounded"></div>
        </div>

        <div className="grid grid-cols-2 gap-6 mt-6">
          {[...Array(4)].map((_, i) => (
            <div key={i}>
              <div className="h-4 w-24 bg-gray-200 rounded mb-2"></div>
              <div className="h-6 w-32 bg-gray-200 rounded"></div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
