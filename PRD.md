# PRD: Monitoring Karhutla & Kualitas Udara Indonesia

**Status:** Final demo release — fitur inti selesai, frontend sudah di cPanel, backend dan database ditargetkan di VPS
**Owner:** Manullang
**Nama project:** sosu

---

## 1. Overview

### 1.1 Problem Statement
Masyarakat umum di wilayah rawan karhutla (Sumatera, Kalimantan, Sulawesi) tidak punya cara cepat dan mudah dipahami untuk mengecek kondisi udara dan risiko kebakaran di sekitar mereka. Platform pemerintah yang ada (SiPongi, ISPU KLHK) bersifat teknis dan tidak dirancang untuk awam.

### 1.2 Value Proposition
"Cek dalam 10 detik — daerah saya aman atau tidak, dan apa yang harus saya lakukan."

### 1.3 Target Pengguna
Masyarakat umum awam di wilayah rawan karhutla. Bukan peneliti/teknis. Tidak familiar istilah AQI/PM2.5/confidence level. Mayoritas akses dari HP.

---

## 2. Goals & Non-Goals

### 2.1 Goals
- Menampilkan status kualitas udara + hotspot kebakaran dalam satu tampilan sederhana
- Memberi rekomendasi aksi konkret, bukan cuma angka mentah
- Mobile-first, ringan, minim klik
- Menggunakan sumber data gratis atau free tier selama masih mencukupi kebutuhan demo

### 2.2 Non-Goals (di luar scope untuk saat ini)
- Bukan pengganti layanan darurat resmi (BPBD, Damkar, layanan kesehatan)
- Bukan alat prediksi/forecasting kebakaran (itu domain KarhutlaFlow, project terpisah)
- Tidak menyediakan aplikasi mobile native di MVP (web-based, mobile-responsive)
- Tidak menjamin data real-time detik-per-detik — mengikuti frekuensi update sumber data asli (FIRMS tiap 30 menit, OpenAQ tiap 1 jam via scheduler, plus fetch awal otomatis saat server startup)

> **Perubahan dari v1:** cakupan geografis awalnya dibatasi Sumatera saja, tapi diperluas jadi **nasional** (bounding box `95,-11,141,6`) karena keterbatasan data kalau dibatasi wilayah kecil.
>
> **Perubahan dari v2:** akun user/login ditambahkan karena kebutuhan sistem notifikasi yang membatasi "1 email = 1 akun". Login memakai Email+Password dengan sesi JWT.

---

## 3. Fitur & Scope per Fase

### Fase 1 — MVP
- Peta interaktif dengan overlay hotspot + zona kualitas udara
- Pencarian lokasi via nama kota/kecamatan (geocoding)
- Auto-detect lokasi user (dengan izin browser)
- Kartu status wilayah: status AQI (kata-kata, bukan cuma angka), jumlah & jarak hotspot terdekat, rekomendasi aksi konkret
- Cakupan geografis: **nasional Indonesia** (bounding box `95,-11,141,6`)
- 1 sumber data hotspot (NASA FIRMS) + 1 sumber data AQI (OpenAQ, difilter `iso=ID`); jika tidak ada stasiun AQI dalam radius 75km dari lokasi yang dicek, kualitas udara ditampilkan sebagai "tidak ada data" alih-alih memaksakan data dari stasiun yang terlalu jauh

### Fase 2 — ✅ Selesai
- Halaman edukasi (FAQ, panduan kesehatan per level ISPU, kontak layanan darurat)
- Sistem akun (register/login Email+Password, JWT) — ditambahkan karena kebutuhan notifikasi
- Notifikasi/alert email berlangganan per wilayah (SMTP), dikirim otomatis sekali per perubahan status ke Waspada/Tidak Aman
- Tren 24 jam terakhir per wilayah (dibandingkan hari sebelumnya)
- Auto-fetch data saat server startup + cleanup data hotspot otomatis (retensi 7 hari, dijalankan harian jam 3 pagi)

### Fase 3 — ✅ Selesai untuk demo
- Grafik tren historis 7 hari (jumlah hotspot & ISPU per hari) di kartu status
- Perbandingan antar kota (tabel status beberapa kota besar sekaligus)
- API publik untuk developer lain (ditunda ke fase pasca-demo)
- Dashboard "Lokasi Saya" untuk user yang login: daftar semua lokasi yang di-subscribe beserta status terkini, dengan opsi unsubscribe
- **Weekly digest email**: ringkasan mingguan otomatis dikirim tiap Senin ke semua user yang punya subscription aktif, mencakup status semua lokasi yang di-subscribe (baik dalam kondisi Aman maupun Waspada/Tidak Aman — tidak hanya saat kondisi bahaya seperti alert biasa)

---

## 4. Tech Stack

| Layer | Pilihan | Alasan |
|---|---|---|
| Frontend | React 19 + Vite + Leaflet + Recharts | Ringan, gratis, dan tidak membutuhkan API key peta berbayar |
| Backend | Go (Gin) | Ringan untuk tugas fetch-cache-serve data eksternal, cocok untuk scheduled job |
| Database | PostgreSQL + PostGIS | Query geospasial radius dan bbox native; development memakai PostgreSQL lokal, production ditargetkan di VPS |
| Cache/Scheduler | `robfig/cron/v3` | FIRMS di-fetch tiap 30 menit, OpenAQ tiap 1 jam, hindari hit API eksternal tiap request user |
| Deployment | Frontend di cPanel; backend + PostgreSQL/PostGIS di VPS | Frontend static dilayani cPanel, API memakai subdomain HTTPS dan database hanya diakses backend |
| Geocoding | Nominatim (OpenStreetMap) | Gratis, dipakai untuk pencarian nama kota/kecamatan |

---

## 5. Sumber Data

| Data | Sumber | Biaya | Cakupan | Status |
|---|---|---|---|---|
| Hotspot | NASA FIRMS API (VIIRS_SNPP_NRT) | Gratis (API key gratis) | Nasional, bbox `95,-11,141,6` | ✅ Live, auto-fetch tiap 30 menit |
| Kualitas udara | OpenAQ v3 | Gratis (API key wajib) | Nasional, filter `iso=ID` | ✅ Live, auto-fetch tiap 1 jam |
| Kualitas udara (opsional tambahan) | IQAir API | Free tier terbatas | — | Belum dipakai, cadangan jika coverage OpenAQ kurang |
| Geocoding | Nominatim OSM | Gratis | Global | ✅ Dipakai di search bar frontend |

**Batasan penting:** ISPU KLHK dan SiPongi+ tidak dipakai di MVP karena tidak ada API publik resmi yang stabil — kemungkinan butuh scraping yang rawan berubah. Ini didokumentasikan sebagai batasan teknis, bukan diabaikan sepenuhnya (bisa dievaluasi ulang di fase lanjut).

---

## 6. Logika Status Aman/Tidak Aman

### 6.1 Kualitas Udara (berdasarkan ISPU, Permen LHK No. 14/2020)

| Rentang ISPU | Kategori | Status Biner |
|---|---|---|
| 0–50 | Baik | Aman |
| 51–100 | Sedang | Aman |
| 101–200 | Tidak Sehat | Tidak Aman |
| 201–300 | Sangat Tidak Sehat | Tidak Aman |
| >300 | Berbahaya | Tidak Aman |

### 6.2 Hotspot (threshold desain internal, belum ada baku mutu nasional)

| Kondisi | Status |
|---|---|
| 0 titik radius 10km, atau semua confidence low | Aman |
| 1–5 titik confidence nominal/high radius 10km | Waspada |
| >5 titik, atau ada titik confidence high radius 5km | Tidak Aman |

### 6.3 Status Gabungan
`status_akhir = MAX(status_aqi, status_hotspot)` — ambil yang paling parah dari dua sinyal.

**Disclaimer wajib ditampilkan:** hotspot adalah indikasi titik panas dari citra satelit, bukan kepastian kebakaran aktif.

---

## 7. Struktur Proyek

### 7.1 Frontend (React)
```
frontend/
├── src/
│   ├── components/
│   │   ├── Map/
│   │   ├── StatusCard/
│   │   ├── SearchBar/
│   │   ├── Legend/
│   │   ├── Education/
│   │   ├── Compare/
│   │   ├── MySubscriptions/
│   │   ├── Auth/
│   │   └── TrendChart/
│   ├── hooks/
│   ├── services/          # pemanggilan API backend
│   ├── utils/
│   │   └── statusColor.js
│   └── App.jsx
├── public/
└── package.json
```

### 7.2 Backend (Go)
```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/            # HTTP handlers (Gin)
│   ├── service/            # business logic (status calculation, dsb)
│   ├── repository/         # akses database
│   ├── fetcher/            # job fetch dari FIRMS, OpenAQ, Nominatim
│   ├── scheduler/          # cron job
│   └── model/
├── migrations/
├── config/
└── go.mod
```

---

## 8. API Design (implemented)

| Endpoint | Deskripsi |
|---|---|
| `GET /api/v1/status?lat=&lng=` | Status gabungan (AQI + hotspot) untuk satu titik |
| `GET /api/v1/hotspots?lat=&lng=&radius_km=` | Titik hotspot dalam radius lokasi |
| `GET /api/v1/hotspots/bbox?min_lat=&min_lng=&max_lat=&max_lng=` | Titik hotspot dalam area peta yang sedang dilihat |
| `GET /api/v1/aqi/nearest?lat=&lng=` | Data AQI terdekat dari titik |
| `GET /api/v1/trend?lat=&lng=&days=` | Tren historis lokasi |
| `POST /api/v1/auth/register` | Registrasi akun |
| `POST /api/v1/auth/login` | Login dan penerbitan JWT |
| `POST /api/v1/subscribe` | Berlangganan notifikasi lokasi; membutuhkan JWT |
| `POST /api/v1/unsubscribe` | Menghapus subscription; membutuhkan JWT |
| `GET /api/v1/my-subscriptions` | Daftar subscription user; membutuhkan JWT |
| `GET /health` | Health check API dan koneksi database |

Pencarian lokasi memakai Nominatim langsung dari frontend, bukan endpoint backend `/api/v1/search`.

---

## 9. Batasan & Constraints

- **Akurasi data:** hotspot ≠ kebakaran pasti; AQI di luar kota besar bisa kosong (ditampilkan jujur sebagai "tidak ada data", bukan estimasi/dipaksakan)
- **Bukan layanan darurat:** tidak menggantikan BPBD/Damkar/layanan kesehatan resmi
- **Cakupan geografis:** nasional Indonesia (diperluas dari rencana awal Sumatera saja)
- **Radius pencarian AQI dibatasi 75km:** jika tidak ada stasiun dalam radius itu, status akhir wilayah dihitung murni dari data hotspot
- **Frekuensi update:** FIRMS tiap 30 menit, OpenAQ tiap 1 jam via scheduler otomatis — bukan live streaming
- **Retensi data hotspot:** data lama dibersihkan otomatis setiap hari pukul 03.00 dengan retensi 7 hari
- **Legalitas:** wajib mencantumkan atribusi sumber data (NASA FIRMS, OpenAQ, BMKG bila dipakai)
- **Akun user:** hanya diperlukan untuk subscription dan notifikasi; peta, status, pencarian, edukasi, dan perbandingan tetap dapat diakses publik
- **Deployment:** frontend cPanel harus memakai URL API production melalui `VITE_API_BASE_URL`; backend harus mengizinkan domain frontend melalui konfigurasi CORS
- **Keamanan production:** database tidak boleh membuka port `5432` ke publik; gunakan HTTPS untuk frontend dan API, secret production, serta backup database berkala
- **Anggaran:** target tetap memakai free tier, tetapi biaya VPS, domain, SMTP, atau layanan eksternal dapat muncul di production

---

## 10. Roadmap Ringkas

1. **Fase 1 (MVP):** peta + status + rekomendasi, cakupan nasional — ✅ selesai
2. **Fase 2:** edukasi + akun/login + notifikasi + tren jangka pendek — ✅ selesai
3. **Fase 3:** grafik tren historis + perbandingan wilayah + dashboard lokasi saya + weekly digest — ✅ selesai untuk demo
4. **Deployment production:** frontend cPanel — ✅ selesai; backend Go + PostgreSQL/PostGIS di VPS — ⬜ langkah berikutnya
5. **Pasca-demo:** API publik terdokumentasi, rate limiting, observability, backup otomatis, dan optimasi bundle frontend

## 11. Status Implementasi Final

| Komponen | Status | Catatan |
|---|---|---|
| Struktur project (`sosu/frontend`, `sosu/backend`) | ✅ Selesai | |
| Model + Migration (PostgreSQL + PostGIS) | ✅ Selesai | Tabel `hotspots`, `aqi_readings`, `users`, `subscriptions` |
| Repository layer | ✅ Selesai | Query radius pakai `ST_DWithin`, `FindNearest` AQI dibatasi radius 75km |
| Fetcher (FIRMS + OpenAQ) | ✅ Selesai | Cakupan nasional + global, filter `iso=ID` untuk hindari data negara tetangga |
| Konversi PM2.5 → ISPU | ✅ Selesai | Sesuai breakpoint resmi Permen LHK No. 14/2020 |
| Service (logika status gabungan) | ✅ Selesai | `MAX(status_aqi, status_hotspot)`, tren 24 jam, AQI diabaikan jika tidak ada stasiun dalam radius |
| Handler / REST API | ✅ Selesai | `GET /status`, `/hotspots`, `/hotspots/bbox`, `/aqi/nearest` |
| Scheduler otomatis | ✅ Selesai | FIRMS tiap 30 menit, OpenAQ tiap 1 jam, FIRMS global tiap 3 jam, fetch awal saat startup, cleanup harian |
| Frontend peta (peta, kartu status, search, geolocation, filter, cluster) | ✅ Selesai | Leaflet + React |
| Halaman edukasi | ✅ Selesai | Modal FAQ + panduan ISPU + kontak darurat |
| Auth (register/login, JWT) | ✅ Selesai | Email+Password, 1 email = 1 akun |
| Notifikasi email (subscribe/alert) | ✅ Selesai | SMTP Gmail, alert sekali per perubahan status |
| Cleanup data hotspot lama (`DeleteOlderThan`) | ✅ Selesai | Dijadwalkan harian jam 3 pagi, retensi 7 hari |
| CORS untuk domain production | ⬜ Konfigurasi deployment | Origin production harus diatur sesuai domain frontend saat backend dijalankan di VPS |
| Grafik tren historis 7 hari | ✅ Selesai | Fase 3 |
| Perbandingan antar kota | ✅ Selesai | Fase 3 |
| API publik + dokumentasi | ⬜ Pasca-demo | Endpoint internal sudah tersedia; rate limit dan dokumentasi publik belum menjadi release demo |
| Dashboard "Lokasi Saya" | ✅ Selesai | Fase 3 |
| Weekly digest email | ✅ Selesai | Fase 3, dikirim tiap Senin ke semua subscription aktif |

## 12. Checklist Deployment Final

### 12.1 Frontend di cPanel

- Build frontend dengan `VITE_API_BASE_URL=https://api.domain.tld/api/v1`.
- Upload isi folder `frontend/dist/` ke document root cPanel.
- Pastikan konfigurasi rewrite SPA mengarahkan route yang tidak berupa file ke `index.html`.
- Pastikan domain frontend memakai HTTPS.

### 12.2 Backend dan Database di VPS

- Siapkan Linux VPS, Go runtime atau binary hasil build, dan PostgreSQL dengan ekstensi PostGIS.
- Buat database production dan jalankan migration dari folder `backend/migrations/`.
- Isi environment production: `DATABASE_URL`, `PORT`, `FIRMS_MAP_KEY`, `OPENAQ_API_KEY`, konfigurasi SMTP, dan secret autentikasi.
- Jalankan backend sebagai service yang otomatis restart.
- Pasang reverse proxy Nginx pada subdomain API, misalnya `api.domain.tld`, lalu aktifkan SSL.
- Atur CORS agar hanya mengizinkan origin frontend production.
- Buka firewall hanya untuk SSH, HTTP, dan HTTPS; jangan expose port database `5432`.
- Uji `GET /health`, login, pencarian lokasi, status, hotspot, subscription, dan pengiriman email.
- Siapkan backup database sebelum demo dan backup terjadwal setelah production.

### 12.3 Kriteria Siap Demo

- Dashboard dapat dibuka dari domain publik tanpa error browser atau mixed content.
- Status, hotspot, AQI, tren, dan pencarian lokasi menampilkan data atau pesan "tidak ada data" yang jelas.
- Halaman Edukasi, Bandingkan Kota, dan Lokasi Saya dapat dibuka dari navigasi.
- Login, subscription, unsubscribe, dan logout sudah diuji dengan akun demo.
- Health check API merespons `status: ok`.
- Ada rencana cadangan bila API FIRMS/OpenAQ atau SMTP sedang tidak tersedia.