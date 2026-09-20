# SOSU

## Monitoring Karhutla dan Kualitas Udara Indonesia

SOSU adalah aplikasi web untuk membantu masyarakat memantau kualitas udara dan risiko kebakaran hutan dan lahan di sekitar lokasi mereka. Informasi yang kompleks dari beberapa sumber data disajikan dalam bentuk peta interaktif, status yang mudah dipahami, rekomendasi tindakan, serta fitur edukasi.

## Pengembang

- Nicky Rotin Suluh Manullang
- Dea Maranatha Siregar

## Fitur Utama

- Peta interaktif Indonesia dengan data hotspot kebakaran.
- Informasi kualitas udara dan konversi PM2.5 ke kategori ISPU.
- Status wilayah: Aman, Waspada, atau Tidak Aman.
- Pencarian kota atau kecamatan melalui geocoding.
- Deteksi lokasi pengguna melalui izin browser.
- Filter dan pengelompokan titik hotspot pada peta.
- Grafik tren kondisi wilayah selama beberapa hari.
- Perbandingan status beberapa kota.
- Halaman edukasi tentang kualitas udara, hotspot, dan tindakan kesehatan.
- Informasi kontak layanan darurat resmi.
- Registrasi dan login pengguna.
- Subscription lokasi dan notifikasi email perubahan status.
- Halaman Lokasi Saya untuk mengelola subscription.

## Sumber Data

- **NASA FIRMS:** data titik panas atau hotspot dari satelit.
- **OpenAQ:** data kualitas udara dari stasiun pemantauan.
- **Nominatim OpenStreetMap:** pencarian lokasi dan geocoding.

Hotspot merupakan indikasi titik panas dari citra satelit dan bukan kepastian bahwa kebakaran sedang berlangsung. Data kualitas udara juga bergantung pada ketersediaan stasiun pemantauan di sekitar lokasi.

## Teknologi

- React 19
- Vite
- Leaflet dan React Leaflet
- Recharts
- Axios
- Go dan Gin untuk backend
- PostgreSQL dan PostGIS untuk database geospasial

## Struktur Deployment

Frontend telah disiapkan sebagai aplikasi static dan dapat dilayani melalui cPanel. Backend Go dan database PostgreSQL/PostGIS direncanakan berjalan di VPS.

```text
Domain frontend di cPanel
	|
	| HTTPS API request
	v
Subdomain API di VPS
	|
	v
Backend Go -> PostgreSQL/PostGIS
```

Database tidak perlu dibuka ke publik. Frontend hanya berkomunikasi dengan backend melalui URL API HTTPS.

## Menjalankan Frontend secara Lokal

Pastikan Node.js dan npm sudah terpasang, lalu jalankan:

```bash
npm install
npm run dev
```

Server development secara default berjalan di `http://localhost:5173`.

## Konfigurasi API

Buat file `.env` di folder `frontend/`:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

Untuk production, gunakan URL backend yang memakai HTTPS:

```env
VITE_API_BASE_URL=https://api.domain.tld/api/v1
```

Nilai environment Vite dibaca pada saat proses build, sehingga frontend perlu dibuild ulang setelah URL API production diubah.

## Build untuk cPanel

```bash
npm run build
```

Upload seluruh isi folder `dist/` ke document root domain pada cPanel. Domain frontend harus menggunakan HTTPS dan backend harus mengizinkan domain tersebut melalui konfigurasi CORS.

## Validasi Build

```bash
npm run build
npm run lint
```

Build production yang berhasil menghasilkan folder `dist/` yang siap diunggah ke hosting static.

## Catatan Produk

SOSU bukan pengganti layanan darurat resmi atau alat prediksi kebakaran. Aplikasi ini berfungsi sebagai alat pemantauan dan edukasi berbasis data yang tersedia dari sumber eksternal.

## React Compiler

The React Compiler is not enabled on this template because of its impact on dev & build performances. To add it, see [this documentation](https://react.dev/learn/react-compiler/installation).

## Expanding the Oxlint configuration

If you are developing a production application, we recommend using TypeScript with type-aware lint rules enabled. Check out the [TS template](https://github.com/vitejs/vite/tree/main/packages/create-vite/template-react-ts) for information on how to integrate TypeScript and Oxlint's TypeScript related rules in your project.
