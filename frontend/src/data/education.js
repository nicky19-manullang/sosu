export const faqItems = [
  {
    q: 'Apa itu hotspot / titik panas?',
    a: 'Hotspot adalah titik di permukaan bumi yang terdeteksi memiliki suhu jauh lebih tinggi dari sekitarnya oleh satelit. Ini bisa menandakan kebakaran hutan/lahan, tapi juga bisa disebabkan hal lain seperti pembakaran sampah, aktivitas industri, atau pantulan panas matahari. Karena itu, hotspot adalah indikasi, bukan kepastian kebakaran.',
  },
  {
    q: 'Apa itu ISPU?',
    a: 'ISPU (Indeks Standar Pencemar Udara) adalah angka yang menunjukkan kondisi kualitas udara di suatu wilayah, ditetapkan oleh Kementerian Lingkungan Hidup dan Kehutanan. Semakin tinggi angkanya, semakin buruk kualitas udaranya.',
  },
  {
    q: 'Kenapa asap karhutla berbahaya?',
    a: 'Asap kebakaran hutan mengandung partikel halus (PM2.5) yang bisa masuk jauh ke saluran pernapasan hingga aliran darah. Paparan jangka pendek bisa menyebabkan iritasi mata dan tenggorokan, sesak napas, dan memperberat penyakit asma/jantung. Paparan jangka panjang berisiko menyebabkan gangguan paru-paru kronis.',
  },
  {
    q: 'Kenapa hotspot terdeteksi tapi ISPU masih Baik?',
    a: 'Hotspot menunjukkan lokasi titik api, sementara ISPU mengukur kualitas udara di sekitar kita. Asap butuh waktu untuk menyebar dan arahnya tergantung angin — jadi ada kemungkinan kebakaran terjadi tapi asapnya belum/tidak terbawa ke wilayah pemukiman.',
  },
  {
    q: 'Seberapa akurat data di website ini?',
    a: 'Data hotspot berasal dari satelit NASA (FIRMS) yang diperbarui tiap 30 menit, dan data kualitas udara dari jaringan stasiun OpenAQ yang diperbarui tiap jam. Akurasi bergantung pada cakupan satelit dan kepadatan stasiun pemantauan — di beberapa wilayah, data AQI mungkin belum tersedia karena belum ada stasiun terdekat.',
  },
];

export const ispuGuide = [
  {
    range: '0–50',
    category: 'Baik',
    color: '#22c55e',
    activity: 'Aman melakukan aktivitas luar ruangan seperti biasa.',
    mask: 'Tidak diperlukan.',
    doctor: 'Tidak ada tindakan khusus.',
  },
  {
    range: '51–100',
    category: 'Sedang',
    color: '#eab308',
    activity: 'Masih dapat beraktivitas normal. Kelompok sensitif (anak, lansia, penderita asma) mulai kurangi aktivitas fisik berat di luar ruangan.',
    mask: 'Opsional bagi kelompok sensitif.',
    doctor: 'Tidak ada tindakan khusus, tetap waspada bagi kelompok sensitif.',
  },
  {
    range: '101–200',
    category: 'Tidak Sehat',
    color: '#f97316',
    activity: 'Kurangi aktivitas luar ruangan yang lama, terutama anak-anak, lansia, ibu hamil, dan penderita gangguan pernapasan.',
    mask: 'Gunakan masker N95 saat keluar rumah.',
    doctor: 'Konsultasi dokter bila muncul gejala sesak napas atau iritasi mata/tenggorokan.',
  },
  {
    range: '201–300',
    category: 'Sangat Tidak Sehat',
    color: '#ef4444',
    activity: 'Hindari aktivitas luar ruangan. Tutup jendela dan pintu rumah.',
    mask: 'Wajib gunakan masker N95 bila terpaksa keluar rumah.',
    doctor: 'Segera ke fasilitas kesehatan bila mengalami sesak napas.',
  },
  {
    range: '>300',
    category: 'Berbahaya',
    color: '#a855f7',
    activity: 'Hentikan seluruh aktivitas luar ruangan. Pertimbangkan evakuasi bila memungkinkan.',
    mask: 'Wajib masker N95 dan minimalkan waktu di luar rumah sama sekali.',
    doctor: 'Segera cari pertolongan medis, terutama untuk kelompok rentan.',
  },
];

export const emergencyContacts = [
  { name: 'Pemadam Kebakaran (Damkar)', number: '113' },
  { name: 'Ambulans dan Layanan Medis', number: '118 / 119' },
  { name: 'Basarnas / SAR', number: '115' },
  { name: 'BNPB (Pusat Pengendalian Operasi)', number: '117' },
  { name: 'Layanan darurat wilayah tertentu', number: '1131' },
];