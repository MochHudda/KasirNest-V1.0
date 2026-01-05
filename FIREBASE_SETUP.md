# Firebase Setup Guide untuk KasirNest

Panduan lengkap untuk mengatur Firebase backend untuk aplikasi KasirNest.

## 📋 Overview

KasirNest menggunakan Firebase sebagai backend dengan layanan berikut:
- **Firebase Authentication** - Untuk login user
- **Firebase Firestore** - Database untuk produk dan transaksi  
- **Firebase Storage** - Penyimpanan gambar produk (opsional)

## 🚀 Langkah 1: Membuat Firebase Project

### 1.1 Buat Project Baru

1. Buka [Firebase Console](https://console.firebase.google.com)
2. Klik **"Create a project"** atau **"Add project"**
3. Masukkan nama project: `kasirnest-[namaanda]`
4. Pilih apakah ingin menggunakan Google Analytics (opsional)
5. Klik **"Create project"**

### 1.2 Catat Project Information

Setelah project dibuat, catat informasi berikut:
- **Project ID**: Akan digunakan di konfigurasi
- **Web API Key**: Terlihat di Project Settings

## 🔐 Langkah 2: Setup Authentication

### 2.1 Enable Authentication

1. Di Firebase Console, pilih project Anda
2. Klik **"Authentication"** di sidebar kiri
3. Pilih tab **"Sign-in method"**
4. Enable **"Email/Password"** provider:
   - Klik pada "Email/Password"
   - Toggle "Enable"
   - Klik "Save"

### 2.2 Tambah User Admin

1. Pilih tab **"Users"**
2. Klik **"Add user"**
3. Masukkan:
   - **Email**: admin@yourdomain.com (ganti dengan email Anda)
   - **Password**: Minimal 6 karakter, gunakan password kuat
4. Klik **"Add user"**

### 2.3 Optional: Tambah User Kasir

Untuk user kasir tambahan:
1. Ulangi langkah 2.2 dengan email berbeda
2. Atau biarkan admin yang mendaftarkan user baru dari aplikasi

## 🗄 Langkah 3: Setup Firestore Database

### 3.1 Create Firestore Database

1. Klik **"Firestore Database"** di sidebar
2. Klik **"Create database"**
3. Pilih **"Start in production mode"** (recommended)
4. Pilih lokasi database (pilih yang terdekat dengan lokasi Anda):
   - **asia-southeast1** (Singapore) - untuk Indonesia
   - **us-central1** (Iowa) - untuk global
5. Klik **"Done"**

### 3.2 Setup Security Rules

1. Pilih tab **"Rules"** di Firestore
2. Replace rules default dengan rules berikut:

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    // Users can only access their own user document
    match /users/{userId} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    
    // Products - authenticated users can read, only admins can write
    match /products/{document} {
      allow read: if request.auth != null;
      allow write: if request.auth != null && 
        exists(/databases/$(database)/documents/users/$(request.auth.uid)) &&
        get(/databases/$(database)/documents/users/$(request.auth.uid)).data.role == 'admin';
    }
    
    // Transactions - authenticated users can read/write
    match /transactions/{document} {
      allow read, write: if request.auth != null;
    }
    
    // Categories - authenticated users can read, only admins can write
    match /categories/{document} {
      allow read: if request.auth != null;
      allow write: if request.auth != null &&
        exists(/databases/$(database)/documents/users/$(request.auth.uid)) &&
        get(/databases/$(database)/documents/users/$(request.auth.uid)).data.role == 'admin';
    }
    
    // Reports - authenticated users can read/write
    match /reports/{document} {
      allow read, write: if request.auth != null;
    }
  }
}
```

3. Klik **"Publish"**

### 3.3 Create Initial Collections

Buat collection dasar (opsional, akan dibuat otomatis saat aplikasi berjalan):

1. Klik **"Start collection"**
2. Collection ID: `categories`
3. Document ID: `food`
4. Fields:
   ```
   category_id: "food"
   name: "Makanan"
   description: "Produk makanan dan minuman"
   ```
5. Klik **"Save"**

Ulangi untuk categories lain: `electronic`, `fashion`, `health`, `household`, `stationery`, `other`.

## 📁 Langkah 4: Setup Storage (Opsional)

Jika Anda ingin menggunakan upload gambar produk:

### 4.1 Enable Storage

1. Klik **"Storage"** di sidebar
2. Klik **"Get started"**
3. Pilih **"Start in production mode"**
4. Pilih lokasi yang sama dengan Firestore
5. Klik **"Done"**

### 4.2 Setup Storage Rules

1. Pilih tab **"Rules"**
2. Replace dengan rules berikut:

```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    // Allow authenticated users to upload/download product images
    match /products/{productId}/{allPaths=**} {
      allow read: if request.auth != null;
      allow write: if request.auth != null && 
        resource.size < 5 * 1024 * 1024 && // Max 5MB
        resource.contentType.matches('image/.*');
    }
  }
}
```

3. Klik **"Publish"**

## 🔑 Langkah 5: Generate Service Account Key

### 5.1 Create Service Account

1. Klik ikon **gear (⚙️)** di samping "Project Overview"
2. Pilih **"Project settings"**
3. Pilih tab **"Service accounts"**
4. Klik **"Generate new private key"**
5. Klik **"Generate key"** di dialog
6. File JSON akan didownload - **SIMPAN FILE INI DENGAN AMAN**

### 5.2 Extract Credentials

Buka file JSON yang didownload, Anda akan melihat struktur seperti ini:

```json
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "key-id",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-xxxxx@your-project-id.iam.gserviceaccount.com",
  "client_id": "123456789",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-xxxxx%40your-project-id.iam.gserviceaccount.com"
}
```

## ⚙️ Langkah 6: Konfigurasi Aplikasi

### 6.1 Copy Configuration Template

```bash
cp config/app.ini.example config/app.ini
```

### 6.2 Edit Configuration File

Buka `config/app.ini` dan isi dengan informasi dari Firebase:

```ini
[firebase]
# Dari file JSON service account
project_id = your-project-id
private_key_id = key-id-from-json
private_key = "-----BEGIN PRIVATE KEY-----\nYour-Private-Key-Content-Here\n-----END PRIVATE KEY-----\n"
client_email = firebase-adminsdk-xxxxx@your-project-id.iam.gserviceaccount.com
client_id = 123456789
auth_uri = https://accounts.google.com/o/oauth2/auth
token_uri = https://oauth2.googleapis.com/token
auth_provider_x509_cert_url = https://www.googleapis.com/oauth2/v1/certs
client_x509_cert_url = https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-xxxxx%40your-project-id.iam.gserviceaccount.com

# Untuk Firebase Storage
storage_bucket = your-project-id.appspot.com

[app]
name = KasirNest
version = 1.0.0
debug = true
window_width = 1200
window_height = 800
theme = light

[security]
# Generate random encryption key
encryption_key = your-32-character-encryption-key-here
session_timeout = 3600

[database]
auto_backup = true
backup_interval = 24
```

### 6.3 Important Notes

**⚠️ Keamanan Private Key:**
- Private key harus dalam format string dengan `\\n` untuk newlines
- Pastikan backslash escape (`\\n`) bukan actual newlines
- Contoh: `"-----BEGIN PRIVATE KEY-----\\nMIIEvQ...\\n-----END PRIVATE KEY-----\\n"`

**🔐 Generate Encryption Key:**
```bash
# Generate random 32-character key
openssl rand -base64 32
```

## 🧪 Langkah 7: Testing Setup

### 7.1 Test Firebase Connection

1. Build dan jalankan aplikasi:
```bash
go run main.go
```

2. Jika konfigurasi benar, aplikasi akan menampilkan login screen
3. Jika ada error, check logs untuk debugging

### 7.2 Test Login

1. Masukkan email dan password admin yang dibuat di Step 2.2
2. Jika berhasil login, Anda akan masuk ke dashboard
3. Coba navigasi ke berbagai tab untuk memastikan semuanya berfungsi

### 7.3 Test Database Operations

1. Coba tambah produk baru di tab "Produk"
2. Buat transaksi di tab "Transaksi"
3. Check di Firebase Console apakah data tersimpan dengan benar

## 🔧 Troubleshooting

### Error: "Project not found"
- Periksa project_id di konfigurasi
- Pastikan project masih aktif di Firebase Console

### Error: "Permission denied"
- Check Firestore rules
- Pastikan user sudah login
- Verify user role di collection `users`

### Error: "Invalid private key"
- Pastikan format private key benar dengan `\\n`
- Jangan ada extra spaces atau characters
- Re-download service account key jika perlu

### Error: "Authentication failed"
- Check apakah Email/Password provider enabled
- Verify email dan password user di Authentication tab
- Pastikan user tidak di-disable

### Firestore Rules Error
```javascript
// Jika ada error dengan rules, coba rules sederhana ini dulu:
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /{document=**} {
      allow read, write: if request.auth != null;
    }
  }
}
```

## 📊 Langkah 8: Setup User Roles (Opsional)

### 8.1 Create Users Collection

Untuk role-based access:

1. Di Firestore, buat document di collection `users`
2. Document ID: [UID dari Authentication]
3. Fields:
```json
{
  "user_id": "authentication-uid",
  "email": "admin@yourdomain.com",
  "name": "Administrator",
  "role": "admin",
  "created_at": "2024-01-01T00:00:00Z",
  "last_login": "2024-01-01T00:00:00Z"
}
```

### 8.2 Get User UID

Untuk mendapatkan UID user:
1. Go to Authentication > Users
2. Klik user yang ingin diberi role
3. Copy UID yang ditampilkan

## 🎯 Best Practices

### Security
- **Jangan** commit file `config/app.ini` ke version control
- **Gunakan** environment variables untuk production
- **Rotate** service account keys secara berkala
- **Monitor** Firebase usage di console

### Performance  
- **Index** fields yang sering di-query di Firestore
- **Limit** query results untuk performance
- **Cache** data yang jarang berubah
- **Compress** images sebelum upload ke Storage

### Monitoring
- **Enable** Firebase Analytics jika diperlukan
- **Setup** alerts untuk usage quotas
- **Monitor** error logs di Firebase Console
- **Backup** data secara berkala

---

## 🎉 Selesai!

Setelah menyelesaikan semua langkah di atas, aplikasi KasirNest Anda sudah siap digunakan dengan Firebase backend yang fully configured.

Jika mengalami masalah, silakan check:
1. [Troubleshooting section](#troubleshooting) di atas
2. [Firebase Documentation](https://firebase.google.com/docs)
3. [Issues di repository](https://github.com/yourrepo/kasirnest/issues)

**Happy coding! 🚀**