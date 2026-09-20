import { useState } from 'react';
import { subscribe } from '../../services/api';

export default function SubscribeForm({ position, locationLabel, onClose, isLoggedIn, onRequireLogin }) {
  const [status, setStatus] = useState('idle');
  const [errorMessage, setErrorMessage] = useState('');

  const handleSubscribe = async () => {
    if (!isLoggedIn) {
      onRequireLogin();
      return;
    }

    if (!position || position.length < 2) {
      setErrorMessage('Pilih lokasi terlebih dahulu sebelum mengaktifkan notifikasi.');
      setStatus('error');
      return;
    }

    setStatus('submitting');
    setErrorMessage('');
    try {
      await subscribe(position[0], position[1], locationLabel);
      setStatus('success');
    } catch (err) {
      setErrorMessage(err.response?.data?.error || 'Gagal mendaftar. Coba lagi.');
      setStatus('error');
    }
  };

  if (status === 'success') {
    return (
      <div className="subscribe-form subscribe-form--success">
        <p>Berhasil. Cek inbox email kamu untuk konfirmasi.</p>
        <button className="subscribe-close" onClick={onClose}>Tutup</button>
      </div>
    );
  }

  return (
    <div className="subscribe-form">
      <p className="subscribe-desc">
        Dapatkan email otomatis jika status wilayah <strong>{locationLabel}</strong> berubah jadi Waspada atau Tidak Aman.
      </p>
      <div className="subscribe-actions">
        <button className="subscribe-cancel" onClick={onClose}>Batal</button>
        <button className="subscribe-submit" onClick={handleSubscribe} disabled={status === 'submitting'}>
          {status === 'submitting' ? 'Mendaftar...' : 'Daftar Notifikasi'}
        </button>
      </div>
      {status === 'error' && <p className="subscribe-error">{errorMessage}</p>}
    </div>
  );
}