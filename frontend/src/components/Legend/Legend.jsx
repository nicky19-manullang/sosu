export default function Legend({ confidenceFilter, onFilterChange }) {
  const options = [
    { value: 'all', label: 'Semua', color: '#64748b' },
    { value: 'high', label: 'Confidence Tinggi', color: '#ef4444' },
    { value: 'nominal', label: 'Confidence Sedang', color: '#eab308' },
    { value: 'low', label: 'Confidence Rendah', color: '#94a3b8' },
  ];

  return (
    <div className="legend">
      <span className="legend-title">Filter Titik Panas</span>
      {options.map((opt) => (
        <button
          key={opt.value}
          className={`legend-item ${confidenceFilter === opt.value ? 'legend-item--active' : ''}`}
          onClick={() => onFilterChange(opt.value)}
        >
          <span className="legend-dot" style={{ backgroundColor: opt.color }} />
          {opt.label}
        </button>
      ))}
    </div>
  );
}