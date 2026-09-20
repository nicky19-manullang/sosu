import { useState, useEffect, useCallback } from 'react';
import { getMySubscriptions, unsubscribe } from '../../services/api';
import { statusColor } from '../../utils/statusColor';

export default function MySubscriptionsPage({ onClose, onGoToLocation }) {
  const [subs, setSubs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [removingId, setRemovingId] = useState(null);

  const loadSubs = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getMySubscriptions();
      setSubs(data.subscriptions || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadSubs();
  }, [loadSubs]);

  const handleUnsubscribe = async (sub) => {
    setRemovingId(sub.id);
    try {
      await unsubscribe(sub.latitude, sub.longitude);
      setSubs((prev) => prev.filter((s) => s.id !== sub.id));
    } catch (err) {
      console.error(err);
    } finally {
      setRemovingId(null);
    }
  };

  return (
    <main className="app-page subscriptions-page">
      <div className="page-shell">
        <button className="page-back" onClick={onClose}>Kembali ke dashboard</button>
        <header className="page-hero subscriptions-hero">
          <div>
            <span className="page-eyebrow">Ruang pemantauan pribadi</span>
            <h1>Lokasi saya</h1>
            <p>Semua wilayah yang kamu pantau, dirangkum dalam satu ruang agar perubahan status tidak terlewat.</p>
          </div>
          <div className="hero-stat"><strong>{subs.length}</strong><span>lokasi aktif</span></div>
        </header>

        <section className="page-section subscriptions-section">
          <div className="section-heading"><span className="page-eyebrow">Notifikasi wilayah</span><h2>Daftar pemantauan</h2></div>
          {loading && <p className="mysub-loading">Memuat lokasi...</p>}

        {!loading && subs.length === 0 && (
          <p className="mysub-empty">
            Belum ada lokasi yang dipantau. Klik lokasi di peta lalu tekan "Aktifkan notifikasi" untuk mulai berlangganan.
          </p>
        )}

        {!loading && subs.length === 0 && <div className="subscription-empty-panel"><strong>Belum ada lokasi tersimpan</strong><p>Pilih wilayah dari peta, lalu aktifkan notifikasi untuk mulai memantau perubahan statusnya.</p><button className="mysub-view-btn" onClick={onClose}>Cari lokasi di peta</button></div>}

        <div className="mysub-list">
          {subs.map((sub) => {
            const color = statusColor[sub.final_status] || '#94a3b8';
            return (
              <div className="mysub-card" key={sub.id} style={{ borderLeftColor: color }}>
                <div className="mysub-card-header">
                  <strong>{sub.location_label}</strong>
                  <span className="mysub-badge" style={{ backgroundColor: color }}>
                    {sub.final_status || 'Tidak diketahui'}
                  </span>
                </div>
                <div className="mysub-card-body">
                  <span>Udara: {sub.aqi_available ? sub.aqi_category : 'Data tidak tersedia'}</span>
                  <span>Hotspot: {sub.hotspot_count} titik</span>
                </div>
                <div className="mysub-card-actions">
                  <button
                    className="mysub-view-btn"
                    onClick={() => {
                      onGoToLocation(sub.latitude, sub.longitude, sub.location_label);
                      onClose();
                    }}
                  >
                    Lihat di Peta
                  </button>
                  <button
                    className="mysub-remove-btn"
                    onClick={() => handleUnsubscribe(sub)}
                    disabled={removingId === sub.id}
                  >
                    {removingId === sub.id ? 'Menghapus...' : 'Berhenti Pantau'}
                  </button>
                </div>
              </div>
            );
          })}
        </div>
        </section>
      </div>
    </main>
  );
}