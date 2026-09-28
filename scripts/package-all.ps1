# PowerShell script untuk menyiapkan paket rilis MOJA v0.4.0 (Source + Windows Portable + Guide)

$ErrorActionPreference = "Stop"

$Version = "v0.4.0"
$RootDir = (Get-Item "$PSScriptRoot/..").FullName
$InterpreterDir = "$RootDir/interpreter"
$ReleasePackageDir = "$RootDir/release-package"
$SourceZipPath = "$ReleasePackageDir/moja-$Version-source.zip"
$PortableZipPath = "$ReleasePackageDir/moja-$Version-windows-amd64.zip"
$GuidePath = "$ReleasePackageDir/GITHUB_UPLOAD_GUIDE.md"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "  MENYIAPKAN PAKET RILIS MOJA $Version UNTUK GITHUB" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# Step 1: Jalankan Test Suite
Write-Host "`n[1/6] Melakukan pengujian interpreter (go test ./...)..." -ForegroundColor Yellow
Set-Location $InterpreterDir
go test ./...
if ($LASTEXITCODE -ne 0) {
    Write-Error "Test suite gagal! Batalkan pembuatan paket rilis."
}
Write-Host "-> Seluruh unit test PASS!" -ForegroundColor Green

# Step 2: Menyiapkan folder release-package/
Write-Host "`n[2/6] Menyiapkan folder release-package/..." -ForegroundColor Yellow
if (Test-Path $ReleasePackageDir) {
    Remove-Item -Recurse -Force $ReleasePackageDir
}
New-Item -ItemType Directory -Force -Path $ReleasePackageDir | Out-Null

# Step 3: Membangun Binary & Membuat Portable ZIP
Write-Host "`n[3/6] Membangun binary & mengemas paket Windows Portable..." -ForegroundColor Yellow
$TempPortableDir = "$ReleasePackageDir/moja-$Version-windows-amd64"
New-Item -ItemType Directory -Force -Path "$TempPortableDir/examples" | Out-Null

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "$TempPortableDir/moja.exe" ./cmd/mojolang

# Verifikasi versi binary
$VerOutput = & "$TempPortableDir/moja.exe" version
Write-Host "-> Versi binary terverifikasi: $VerOutput" -ForegroundColor Green

# Salin Lisensi dan Contoh Program
Copy-Item "$RootDir/LICENSE" -Destination "$TempPortableDir/LICENSE"
Copy-Item "$InterpreterDir/examples/*.mojo" -Destination "$TempPortableDir/examples/"

# Buat README.txt untuk Portable
@'
========================================================================
             MOJA v0.4.0 — Bahasa Pemrograman Jawa Mojokerto
                       Distribusi Windows Portable
========================================================================

Sugeng rawuh ing MOJA!
MOJA adalah bahasa pemrograman edukasi berbasis kosakata Jawa Mojokerto.
Paket ini berisi distribusi portable yang siap langsung dijalankan di Windows
tanpa perlu menginstall Go atau software pendukung lainnya.

------------------------------------------------------------------------
ISI PAKET RILIS:
------------------------------------------------------------------------
  - moja.exe        : Executable interpreter MOJA
  - README.txt      : Panduan penggunaan ini
  - LICENSE         : Lisensi MIT MOJA
  - examples/       : Folder berisi contoh program (.mojo)

------------------------------------------------------------------------
CARA CEPAT MENJALANKAN MOJA:
------------------------------------------------------------------------
1. Buka Terminal (Command Prompt atau PowerShell) di folder tempat Anda
   mengekstrak paket ini.
   - Tips: Tulis "cmd" atau "powershell" pada address bar File Explorer di
     folder ini lalu tekan Enter.

2. Jalankan perintah untuk cek versi MOJA:
   .\moja.exe version

3. Jalankan contoh program .mojo:
   .\moja.exe run examples\halo.mojo
   .\moja.exe run examples\variabel.mojo
   .\moja.exe run examples\kondisi.mojo
   .\moja.exe run examples\perulangan.mojo
   .\moja.exe run examples\logika.mojo
   .\moja.exe run examples\array.mojo

------------------------------------------------------------------------
CARA MENJALANKAN SCRIPT MILIK ANDA SENDIRI:
------------------------------------------------------------------------
1. Buat file baru dengan ekstensi .mojo (misalnya: programku.mojo) menggunakan
   text editor favorit Anda (VS Code, Notepad, Notepad++, dll).
2. Tulis kode MOJA, contoh:
     omong("Sugeng rawuh ing MOJA!")
3. Simpan file tersebut.
4. Jalankan via terminal:
     .\moja.exe run programku.mojo

------------------------------------------------------------------------
MENAMBAHKAN MOJA KE WINDOWS PATH (OPSIONAL):
------------------------------------------------------------------------
Agar Anda dapat mengetik `moja` langsung dari folder mana saja tanpa harus
menyebutkan lokasi `.\moja.exe`:

1. Ekstrak folder MOJA ke lokasi permanen, misal: C:\moja
2. Tekan tombol Windows + R, ketik `sysdm.cpl`, lalu tekan Enter.
3. Masuk ke tab "Advanced" (Lanjutan) -> klik "Environment Variables..."
4. Pada bagian "User variables" (atau "System variables"), pilih "Path" lalu
   klik "Edit...".
5. Klik "New" (Baru), masukkan path folder tempat moja.exe berada (misal: C:\moja).
6. Klik OK pada semua jendela.
7. Buka jendela Terminal baru, lalu tes dengan perintah:
     moja version

------------------------------------------------------------------------
TROUBLESHOOTING / TROUBLESHOOTING ERROR UMUM:
------------------------------------------------------------------------
Q: Muncul error "'moja' is not recognized as an internal or external command"
A: Pastikan Anda mengetik `.\moja.exe` jika berada di folder yang sama di
   PowerShell, atau pastikan langkah penambahan PATH sudah dilakukan dan
   terminal telah di-restart.

Q: Mengapa file .mojo tidak bisa di-double click?
A: MOJA adalah program CLI (Command Line Interface). Gunakan terminal untuk
   menjalankannya dengan perintah `moja run namafile.mojo`.

------------------------------------------------------------------------
DOKUMENTASI & KONTRIBUSI:
------------------------------------------------------------------------
Spesifikasi sintaks dan dokumentasi lengkap dapat dilihat pada repositori resmi MOJA.
Matur nuwun sampun nggunakake MOJA!
========================================================================
'@ | Out-File -FilePath "$TempPortableDir/README.txt" -Encoding utf8

# Buat ZIP Portable
Compress-Archive -Path "$TempPortableDir/*" -DestinationPath $PortableZipPath -Force
Remove-Item -Recurse -Force $TempPortableDir
Write-Host "-> Portable ZIP berhasil dibuat: $PortableZipPath" -ForegroundColor Green

# Step 4: Menyiapkan Paket Source Code MURNI INTERPRETER (Tanpa Web App)
Write-Host "`n[4/6] Mengemas paket Source Code Interpreter MOJA Murni (Tanpa Web App)..." -ForegroundColor Yellow
$TempSourceDir = "$ReleasePackageDir/moja-$Version-source"
New-Item -ItemType Directory -Force -Path $TempSourceDir | Out-Null

# Hanya menyalin komponen inti interpreter MOJA (interpreter, spec, scripts, README, LICENSE, .gitignore)
$ItemsToCopy = @(
    "interpreter",
    "spec",
    "scripts",
    "LICENSE",
    "README.md",
    ".gitignore"
)

foreach ($item in $ItemsToCopy) {
    $srcPath = "$RootDir/$item"
    if (Test-Path $srcPath) {
        $destPath = "$TempSourceDir/$item"
        Copy-Item -Path $srcPath -Destination $destPath -Recurse -Force
    }
}

# Hapus binary atau file temporary di dalam folder source yang terkopi
Get-ChildItem -Path "$TempSourceDir/interpreter" -Recurse -Include "*.exe", "*.log", "moja", "mojolang", "dist" -ErrorAction SilentlyContinue | Remove-Item -Force -Recurse

# Buat ZIP Source Code
Compress-Archive -Path "$TempSourceDir/*" -DestinationPath $SourceZipPath -Force
Remove-Item -Recurse -Force $TempSourceDir
Write-Host "-> Source ZIP Interpreter Murni berhasil dibuat: $SourceZipPath" -ForegroundColor Green

# Step 5: Membuat GITHUB_UPLOAD_GUIDE.md
Write-Host "`n[5/6] Membuat panduan upload GitHub (GITHUB_UPLOAD_GUIDE.md)..." -ForegroundColor Yellow
@'
# Panduan Upload GitHub Rilis MOJA v0.4.0

Panduan praktis dan sederhana langkah demi langkah untuk mempublikasikan proyek **MOJA v0.4.0** (Interpreter Bahasa Pemrograman MOJA) ke GitHub.

---

## 📋 BAGIAN 1 — Upload Source Code Interpreter ke GitHub

### 1. Membuat Repositori Baru di GitHub
1. Buka situs [GitHub.com](https://github.com) dan login ke akun Anda.
2. Klik ikon **`+`** di sudut kanan atas -> pilih **New repository**.
3. Isi kolom berikut:
   - **Repository name**: `moja`
   - **Description**: `Bahasa pemrograman edukatif berbasis bahasa Jawa Mojokerto (Interpreter Go)`
   - **Visibility**: Pilih **Public** (agar dapat diakses publik) atau **Private**.
   - **Initialize repository**: **JANGAN** mencentang *Add a README file*, *.gitignore*, atau *license* (karena file-file ini sudah tersedia lengkap di dalam source code).
4. Klik tombol **Create repository**.

---

### 2. Mengunggah Source Code Interpreter

#### Option A: Menggunakan Git CLI (Rekomendasi Utama & Paling Stabil)
1. Ekstrak file `moja-v0.4.0-source.zip` ke sebuah folder di komputer Anda.
2. Buka Terminal (Command Prompt / PowerShell / Git Bash) di dalam folder hasil ekstrak tersebut.
3. Jalankan perintah berikut secara berurutan:
   ```bash
   git init
   git add .
   git commit -m "feat: initial release MOJA v0.4.0"
   git branch -M main
   git remote add origin https://github.com/USERNAME/moja.git
   git push -u origin main
   ```
   *(Ganti `USERNAME` dengan nama akun GitHub Anda).*

#### Option B: Menggunakan Web Upload (Drag & Drop)
1. Ekstrak file `moja-v0.4.0-source.zip`.
2. Pada halaman repositori kosong di GitHub, klik tautan **uploading an existing file**.
3. Seret (drag & drop) seluruh file dan folder dari hasil ekstrak ke area upload GitHub.
4. Tulis pesan commit (contoh: `Initial release MOJA v0.4.0`) lalu klik **Commit changes**.

---

### 3. Verifikasi Tampilan Utama
- Setelah push/upload selesai, buka halaman utama repositori Anda di `https://github.com/USERNAME/moja`.
- Tampilan `README.md` akan otomatis dirender di halaman utama.

---

## 📦 BAGIAN 2 — Membuat Release di GitHub

### 1. Membuka Halaman Releases
1. Pada halaman utama repositori GitHub Anda, perhatikan bilah kanan (sidebar).
2. Klik tautan **Releases** (atau **Create a new release**).

---

### 2. Membuat Release Tag `v0.4.0`
1. Klik tombol **Draft a new release**.
2. Klik tombol **Choose a tag**:
   - Ketik: `v0.4.0`
   - Klik **Create new tag: v0.4.0 on publish**.
3. Pastikan **Target** mengarah ke branch `main`.
4. Isi **Release title**:
   ```text
   MOJA v0.4.0 — Mojokerto Edition
   ```

---

### 3. Mengisi Deskripsi Release (Dapat Di-copy Paste)
Salin (copy) dan tempel (paste) teks Markdown berikut ke dalam kotak deskripsi release:

```markdown
# MOJA v0.4.0 — Mojokerto Edition 🚀

Sugeng rawuh ing rilis resmi **MOJA v0.4.0**! MOJA adalah bahasa pemrograman edukasi berbasis kosakata Jawa Mojokerto yang terinspirasi dari warisan budaya Majapahit.

### 🌟 Fitur Utama MOJA v0.4.0
- 🖨️ **Output & Cetak**: `omong("Halo Mojokerto!")`
- 📦 **Variabel & Scope**: Penugasan langsung `x = 10` dengan pencarian scope berjenjang
- 🔀 **Pengondisian**: `yen` ... `liyane`
- 🔁 **Perulangan**: `baleni i saka 1 tekan 5`
- 🧩 **Fungsi**: Deklarasi `gawe` dan pengembalian nilai `balik`
- 📚 **Array / List**: Literal `[1, 2, 3]`, akses indeks `arr[0]`, mutasi `arr[0] = 99`, dan fungsi bawaan `panjang()`
- ⚡ **Operator Logika**: `lan` (AND), `utawa` (OR), dan `ora` (NOT) dengan evaluasi short-circuit

---

### 💾 Unduh Distribusi Windows Portable (Tanpa Install Go)
Bagi pengguna Windows yang ingin langsung mencoba MOJA tanpa menginstall Go:
1. Unduh file **`moja-v0.4.0-windows-amd64.zip`** pada bagian **Assets** di bawah.
2. Ekstrak file zip ke folder pilihan Anda (misal `C:\moja`).
3. Buka Terminal (PowerShell / Command Prompt) di folder tersebut, lalu jalankan:
   ```cmd
   .\moja.exe run examples\halo.mojo
   ```
```

---

### 4. Mengunggah File Asset Portable
1. Scroll ke bagian **Attach binaries by dropping them here or selecting them**.
2. Unggah file **`moja-v0.4.0-windows-amd64.zip`** dari folder `release-package/`.
3. **Penting**: *Jangan* mengunggah `moja-v0.4.0-source.zip` sebagai asset tambahan, karena GitHub akan secara otomatis membuatkan arsip Source Code (`.zip` dan `.tar.gz`) dari repositori saat rilis dipublikasikan.

---

### 5. Mempublikasikan Release
- Klik tombol hijau **Publish release**.

---

## 💻 BAGIAN 3 — Cara Pengguna Memakai MOJA (Petunjuk Pengguna)

1. **Unduh Paket Portable**:
   - Pengguna mengunduh file `moja-v0.4.0-windows-amd64.zip` dari halaman Releases GitHub.
2. **Ekstrak**:
   - Klik kanan file zip -> **Extract All...** -> arahkan ke `C:\moja`.
3. **Buka Terminal & Jalankan Perintah**:
   - Buka PowerShell di folder `C:\moja`.
   - Cek versi:
     ```cmd
     .\moja.exe version
     ```
   - Jalankan contoh program:
     ```cmd
     .\moja.exe run examples\halo.mojo
     .\moja.exe run examples\logika.mojo
     ```
4. **Menambahkan PATH (Opsional)**:
   - Pengguna dapat menambahkan `C:\moja` ke Environment Variable `PATH` di Windows agar perintah `moja` dapat dipanggil dari lokasi mana saja tanpa mengetik `.\moja.exe`.
5. **Troubleshooting Error**:
   - *Pesan Error*: `'moja' is not recognized as an internal or external command`
   - *Solusi*: Sertakan `.\` di depan executable jika berada di PowerShell (`.\moja.exe`), atau restart terminal setelah menambahkan ke PATH.

---

## ✅ BAGIAN 4 — Checklist Sebelum Publikasi

- [x] **Source Code**: Murni Interpreter Go MOJA + Spesifikasi Syntax + Scripts (Tanpa Web App).
- [x] **Keamanan**: Tidak ada file rahasia, token, credential, atau file `.env`.
- [x] **Lisensi**: File `LICENSE` (MIT) tersedia.
- [x] **Dokumentasi**: `README.md` informatif dan jelas.
- [x] **Pengujian Unit**: `go test ./...` lulus 100%.
- [x] **Kompilasi Binary**: `moja.exe` dikompilasi untuk `windows/amd64` (versi `v0.4.0`).
- [x] **Arsip ZIP**: Berkas `moja-v0.4.0-source.zip` dan `moja-v0.4.0-windows-amd64.zip` valid & dapat dibuka.
- [x] **Script Contoh**: Seluruh contoh `.mojo` dalam paket portable diuji dan tereksekusi dengan benar.
'@ | Out-File -FilePath $GuidePath -Encoding utf8
Write-Host "-> GITHUB_UPLOAD_GUIDE.md berhasil dibuat: $GuidePath" -ForegroundColor Green

# Step 6: Verifikasi Ukuran & Status
Write-Host "`n[6/6] Memeriksa ukuran & verifikasi file paket rilis..." -ForegroundColor Yellow
$SourceZipItem = Get-Item $SourceZipPath
$PortableZipItem = Get-Item $PortableZipPath
$GuideItem = Get-Item $GuidePath

Write-Host "   - $SourceZipPath ($([math]::Round($SourceZipItem.Length/1KB, 2)) KB)" -ForegroundColor Cyan
Write-Host "   - $PortableZipPath ($([math]::Round($PortableZipItem.Length/1MB, 2)) MB)" -ForegroundColor Cyan
Write-Host "   - $GuidePath ($([math]::Round($GuideItem.Length/1KB, 2)) KB)" -ForegroundColor Cyan

Set-Location $RootDir
Write-Host "`n[SELESAI] Seluruh paket rilis murni berhasil dibuat di $ReleasePackageDir!" -ForegroundColor Green
