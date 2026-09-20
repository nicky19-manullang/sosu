import { useState } from 'react';
import { statusActions, statusColor } from '../../utils/statusColor';
import SubscribeForm from '../Subscribe/SubscribeForm';

const statusIcon = {
  Aman: '✓',
  Waspada: '!',
  'Tidak Aman': '✕',
};

const trendLabel = {
  meningkat: 'Meningkat dibanding kemarin',
  menurun: 'Menurun dibanding kemarin',
  stabil: 'Stabil dibanding kemarin',
};

export default function StatusCard({ status, loading, position, locationLabel, isLoggedIn, onRequireLogin }) {
  const [showSubscribe, setShowSubscribe] = useState(false);
  const [renderedAt] = useState(() => Date.now());

  if (loading) {
    return (
      <div className="status-card status-card--empty">
        <div className="spinner" />
        <p>Memuat data lokasi...</p>
      </div>
    );
  }

  if (!status) {
    return (
      <div className="status-card status-card--empty">
        <p>Klik lokasi di peta atau cari nama kota untuk melihat status kualitas udara dan risiko karhutla.</p>
      </div>
    );
  }

  const color = statusColor[status.final_status];
  const actions = statusActions[status.final_status] || [];
  const updatedAt = status.data_updated_at ? new Date(status.data_updated_at) : null;
  const isDataStale = updatedAt && renderedAt - updatedAt.getTime() > 6 * 60 * 60 * 1000;
  const updatedLabel = updatedAt && !Number.isNaN(updatedAt.getTime())
    ? updatedAt.toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })
    : 'waktu pembaruan tidak tersedia';

  const handleShare = async () => {
    const shareUrl = window.location.href;
    const text = `Status wilayah: ${status.final_status}\nKualitas udara: ${
      status.aqi_available ? `${status.aqi_category} (ISPU ${status.aqi_value})` : 'tidak ada data'
    }\nTitik panas terdeteksi: ${status.hotspot_count}\n\nCek daerahmu juga di Cek Karhutla: ${shareUrl}`;

    if (navigator.share) {
      try {
        await navigator.share({ text, title: 'Status Karhutla', url: shareUrl });
      } catch {
        // dibatalkan user
      }
    } else {
      window.open(`https://wa.me/?text=${encodeURIComponent(text)}`, '_blank');
    }
  };

  return (
    <div className="status-card" style={{ '--status-color': color }}>
      <div className="status-header">
        <div className="status-icon" style={{ backgroundColor: color }}>
          {statusIcon[status.final_status]}
        </div>
        <div>
          <div className="status-label">Status Wilayah</div>
          <div className="status-value" style={{ color }}>{status.final_status}</div>
        </div>
      </div>

      <p className={`status-freshness${isDataStale ? ' status-freshness--stale' : ''}`}>
        Data sumber terbaru: {updatedLabel}{isDataStale ? ' (perlu diperbarui)' : ''}
      </p>

      <div className="status-details">
        <div className="status-row">
          <span className="status-row-label">Kualitas udara</span>
          <span className="status-row-value">
            {status.aqi_available
              ? `${status.aqi_category} (ISPU ${status.aqi_value})`
              : 'Tidak ada data terdekat'}
          </span>
        </div>
        {status.aqi_available && status.aqi_distance_km > 5 && (
          <p className="status-note">
            *Data dari stasiun {status.aqi_distance_km.toFixed(1)} km dari lokasi ini
          </p>
        )}
        <div className="status-row">
          <span className="status-row-label">Titik panas (radius 10 km)</span>
          <span className="status-row-value">{status.hotspot_count} titik</span>
        </div>
        {status.hotspot_trend && (
          <p className="status-note">{trendLabel[status.hotspot_trend]}</p>
        )}
      </div>

      <div className="status-recommendation">
        <strong>Yang bisa dilakukan</strong>
        <ul>
          {actions.map((action) => <li key={action}>{action}</li>)}
        </ul>
      </div>

      <div className="status-buttons">
        <button className="share-btn" onClick={handleShare}>
          Bagikan
        </button>
        <button className="subscribe-btn" onClick={() => setShowSubscribe(true)}>
          Aktifkan notifikasi
        </button>
      </div>

      {showSubscribe && (
      <SubscribeForm
        position={position}
        locationLabel={locationLabel}
        onClose={() => setShowSubscribe(false)}
        isLoggedIn={isLoggedIn}
        onRequireLogin={onRequireLogin}
      />
    )}
      <p className="status-disclaimer">Catatan: {status.disclaimer}</p>
    </div>
  );
}