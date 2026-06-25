# 🛡️ ReksaFel — Sistem Pengawasan & Anti-Kecurangan Ujian Lab

**ReksaFel** adalah platform pengawas ujian lokal berkinerja tinggi yang dirancang untuk menjaga kejujuran akademik di laboratorium komputer sekolah maupun universitas. Melalui jaringan tertutup (*peer-to-peer mesh*) terenkripsi, pengawas dapat memantau aktivitas, mendeteksi potensi kecurangan, dan mengumpulkan bukti tangkapan layar dari puluhan komputer siswa secara real-time dari satu dashboard utama.

Repositori ini berfungsi sebagai pusat dokumentasi publik, spesifikasi API inti, dan tinjauan teknis untuk sistem ReksaFel.

---

## 🌌 Konsep Utama & Arsitektur

ReksaFel dioptimalkan untuk jaringan lokal tanpa bergantung pada server cloud (offline-first), mengubah laboratorium komputer standar menjadi zona ujian yang aman dan terkontrol.

```text
  +-------------------------------+                 +-----------------------------------+
  |    PC Siswa (Agent Client)    |                 |     PC Pengawas (Admin Dashboard) |
  |   - Servis telemetri latar    |                 |     - Pengendali API berbasis Go  |
  |   - Agen pemantauan layar     |                 |     - Desktop Wrapper (Wails)     |
  +---------------+---------------+                 +-----------------+-----------------+
                  |                                                   |
                  |     Unggah Log Telemetri & Bukti Screenshot       |
                  +-------------------------------------------------->| [HTTP API (Port Telemetri)]
                  |     (Real-time alert, status aplikasi aktif)      |
                  |                                                   |
                  |     Verifikasi & Rotasi Kunci Enkripsi            |
                  |<--------------------------------------------------+ [Push Service (Port Komunikasi)]
                  |                                                   |
                  |     Simulasi Performa Jaringan & Beban            |
                  |<--------------------------------------------------+ [Load Simulator (Port Simulasi)]
```

### 1. Jaringan Zero-Trust Mesh
Seluruh komunikasi berjalan di atas jaringan terowongan (*overlay tunnel*) terenkripsi. Hanya node klien yang sah dengan kunci otentikasi aktif yang dapat bertukar data dengan dashboard pengawas, mencegah pemalsuan identitas klien.

### 2. Peta Tata Letak Tempat Duduk Live (8x5)
Peta grid visual interaktif yang mencerminkan tata letak fisik meja lab komputer. Pengawas dapat memantau status siswa secara langsung:
* ⚪ **Offline**: Klien tidak terhubung.
* 🟢 **Online (Aman)**: Klien terhubung dan berjalan dalam kondisi lingkungan ujian aman.
* 🟡 **Alert (Aktif)**: Klien mendeteksi adanya aktivitas mencurigakan (membuka browser dilarang, dll.).
* 🔴 **Violated (Pernah Melanggar)**: Klien pernah memicu alert sebelumnya pada sesi yang sedang berjalan.

### 3. Pengumpulan Bukti Proaktif
Ketika agen mendeteksi pembukaan proses terlarang atau akses URL mencurigakan, sistem secara otomatis menangkap screenshot layar siswa. Bukti dikirimkan secara aman ke PC Admin untuk diarsipkan secara lokal tanpa melibatkan cloud.

### 4. Simulator Beban Jaringan
Modul bawaan untuk menguji performa throughput jaringan dan kapasitas I/O dashboard. Pengawas dapat memicu pengiriman paket data tiruan dalam jumlah besar secara simultan dari klien untuk menguji keandalan rendering dashboard di bawah beban ekstrem.

---

## 🎓 Alur Kerja Sistem (System Flow)

ReksaFel mengadopsi model **Hybrid Centralized-Mesh** untuk pengawasan terisolasi:
* **Non-Intrusive Agent:** Agen berjalan sebagai *background service* ringan di PC siswa tanpa membebani performa ujian.
* **Asynchronous Concurrency:** Backend Go menggunakan *Goroutines* untuk memproses unggahan bukti tangkapan layar secara paralel tanpa menghambat antrean input/output (I/O) utama.

### Diagram Alur Sistem (System Flowchart)

```mermaid
flowchart TD
    subgraph Klien [PC Siswa / Agent Client]
        A[Agent Background Service] -->|1. Register & Handshake| B(Kanal Komunikasi Terenkripsi)
        A -->|2. Loop Deteksi Proaktif| C{Deteksi Pelanggaran?}
        C -->|Ya| D[Ambil Bukti Layar Screenshot]
        D -->|3. POST Log & Image| B
        C -->|Tidak| A
    end

    subgraph Server [PC Pengawas / Admin Dashboard]
        B -->|4. Terima Data Telemetri| E[API Listener & State Manager]
        E -->|Simpan Bukti Lokal| F[(Penyimpanan Lokal / Disk)]
        E -->|Update Cache| G[Wails Binding Controller]
        
        subgraph UI [Frontend Dashboard]
            G -->|Status Update| H[8x5 Seating Map Grid]
            G -->|Alert Trigger| I[Live Alerts Sidebar]
        end
    end

    subgraph Stress [Modul Simulasi & Stress Test]
        K[Dashboard Admin] -->|5. Trigger Stress Test| L[PC Siswa / Agent Client]
        L -->|6. Kirim Paket Data Paralel| K
        K -->|7. Hitung Throughput / RTT| M[Metrik Grafik & SVG Animasi]
    end
```

---

## 💻 Struktur Kode & Spesifikasi Publik

Repositori ini menyajikan struktur API, skema data telemetri, dan logika enkripsi lokal yang digunakan dalam sistem ReksaFel:

* **`pkg/telemetry/`** — Struktur data, aturan validasi, daftar hitam proses aplikasi (`process.go`), informasi perangkat keras klien (`sysinfo.go`), dan log aktivitas.
* **`pkg/net/`** — Deteksi adapter jaringan (`interface.go`), validasi rute subnet (`routing.go`), dan pengukur latensi TCP.
* **`pkg/api/`** — Spesifikasi router API utama, middleware pencatatan aktivitas, pemetaan endpoint, dan klien pengirim HTTP (`client.go`).
* **`pkg/config/`** — Skema validator file konfigurasi (`validation.go`), pembaca/penulis file lokal (`io.go`), penandatangan integritas payload (`signing.go`), dan mekanisme segel enkripsi AES-256-GCM.

### 🧪 Menjalankan Unit Test
Untuk memastikan keandalan logika kode dan parameter keamanan secara lokal:
```bash
# Menjalankan seluruh pengujian unit
go test -v ./...
```

---

## 🛠️ Stack Teknologi

* **Backend Core:** Go (Golang) — Penanganan konkurensi tingkat tinggi, manajemen soket jaringan, dan kompilasi biner Windows.
* **Frontend Shell:** Vanilla JS, CSS3, & HTML5 — Dirender sebagai aplikasi desktop native OS menggunakan **Wails v2**.
* **Keamanan:** Enkripsi AES-256-GCM untuk mengunci konfigurasi lokal klien, memanfaatkan pengidentifikasi mesin unik berbasis perangkat keras Windows.

---

## 📋 Status Pengembangan & Fitur

Untuk rincian lengkap mengenai status fungsionalitas fitur, visualisasi tata letak navigasi tab dashboard, detail cara kerja sistem di belakang layar, serta riwayat log pembaruan fitur, silakan merujuk pada dokumen pelacak status:

👉 **[Papan Status Fitur & Pengembangan ReksaFel (reksafel_feature_status.md)](reksafel_feature_status.md)**

---

## 💡 Informasi Kolaborasi

Untuk informasi lebih lanjut mengenai evaluasi arsitektur jaringan, integrasi komersial modul Zero-Trust mesh, atau kolaborasi pengembangan teknis sistem ujian laboratorium, silakan menghubungi tim pengembang terkait.

<!-- Readme status checked on 2026-06-25 -->
