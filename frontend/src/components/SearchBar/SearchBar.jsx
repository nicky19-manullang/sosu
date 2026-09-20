import { useState } from 'react';
import { searchLocation } from '../../services/api';

export default function SearchBar({ onLocationSelect }) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!query.trim()) return;
    const data = await searchLocation(query);
    setResults(data);
  };

  const handlePick = (result) => {
  onLocationSelect(parseFloat(result.lat), parseFloat(result.lon), result.display_name);
  setResults([]);
  setQuery(result.display_name);
  };

  return (
    <div className="search-bar">
      <form onSubmit={handleSearch}>
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Cari nama kota atau kecamatan..."
        />
        <button type="submit">Cari</button>
      </form>
      {results.length > 0 && (
        <ul className="search-results">
          {results.map((r) => (
            <li key={r.place_id} onClick={() => handlePick(r)}>
              {r.display_name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}