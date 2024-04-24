# Cara Menjalankan

Ikuti langkah-langkah berikut untuk menjalankan aplikasi dengan Docker dan Docker Compose.

## Yang Perlu Dipastikan
- Pastikan Anda telah menginstal Docker dan Docker Compose pada sistem Anda.

## Langkah-Langkah

### 1. Menginstal Docker dan Docker Compose
Pastikan Docker dan Docker Compose telah diinstal pada sistem Anda.

### 2. Membuat Docker Network
Setelah instalasi Docker dan Docker Compose selesai, buat jaringan Docker dengan perintah berikut:
```bash
docker network create sail
```

### 3. Menjalankan Front-End

Untuk menjalankan bagian front-end, lakukan hal berikut:
- Masuk ke folder front_end.
- Jalankan Docker Compose untuk front-end dengan perintah berikut:
```bash
./vendor/bin/sail up -d
```

### 4. Menjalankan Back-End

Untuk menjalankan bagian back-end, ikuti langkah-langkah berikut:
- Masuk ke folder back_end.
- Jalankan Docker Compose dengan perintah berikut:
```bash
docker compose up -d
```

### 5. Verifikasi
Pastikan semua kontainer Docker berjalan dengan benar. Anda dapat memeriksa status kontainer dengan perintah:
```bash
docker ps
```

Jika semuanya berjalan dengan baik, applikasi seharusnya sudah bisa di akses di
- http://localhost

#### Catatan Tambahan

- Untuk menghentikan semua kontainer, gunakan `docker compose down`.
- Untuk menjalankannya kembali, gunakan `docker compose up -d`.
- Untuk informasi lebih lanjut tentang Docker dan Docker Compose, kunjungi situs web resmi mereka.