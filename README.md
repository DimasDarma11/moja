# MOJA — Bahasa Pemrograman Jawa Mojokerto

**MOJA** adalah bahasa pemrograman edukatif berbasis kosakata bahasa Jawa Mojokerto yang terinspirasi dari warisan budaya Majapahit.

> **Sinau coding nganggo basa dhewe.**

---

## 🚀 Panduan Rilis Windows Portable (Tanpa Install Go)

Bagi pengguna Windows yang ingin langsung menggunakan MOJA tanpa mengkompilasi source code:

### 1. Persyaratan Sistem
- **Sistem Operasi**: Windows 10 atau Windows 11 (64-bit).
- **Prasyarat**: Tidak memerlukan Go atau perangkat lunak tambahan lain.

### 2. Cara Download & Ekstrak
1. Unduh arsip zip rilis terbaru (`moja-v0.4.0-windows-amd64.zip`) dari halaman [GitHub Releases](../../releases).
2. Klik kanan file `.zip` lalu pilih **Extract All...** (Ekstrak Semua).
3. Ekstrak ke folder yang diinginkan, misalnya `C:\moja`.

### 3. Cara Menjalankan MOJA di Terminal
1. Buka Terminal (**Command Prompt** atau **PowerShell**) di dalam folder hasil ekstrak `C:\moja`.
   > **Tips**: Pada File Explorer di folder `C:\moja`, ketik `cmd` atau `powershell` di address bar atas lalu tekan **Enter**.
2. Periksa versi MOJA untuk memastikan executable bekerja:
   ```cmd
   .\moja.exe version
   ```
3. Menjalankan file `.mojo`:
   ```cmd
   .\moja.exe run program.mojo
   ```
4. Melihat bantuan perintah CLI:
   ```cmd
   .\moja.exe help
   ```

### 4. Menjalankan Contoh Program yang Tersedia
Di dalam folder rilis terdapat folder `examples/` berisi contoh script MOJA:
```cmd
.\moja.exe run examples\halo.mojo
.\moja.exe run examples\variabel.mojo
.\moja.exe run examples\kondisi.mojo
.\moja.exe run examples\perulangan.mojo
.\moja.exe run examples\logika.mojo
.\moja.exe run examples\array.mojo
```

### 5. Menambahkan MOJA ke Windows PATH (Opsional)
Agar perintah `moja` dapat dipanggil dari folder mana saja di terminal:
1. Tekan `Win + R`, ketik `sysdm.cpl`, lalu tekan **Enter**.
2. Buka tab **Advanced** -> klik **Environment Variables...**.
3. Pada **User variables**, pilih **Path** lalu klik **Edit...**.
4. Klik **New**, lalu masukkan lokasi folder tempat `moja.exe` berada (misal `C:\moja`).
5. Klik **OK** pada seluruh jendela dialog.
6. Buka jendela terminal baru, lalu verifikasi dengan perintah:
   ```cmd
   moja version
   ```

### 6. Pertanyaan Umum & Troubleshooting Error
- **Error: `'moja' is not recognized as an internal or external command`**
  - **Penyebab**: Terminal tidak dapat menemukan lokasi `moja.exe`.
  - **Solusi**: Jika belum menambahkan PATH, sertakan tanda `.\` di depan nama program saat berada di foldernya: `.\moja.exe run file.mojo`. Jika sudah menambah PATH, pastikan Anda telah membuka jendela terminal baru setelah menyimpan konfigurasi PATH.
- **File `.mojo` tidak berjalan saat di-double click**
  - MOJA adalah aplikasi berbasis baris perintah (CLI). Eksekusi script harus dilakukan melalui terminal menggunakan perintah `moja run namafile.mojo`.

---

## 🛠️ Pengembangan & Build dari Source Code

### Struktur Monorepo

```
/                      ← Website dokumentasi (React + Vite + TypeScript)
├── src/               ← Source code React
├── public/            ← Aset statis, termasuk mojolang.wasm & mojoWorker.js
├── interpreter/       ← Interpreter Go (CLI + WASM)
│   ├── cmd/mojolang/  ← Entry point CLI
│   ├── cmd/wasm/      ← Entry point WebAssembly
│   ├── internal/      ← Lexer, Parser, AST, Evaluator
│   └── examples/      ← Contoh program .mojo
├── scripts/           ← Skrip otomatisasi (build-release.ps1)
└── spec/              ← Spesifikasi sintaks MOJA
```

### Persyaratan Build
- **Go 1.22+**

### Kompilasi CLI Manual
```bash
cd interpreter

# Windows
go build -ldflags="-s -w" -o moja.exe ./cmd/mojolang

# Linux / macOS
go build -ldflags="-s -w" -o moja ./cmd/mojolang
```

### Membangun Paket Rilis Windows Portable Automated
Jalankan skrip PowerShell berikut dari root repositori:
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\build-release.ps1
```
Skrip akan mengompilasi binary, menyalin `LICENSE` & `examples/`, membuat `README.txt`, dan mengompresnya menjadi `dist/moja-v0.4.0-windows-amd64.zip`.

### Jalankan Test Suite
```bash
cd interpreter
go test -v ./...
```

---

## ✨ Fitur Bahasa MOJA v0.4

| Fitur | Sintaks / Keyword | Status |
|---|---|---|
| Output teks | `omong("Teks")` | ✅ |
| Variabel | penugasan langsung (`x = 10`) | ✅ |
| Kondisi | `yen` / `liyane` | ✅ |
| Perulangan | `baleni saka tekan` | ✅ |
| Fungsi | `gawe` / `balik` | ✅ |
| Array / List | `[1, 2, 3]`, `arr[0]`, `arr[0] = 99` | ✅ |
| Fungsi Bawaan | `panjang(arr)` / `panjang(str)` | ✅ |
| Auto-Conversion | `+` (String + Number) | ✅ |
| Operator Logika | `lan` (AND), `utawa` (OR), `ora` (NOT) | ✅ |

---

## 👥 Pembuat

- **Dimaz Darma Praja** — Co-Creator
- **Tirta Eka Prawita Sari** — Co-Editor

---

## 📜 Lisensi

MIT License © 2026 Dimaz Darma Praja and Tirta Eka Prawita Sari

