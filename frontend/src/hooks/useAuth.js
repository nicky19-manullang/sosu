import { useState, useCallback } from 'react';
import { login as apiLogin, register as apiRegister } from '../services/api';

export function useAuth() {
  const [email, setEmail] = useState(localStorage.getItem('email') || null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const doLogin = useCallback(async (emailInput, password) => {
    setLoading(true);
    setError('');
    try {
      const data = await apiLogin(emailInput, password);
      localStorage.setItem('token', data.token);
      localStorage.setItem('email', data.email);
      setEmail(data.email);
      return true;
    } catch (err) {
      setError(err.response?.data?.error || 'Login gagal. Periksa email dan password.');
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const doRegister = useCallback(async (emailInput, password) => {
    setLoading(true);
    setError('');
    try {
      const data = await apiRegister(emailInput, password);
      localStorage.setItem('token', data.token);
      localStorage.setItem('email', data.email);
      setEmail(data.email);
      return true;
    } catch (err) {
      setError(err.response?.data?.error || 'Registrasi gagal.');
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('email');
    setEmail(null);
  }, []);

  return { email, loading, error, doLogin, doRegister, logout, isLoggedIn: !!email };
}