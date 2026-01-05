# KasirNest - Sistem Kasir Desktop Modern

KasirNest adalah aplikasi kasir (Point of Sale) desktop yang dibuat dengan Go, Fyne UI framework, dan Firebase backend. Aplikasi ini dirancang untuk bisnis kecil hingga menengah yang membutuhkan sistem kasir yang handal, aman, dan mudah digunakan.

![KasirNest Logo](assets/logo.png)

## ✨ Fitur Utama

- 🔐 **Sistem Login Aman** - Autentikasi menggunakan Firebase Auth
- 📦 **Manajemen Produk** - CRUD produk dengan support gambar
- 💳 **Transaksi POS** - Interface kasir yang intuitif dan responsif
- 📊 **Laporan Penjualan** - Laporan harian, mingguan, dan bulanan
- 🎨 **UI Modern** - Antarmuka yang bersih dan mudah digunakan
- 🔒 **Keamanan Tinggi** - Binary obfuscation dan proteksi anti-inspeksi
- 🌐 **Cloud Database** - Data tersimpan aman di Firebase Firestore
- 📱 **Cross Platform** - Berjalan di Windows, macOS, dan Linux

## 🛠 Teknologi yang Digunakan

- **Go 1.21+** - Bahasa pemrograman utama
- **Fyne v2** - Framework UI desktop
- **Firebase Auth** - Sistem autentikasi
- **Firebase Firestore** - Database NoSQL
- **Firebase Storage** - Penyimpanan gambar produk
- **Garble** - Obfuscation untuk keamanan

## 📋 Persyaratan Sistem

### Minimum System Requirements
- **OS**: Windows 10, macOS 10.15, Ubuntu 18.04 atau yang lebih baru
- **RAM**: 4GB (8GB direkomendasikan)
- **Storage**: 100MB ruang kosong
- **Network**: Koneksi internet untuk sinkronisasi data

### Development Requirements
- **Go**: 1.21 atau lebih baru
- **Git**: Untuk version control
- **Firebase Project**: Akun Google dengan akses Firebase Console

## 🚀 Instalasi

### Untuk Pengguna (Binary Release)

1. **Download** aplikasi dari halaman [Releases](https://github.com/yourrepo/kasirnest/releases)
2. **Extract** file ZIP/TAR.GZ yang didownload
3. **Copy** file `app.ini.example` ke `app.ini` dalam folder `config`
4. **Edit** file `config/app.ini` dengan konfigurasi Firebase Anda
5. **Jalankan** aplikasi dengan double-click atau dari terminal

### Untuk Developer (Build from Source)

```bash
# 1. Clone repository
git clone https://github.com/yourrepo/kasirnest.git
cd kasirnest

# 2. Install dependencies
go mod download

# 3. Install build tools
go install mvdan.cc/garble@latest
go install fyne.io/fyne/v2/cmd/fyne@latest

# 4. Copy dan edit konfigurasi
cp config/app.ini.example config/app.ini
# Edit config/app.ini dengan konfigurasi Firebase Anda

# 5. Build aplikasi
# Linux/macOS:
chmod +x build/build.sh
./build/build.sh

# Windows:
build\\build.bat
```

## ⚙ Konfigurasi

### Firebase Setup

Sebelum menjalankan aplikasi, Anda harus setup project Firebase. Lihat panduan lengkap di [FIREBASE_SETUP.md](FIREBASE_SETUP.md).

### Konfigurasi Aplikasi

Edit file `config/app.ini` dengan konfigurasi yang sesuai:

```ini
[firebase]
project_id = your-firebase-project-id
private_key_id = your-private-key-id
private_key = "-----BEGIN PRIVATE KEY-----\nYour-Private-Key-Content-Here\n-----END PRIVATE KEY-----\n"
client_email = your-service-account@your-project-id.iam.gserviceaccount.com
# ... dan seterusnya

[app]
name = KasirNest
version = 1.0.0
debug = false
window_width = 1200
window_height = 800
theme = light

[security]
encryption_key = your-encryption-key-here
session_timeout = 3600
```

## 📖 Panduan Penggunaan

### Login
1. Jalankan aplikasi KasirNest
2. Masukkan email dan password yang terdaftar di Firebase Auth
3. Klik tombol "Masuk"

### Manajemen Produk
1. Pilih tab "Produk" di dashboard
2. Untuk menambah produk baru, klik "Tambah Produk"
3. Isi informasi produk (nama, harga, stok, kategori, barcode)
4. Klik "Simpan"

### Transaksi Kasir
1. Pilih tab "Transaksi" 
2. Di tab "Kasir (POS)", cari produk dengan nama atau scan barcode
3. Produk akan ditambahkan ke keranjang
4. Pilih metode pembayaran
5. Klik "Proses Pembayaran"

### Laporan
1. Pilih tab "Laporan"
2. Pilih rentang tanggal dan jenis laporan
3. Klik "Generate" untuk membuat laporan
4. Gunakan "Export" untuk menyimpan laporan

## 🏗 Arsitektur Proyek

```
KasirNest/
├── main.go                 # Entry point aplikasi
├── go.mod                  # Go module definition
├── .gitignore             # Git ignore rules
│
├── config/                # Konfigurasi
│   ├── app.ini           # File konfigurasi (git-ignored)
│   ├── app.ini.example   # Template konfigurasi
│   └── config.go         # Config loader
│
├── models/               # Data models
│   ├── user.go          # Model user
│   ├── product.go       # Model produk
│   ├── transaction.go   # Model transaksi
│   ├── category.go      # Model kategori
│   ├── report.go        # Model laporan
│   └── errors.go        # Error definitions
│
├── firebase/            # Firebase integration
│   ├── client.go        # Firebase client
│   ├── auth.go          # Authentication service
│   ├── firestore.go     # Firestore service
│   └── storage.go       # Storage service
│
├── ui/                  # User interface
│   ├── login.go         # Login screen
│   ├── dashboard.go     # Main dashboard
│   ├── products.go      # Products management
│   ├── transactions.go  # Transactions/POS
│   └── reports.go       # Reports screen
│
├── utils/               # Utility functions
│   ├── validator.go     # Input validation
│   ├── formatter.go     # Data formatting
│   └── crypto.go        # Encryption utilities
│
├── internal/            # Internal packages
│   └── secure.go        # Security utilities
│
├── assets/              # Application assets
│   └── logo.png         # Application logo
│
└── build/               # Build scripts
    ├── build.sh         # Unix build script
    └── build.bat        # Windows build script
```

## 🔧 Development

### Running in Development Mode

```bash
# Set debug mode in config/app.ini
debug = true

# Run without building
go run main.go

# Or build and run
go build -o kasirnest .
./kasirnest
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./utils -v
```

### Adding New Features

1. **Models**: Tambah struct di folder `models/`
2. **Firebase**: Extend service di folder `firebase/`
3. **UI**: Buat komponen UI baru di folder `ui/`
4. **Utils**: Tambah utility functions di folder `utils/`

## 📦 Build & Distribution

### Build untuk Development

```bash
# Build sederhana
go build -o kasirnest .

# Build dengan optimasi
go build -ldflags="-s -w" -o kasirnest .
```

### Build untuk Production

```bash
# Linux/macOS
./build/build.sh --all-platforms

# Windows
build\build.bat --all-platforms

# Build dengan obfuscation
./build/build.sh --obfuscate
```

### Packaging

Script build otomatis membuat package distribusi:
- **Windows**: ZIP file dengan executable dan dependencies
- **Linux/macOS**: TAR.GZ file dengan binary dan assets
- **Include**: Binary, config template, assets, documentation

## 🔒 Keamanan

KasirNest menerapkan beberapa tingkat keamanan:

### 1. Data Security
- ✅ Komunikasi terenkripsi dengan Firebase
- ✅ API keys tidak di-hardcode dalam binary
- ✅ Session timeout untuk keamanan login
- ✅ Input validation untuk mencegah injection

### 2. Binary Protection  
- ✅ Code obfuscation menggunakan garble
- ✅ Binary stripping untuk mengurangi ukuran
- ✅ Anti-debugging measures (basic)

### 3. Firebase Security
- ✅ Authentication required untuk semua operasi
- ✅ Firestore security rules
- ✅ Role-based access control

## ⚠ Troubleshooting

### Masalah Umum

**1. Aplikasi tidak bisa login**
- Periksa konfigurasi Firebase di `config/app.ini`
- Pastikan user terdaftar di Firebase Auth
- Cek koneksi internet

**2. Error "Firebase not initialized"**
- Pastikan semua field Firebase di konfigurasi terisi
- Periksa format private key (harus include \\n untuk newlines)
- Verifikasi project ID dan credentials

**3. Produk tidak tersimpan**
- Cek Firestore security rules
- Pastikan user memiliki permission write
- Lihat log error di console

**4. Build error**
- Update Go ke versi terbaru
- Jalankan `go mod tidy`
- Install ulang build tools

### Debug Mode

Untuk debugging, set `debug = true` di `config/app.ini`. Ini akan:
- Menampilkan log lebih detail
- Disable beberapa security measures
- Enable development features

### Log Files

Aplikasi menyimpan log di:
- **Windows**: `%APPDATA%/KasirNest/logs/`
- **macOS**: `~/Library/Application Support/KasirNest/logs/`
- **Linux**: `~/.local/share/KasirNest/logs/`

## 🤝 Contributing

Kami menyambut kontribusi dari community! 

### How to Contribute

1. **Fork** repository ini
2. **Create** feature branch (`git checkout -b feature/AmazingFeature`)
3. **Commit** changes (`git commit -m 'Add some AmazingFeature'`)
4. **Push** to branch (`git push origin feature/AmazingFeature`)
5. **Open** Pull Request

### Development Guidelines

- Ikuti Go coding standards
- Tambahkan tests untuk fitur baru
- Update dokumentasi jika diperlukan
- Pastikan build berhasil di semua platform

## 📄 License

Distributed under the MIT License. See `LICENSE` file for more information.

## 👥 Team

- **Developer**: [Your Name]
- **Email**: your.email@example.com
- **Website**: https://yourwebsite.com

## 🙏 Acknowledgments

- [Fyne](https://fyne.io/) - Go-based UI framework
- [Firebase](https://firebase.google.com/) - Backend services
- [Garble](https://github.com/burrowers/garble) - Go code obfuscation
- [Go Team](https://golang.org/) - Programming language

---

**KasirNest** - Modern POS System for Modern Business 🚀