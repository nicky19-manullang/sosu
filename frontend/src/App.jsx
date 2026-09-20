import { useState, useEffect, useCallback, useRef } from 'react';
import MapView from './components/Map/MapView';
import StatusCard from './components/StatusCard/StatusCard';
import SearchBar from './components/SearchBar/SearchBar';
import Legend from './components/Legend/Legend';
import { useGeolocation } from './hooks/useGeolocation';
import { getStatus, getHotspots, getHotspotsInBBox } from './services/api';
import './App.css';
import EducationPage from './components/Education/EducationPage';
import { useAuth } from './hooks/useAuth';
import AuthModal from './components/Auth/AuthModal';
import MySubscriptionsPage from './components/MySubscriptions/MySubscriptionsPage';
import TrendChart from './components/TrendChart/TrendChart';
import ComparePage from './components/Compare/ComparePage';

const SUMATRA_CENTER = [0.5897, 101.3431];

export default function App() {
  const [status, setStatus] = useState(null);
  const [showEducation, setShowEducation] = useState(false);
  const [loading, setLoading] = useState(false);
  const [confidenceFilter, setConfidenceFilter] = useState('all');
  const [hotspots, setHotspots] = useState([]);
  const [selectedPosition, setSelectedPosition] = useState(null);
  const boundsRequestTimer = useRef(null);
  const { status: geoStatus, error: geoError, detect } = useGeolocation();
  const hasAutoDetected = useRef(false);
  const [locationLabel, setLocationLabel] = useState('lokasi pilihan');
  const { email, loading: authLoading, error: authError, doLogin, doRegister, logout, isLoggedIn } = useAuth();
  const [showAuth, setShowAuth] = useState(false);
  const [showMySubs, setShowMySubs] = useState(false);
  const [showCompare, setShowCompare] = useState(false);
  const openPage = (page) => {
    setShowEducation(page === 'education');
    setShowCompare(page === 'compare');
    setShowMySubs(page === 'subscriptions');
  };
  const handleBoundsChange = useCallback(async (minLat, minLng, maxLat, maxLng) => {
    window.clearTimeout(boundsRequestTimer.current);
    boundsRequestTimer.current = window.setTimeout(async () => {
      try {
        const data = await getHotspotsInBBox(minLat, minLng, maxLat, maxLng);
        setHotspots(data.hotspots || []);
      } catch (err) {
        console.error(err);
      }
    }, 300);
  }, []);
  const handleLocationSelect = useCallback(async (lat, lng, label) => {
    window.clearTimeout(boundsRequestTimer.current);
    const shareUrl = new URL(window.location.href);
    shareUrl.searchParams.set('lat', lat.toFixed(6));
    shareUrl.searchParams.set('lng', lng.toFixed(6));
    window.history.replaceState({}, '', shareUrl);
    setLoading(true);
    setSelectedPosition([lat, lng]);
    if (label) setLocationLabel(label);
    try {
      const [statusData, hotspotData] = await Promise.all([
        getStatus(lat, lng),
        getHotspots(lat, lng, 50),
      ]);
      setStatus(statusData);
      setHotspots(hotspotData.hotspots || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    getHotspots(SUMATRA_CENTER[0], SUMATRA_CENTER[1], 75).then((data) =>
      setHotspots(data.hotspots || [])
    );
  }, []);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const lat = Number(params.get('lat'));
    const lng = Number(params.get('lng'));
    if (!Number.isFinite(lat) || !Number.isFinite(lng)) return;

    hasAutoDetected.current = true;
    handleLocationSelect(lat, lng);
  }, [handleLocationSelect]);

  useEffect(() => {
    if (hasAutoDetected.current) return;
    hasAutoDetected.current = true;
    detect(handleLocationSelect);
  }, [detect, handleLocationSelect]);

    return (
        <div className="app-layout">
        <div className="topbar">
        <div className="topbar-brand">
          <img className="brand-logo" src="/logo-karhutla.svg" alt="" />
          <span>Cek Karhutla</span>
        </div>
        <div className="topbar-actions">
        <button className="topbar-btn" onClick={() => openPage('education')}>
          Edukasi
        </button>
        <button className="topbar-btn" onClick={() => openPage('compare')}>
          Bandingkan Kota
        </button>
        {isLoggedIn && (
          <button className="topbar-btn" onClick={() => openPage('subscriptions')}>
          Lokasi Saya
          </button>
        )}
        {isLoggedIn ? (
          <button className="topbar-btn" onClick={logout}>
          {email} (Keluar)
          </button>
        ) : (
          <button className="topbar-btn" onClick={() => setShowAuth(true)}>
          Masuk / Daftar
          </button>
        )}
        <div className="topbar-stats">
          <div className="topbar-stat">
            <strong>{hotspots.length}</strong>
           Titik panas terlihat
          </div>
        </div>
      </div>

      </div>

      {showAuth && (
        <AuthModal
          onLogin={doLogin}
          onRegister={doRegister}
          loading={authLoading}
          error={authError}
          onClose={() => setShowAuth(false)}
        />
      )}
        {showEducation ? (
          <EducationPage onClose={() => setShowEducation(false)} />
        ) : showCompare ? (
          <ComparePage
            onClose={() => setShowCompare(false)}
            onGoToCity={handleLocationSelect}
          />
        ) : showMySubs ? (
          <MySubscriptionsPage
            onClose={() => setShowMySubs(false)}
            onGoToLocation={handleLocationSelect}
          />
        ) : (
        <div className="body-layout">
          <div className="map-container">
           <MapView
            hotspots={hotspots}
            onLocationSelect={handleLocationSelect}
            onBoundsChange={handleBoundsChange}
            selectedPosition={selectedPosition}
            confidenceFilter={confidenceFilter}
          />
          <Legend confidenceFilter={confidenceFilter} onFilterChange={setConfidenceFilter} />
          </div>
          <div className="sidebar">
            <div className="sidebar-intro">
              <p>Cek kondisi udara dan risiko karhutla di sekitarmu</p>
            </div>

            <SearchBar onLocationSelect={handleLocationSelect} />

            <button
              className="geolocate-btn"
              onClick={() => detect(handleLocationSelect)}
              disabled={geoStatus === 'detecting'}
            >
              {geoStatus === 'detecting' ? 'Mendeteksi lokasi...' : 'Gunakan Lokasi Saya'}
            </button>

            {geoError && <p className="geolocate-error">{geoError}</p>}
          <TrendChart position={selectedPosition} />
           <StatusCard
            status={status}
            loading={loading}
            position={selectedPosition}
            locationLabel={locationLabel}
            isLoggedIn={isLoggedIn}
            onRequireLogin={() => setShowAuth(true)}
          />
          </div>
        </div>
        )}
      </div>
    );
}