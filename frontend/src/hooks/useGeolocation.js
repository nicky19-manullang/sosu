import { useState, useCallback } from 'react';

export function useGeolocation() {
  const [status, setStatus] = useState('idle');
  const [error, setError] = useState(null);

  const detect = useCallback((onSuccess) => {
    if (!navigator.geolocation) {
      setStatus('unsupported');
      setError('Browser tidak mendukung deteksi lokasi.');
      return;
    }

    setStatus('detecting');
    setError(null);

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setStatus('success');
        onSuccess(position.coords.latitude, position.coords.longitude);
      },
      (err) => {
        setStatus('error');
        if (err.code === err.PERMISSION_DENIED) {
          setError('Izin lokasi ditolak. Silakan pilih lokasi manual di peta.');
        } else {
          setError('Gagal mendeteksi lokasi. Silakan pilih lokasi manual di peta.');
        }
      },
      { enableHighAccuracy: true, timeout: 10000 }
    );
  }, []);

  return { status, error, detect };
}