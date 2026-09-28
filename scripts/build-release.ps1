# PowerShell script untuk mengemas rilis Windows Portable MOJA v0.4.0

$ErrorActionPreference = "Stop"

$Version = "v0.4.0"
$Arch = "windows-amd64"
$ReleaseName = "moja-$Version-$Arch"
$RootDir = Resolve-Path "$PSScriptRoot/.."
$InterpreterDir = "$RootDir/interpreter"
$DistDir = "$RootDir/dist"
$PackageDir = "$DistDir/$ReleaseName"
$ZipPath = "$DistDir/$ReleaseName.zip"

Write-Host "==> Menyiapkan folder rilis $ReleaseName..." -ForegroundColor Cyan

# Bersihkan dist/ jika ada
if (Test-Path $DistDir) {
    Remove-Item -Recurse -Force $DistDir
}

New-Item -ItemType Directory -Force -Path $PackageDir | Out-Null
New-Item -ItemType Directory -Force -Path "$PackageDir/examples" | Out-Null

Write-Host "==> Membangun binary moja.exe..." -ForegroundColor Cyan
Set-Location $InterpreterDir
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "$PackageDir/moja.exe" ./cmd/mojolang

if (-not (Test-Path "$PackageDir/moja.exe")) {
    Write-Error "Gagal membuat binary moja.exe"
}

Write-Host "==> Menyalin LICENSE dan contoh program..." -ForegroundColor Cyan
Copy-Item "$RootDir/LICENSE" -Destination "$PackageDir/LICENSE"
Copy-Item "$InterpreterDir/examples/*.mojo" -Destination "$PackageDir/examples/"

Write-Host "==> Membuat README.txt untuk paket portable..." -ForegroundColor Cyan
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
'@ | Out-File -FilePath "$PackageDir/README.txt" -Encoding utf8

Write-Host "==> Membuat file kompresi ZIP..." -ForegroundColor Cyan
Compress-Archive -Path $PackageDir -DestinationPath $ZipPath -Force

Write-Host "==> Paket rilis $ReleaseName.zip berhasil dibuat!" -ForegroundColor Green
Write-Host "    Lokasi folder : $PackageDir"
Write-Host "    Lokasi archive: $ZipPath"
