import axios from 'axios';

const apiBaseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: apiBaseURL,
  timeout: 15000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

const readResponseData = (promise) => promise.then((res) => res.data);

export const getStatus = (lat, lng) =>
  readResponseData(api.get('/status', { params: { lat, lng } }));

export const getHotspots = (lat, lng, radiusKm = 100) =>
  readResponseData(api.get('/hotspots', { params: { lat, lng, radius_km: radiusKm } }));

export const getHotspotsInBBox = (minLat, minLng, maxLat, maxLng) =>
  readResponseData(
    api.get('/hotspots/bbox', {
      params: { min_lat: minLat, min_lng: minLng, max_lat: maxLat, max_lng: maxLng },
    })
  );

export const searchLocation = (query) =>
  axios
    .get('https://nominatim.openstreetmap.org/search', {
      params: { q: `${query}, Indonesia`, format: 'json', limit: 5 },
    })
    .then((res) => res.data);

export const register = (email, password) =>
  readResponseData(api.post('/auth/register', { email, password }));

export const login = (email, password) =>
  readResponseData(api.post('/auth/login', { email, password }));

export const subscribe = (lat, lng, locationLabel) =>
  readResponseData(api.post('/subscribe', { latitude: lat, longitude: lng, location_label: locationLabel }));

export const unsubscribe = (lat, lng) =>
  readResponseData(api.post('/unsubscribe', { latitude: lat, longitude: lng }));

export const getMySubscriptions = () => readResponseData(api.get('/my-subscriptions'));

export const getTrend = (lat, lng, days = 7) =>
  readResponseData(api.get('/trend', { params: { lat, lng, days } }));