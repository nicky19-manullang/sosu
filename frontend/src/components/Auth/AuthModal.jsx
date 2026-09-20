import { useState } from 'react';

export default function AuthModal({ onLogin, onRegister, loading, error, onClose }) {
  const [mode, setMode] = useState('login');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    const success = mode === 'login' ? await onLogin(email, password) : await onRegister(email, password);
    if (success) onClose();
  };

  return (
    <div className="auth-overlay">
      <div className="auth-panel">
        <div className="auth-branding">
          <img className="auth-logo" src="/logo-karhutla.svg" alt="Logo Cek Karhutla" />
          <div>
            <p className="auth-eyebrow">CEK KARHUTLA</p>
            <h1>Pantau udara,<br />jaga sekitar.</h1>
            <p className="auth-brand-copy">
              Dapatkan informasi kualitas udara dan risiko kebakaran hutan di lokasimu.
            </p>
          </div>
          <p className="auth-brand-footer">Informasi berbasis data untuk Indonesia.</p>
        </div>

        <div className="auth-form-area">
          <div className="auth-header">
            <div>
              <p className="auth-form-kicker">Selamat datang kembali</p>
              <h2>{mode === 'login' ? 'Masuk ke akunmu' : 'Buat akun baru'}</h2>
            </div>
            <button className="close-btn" onClick={onClose} aria-label="Tutup">✕</button>
          </div>

          <form onSubmit={handleSubmit}>
            <label htmlFor="auth-email">Email</label>
            <input
              id="auth-email"
              type="email"
              placeholder="nama@email.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <label htmlFor="auth-password">Password</label>
            <input
              id="auth-password"
              type="password"
              placeholder="Minimal 6 karakter"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              minLength={6}
              required
            />

            {error && <p className="auth-error">{error}</p>}

            <button type="submit" className="auth-submit" disabled={loading}>
              {loading ? 'Memproses...' : mode === 'login' ? 'Masuk' : 'Daftar'}
            </button>
          </form>

          <p className="auth-switch">
            {mode === 'login' ? 'Belum punya akun?' : 'Sudah punya akun?'}{' '}
            <button onClick={() => setMode(mode === 'login' ? 'register' : 'login')}>
              {mode === 'login' ? 'Daftar di sini' : 'Masuk di sini'}
            </button>
          </p>
        </div>
      </div>
    </div>
  );
}