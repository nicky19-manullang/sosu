import { MapContainer, TileLayer, CircleMarker, Marker, Popup, useMapEvents } from 'react-leaflet';
import MarkerClusterGroup from 'react-leaflet-cluster';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { statusColor } from '../../utils/statusColor';

const selectedIcon = new L.DivIcon({
  html: `<div style="
    width: 20px; height: 20px; border-radius: 50%;
    background: #2563eb; border: 3px solid white;
    box-shadow: 0 2px 8px rgba(0,0,0,0.4);
  "></div>`,
  className: '',
  iconSize: [20, 20],
  iconAnchor: [10, 10],
});

const formatDate = (value) => {
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? 'Waktu tidak valid' : date.toLocaleString('id-ID');
};

function ClickHandler({ onLocationSelect }) {
  useMapEvents({
    click(e) {
      onLocationSelect(e.latlng.lat, e.latlng.lng);
    },
  });
  return null;
}

function ViewportWatcher({ onBoundsChange }) {
  const map = useMapEvents({
    moveend() {
      const b = map.getBounds();
      onBoundsChange(b.getSouth(), b.getWest(), b.getNorth(), b.getEast());
    },
  });
  return null;
}

const isLikelyLandHotspot = (lat, lng) => {
  if (lat < -12 || lat > 8 || lng < 94 || lng > 141) return false;

  const blocks = [
    { minLat: -4.5, maxLat: 6.2, minLng: 95.0, maxLng: 106.5 },
    { minLat: -10.0, maxLat: -5.0, minLng: 104.0, maxLng: 116.0 },
    { minLat: -4.8, maxLat: 4.8, minLng: 108.0, maxLng: 119.5 },
    { minLat: -8.5, maxLat: 2.2, minLng: 118.2, maxLng: 126.5 },
    { minLat: -10.5, maxLat: 3.0, minLng: 126.0, maxLng: 141.0 },
    { minLat: -9.5, maxLat: -0.5, minLng: 130.0, maxLng: 141.0 },
  ];

  return blocks.some((block) => lat >= block.minLat && lat <= block.maxLat && lng >= block.minLng && lng <= block.maxLng);
};

export default function MapView({ hotspots, onLocationSelect, onBoundsChange, selectedPosition, confidenceFilter }) {
  const filteredHotspots = (hotspots || []).filter((h) => {
    if (!isLikelyLandHotspot(h.Latitude, h.Longitude)) return false;
    return confidenceFilter === 'all' ? true : h.Confidence === confidenceFilter;
  });

  return (
    <MapContainer center={[0.5897, 101.3431]} zoom={5} style={{ height: '100%', width: '100%' }}>
      <TileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution='&copy; OpenStreetMap contributors'
      />
      <ClickHandler onLocationSelect={onLocationSelect} />
      <ViewportWatcher onBoundsChange={onBoundsChange} />

      <MarkerClusterGroup chunkedLoading maxClusterRadius={40}>
        {filteredHotspots.map((h) => (
          <CircleMarker
            key={h.ID}
            center={[h.Latitude, h.Longitude]}
            radius={6}
            pathOptions={{
              color: h.Confidence === 'high' ? statusColor['Tidak Aman'] : statusColor.Waspada,
              fillColor: h.Confidence === 'high' ? statusColor['Tidak Aman'] : statusColor.Waspada,
              fillOpacity: 0.7,
              weight: 2,
            }}
          >
            <Popup>
              <strong>Titik Panas</strong>
              <br />
              Confidence: {h.Confidence}
              <br />
              FRP: {h.FRP}
              <br />
              Terdeteksi: {formatDate(h.AcquiredAt)}
            </Popup>
          </CircleMarker>
        ))}
      </MarkerClusterGroup>

      {selectedPosition && (
        <Marker position={selectedPosition} icon={selectedIcon}>
          <Popup>Lokasi yang dipilih</Popup>
        </Marker>
      )}
    </MapContainer>
  );
}