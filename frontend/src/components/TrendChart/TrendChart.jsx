import { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { getTrend } from '../../services/api';

export default function TrendChart({ position }) {
  const [trend, setTrend] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!position) return;
    setLoading(true);
    getTrend(position[0], position[1], 7)
      .then((data) => setTrend(data.trend || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [position]);

  if (!position || loading) return null;
  if (trend.length === 0) return null;

  const chartData = trend.map((t) => ({
    date: t.date.slice(5),
    'Titik Panas': t.hotspot_count,
    ISPU: t.aqi_available ? Math.round(t.aqi_value) : null,
  }));

  return (
    <div className="trend-chart">
      <div className="trend-heading">
        <div>
          <span className="eyebrow">Riwayat wilayah</span>
          <h4>Tren 7 hari terakhir</h4>
        </div>
        <span className="trend-period">7 hari</span>
      </div>
      <ResponsiveContainer width="100%" height={180}>
        <LineChart data={chartData} margin={{ top: 5, right: 10, left: -20, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
          <XAxis dataKey="date" tick={{ fontSize: 11 }} />
          <YAxis tick={{ fontSize: 11 }} />
          <Tooltip />
          <Line type="monotone" dataKey="Titik Panas" stroke="#ef4444" strokeWidth={2} dot={{ r: 3 }} />
          <Line type="monotone" dataKey="ISPU" stroke="#eab308" strokeWidth={2} dot={{ r: 3 }} connectNulls />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}