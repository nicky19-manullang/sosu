import { useState, useEffect } from 'react';
import { getStatus } from '../../services/api';
import { compareCities } from '../../data/cities';
import { statusColor } from '../../utils/statusColor';

const statusRank = { Aman: 0, Waspada: 1, 'Tidak Aman': 2 };

export default function ComparePage({ onClose, onGoToCity }) {
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    Promise.all(
      compareCities.map((city) =>
        getStatus(city.lat, city.lng)
          .then((data) => ({ ...city, ...data, failed: false }))
          .catch(() => ({ ...city, failed: true }))
      )
    ).then((data) => {
      const sorted = data
        .filter((d) => !d.failed)
        .sort((a, b) => statusRank[b.final_status] - statusRank[a.final_status]);
      setResults(sorted);
      setLoading(false);
    });
  }, []);

  const summary = results.reduce((acc, city) => {
    if (city.final_status) acc[city.final_status] = (acc[city.final_status] || 0) + 1;
    return acc;
  }, {});

  return (
    <main className="app-page compare-page">
      <div className="page-shell">
        <button className="page-back" onClick={onClose}>Kembali ke dashboard</button>
        <header className="page-hero compare-hero">
          <div>
            <span className="page-eyebrow">Pantauan nasional</span>
            <h1>Bandingkan kondisi kota</h1>
            <p>Bandingkan kualitas udara dan titik panas di kota-kota prioritas dalam satu pandangan yang mudah dibaca.</p>
          </div>
          <div className="hero-stat"><strong>{results.length || compareCities.length}</strong><span>kota dipantau</span></div>
        </header>

        <div className="summary-strip">
          <div className="summary-card summary-card--danger"><span className="summary-dot summary-dot--danger" /><strong>{summary['Tidak Aman'] || 0}</strong><span>Tidak aman</span></div>
          <div className="summary-card summary-card--warning"><span className="summary-dot summary-dot--warning" /><strong>{summary.Waspada || 0}</strong><span>Waspada</span></div>
          <div className="summary-card summary-card--safe"><span className="summary-dot summary-dot--safe" /><strong>{summary.Aman || 0}</strong><span>Aman</span></div>
        </div>

        <section className="page-section compare-section">
          <div className="compare-header-row">
            <div className="section-heading"><span className="page-eyebrow">Urutan prioritas</span><h2>Status kota hari ini</h2></div>
            <div className="compare-pill">{results.length} kota terpantau</div>
          </div>

          {loading && <p className="mysub-loading">Memuat data semua kota...</p>}

          {!loading && (
            <div className="compare-table">
              {results.map((city) => {
                const color = statusColor[city.final_status] || '#94a3b8';
                const aqiText = city.aqi_available ? `${city.aqi_category} (${city.aqi_value})` : 'Data AQI belum tersedia';
                return (
                  <article key={city.name} className="compare-card" style={{ borderTopColor: color }}>
                    <div className="compare-card-top">
                      <div>
                        <span className="compare-card-label">Kota</span>
                        <h3 className="compare-city">{city.name}</h3>
                      </div>
                      <span className="compare-badge" style={{ backgroundColor: color }}>
                        {city.final_status}
                      </span>
                    </div>

                    <div className="compare-card-metrics">
                      <div className="compare-card-metric">
                        <span className="compare-card-label">Kualitas udara</span>
                        <strong>{aqiText}</strong>
                      </div>
                      <div className="compare-card-metric">
                        <span className="compare-card-label">Titik panas</span>
                        <strong>{city.hotspot_count} titik</strong>
                      </div>
                    </div>

                    <button
                      className="compare-view-btn"
                      onClick={() => {
                        onGoToCity(city.lat, city.lng, city.name);
                        onClose();
                      }}
                    >
                      Lihat detail
                    </button>
                  </article>
                );
              })}
            </div>
          )}
        </section>
      </div>
    </main>
  );
}