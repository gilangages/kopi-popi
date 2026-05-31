# API Documentation - Multi Outlet POS Kopi-Popi & Inventory

Dokumentasi API ini ditulis menggunakan standar **OpenAPI 3.0.3**. Untuk menjaga agar dokumen tetap mudah dibaca, dikelola (_maintainable_), dan *scalable*, dokumentasi ini dipisah berdasarkan konsep **Domain-Driven Design**.

---

## 📂 Struktur Direktori Spesifikasi

Seluruh file spesifikasi JSON dibagi menjadi komponen-komponen berikut:

| Folder / File | Deskripsi | Contoh Isi / Domain |
| :--- | :--- | :--- |
| **`openapi.json`** | File **Root / Utama**. Tempat mendeklarasikan Info API, Server, Security Global, dan mapping Path. | `"title": "Multi Outlet POS"`, mapping `$ref` ke folder `paths/`. |
| **`paths/`** | Berisi definisi semua endpoint (Method, Params, Request Body). Dipisah per fitur/domain. | `auth.json`, `users.json`, `branches.json`, `catalogues.json`, `media.json` |
| **`schemas/`** | Berisi definisi **Model/Struct** JSON agar *reusable* dan tidak perlu ditulis berulang-ulang. | `User`, `Product`, `Branch`, `ErrorResponse` |
| **`responses/`** | Berisi definisi **Format Response Standard** aplikasi. | `GenericSuccess`, `GenericError` |

---

## ✅ Modul / Domain yang Telah Diimplementasikan

Hingga saat ini, sistem *backend* telah memiliki implementasi penuh untuk domain berikut:

1. **Auth (`internal/auth`)**: Menangani otentikasi (Register, Login, Forgot Password).
2. **User (`internal/user`)**: Manajemen pengguna, staf, hak akses, dan manajemen profil pribadi.
3. **Branch (`internal/branch`)**: CRUD untuk manajemen cabang kedai kopi.
4. **Catalog (`internal/catalog`)**: Manajemen Kategori, Material (Bahan Baku), dan Produk (termasuk Resep / Bill of Materials). Dilengkapi proteksi rahasia resep.
5. **Media (`internal/media`)**: Sentralisasi unggah gambar/file dengan proteksi ekstensi (JPG, PNG) dan batas ukuran (5MB). Mendukung folder dinamis.

---

## 🔐 Authentication Strategy

Sistem menggunakan global authentication berbasis **JWT (JSON Web Token)**. Secara *default*, semua endpoint membutuhkan autentikasi kecuali jika di-override dengan `"security": []`.

**Format Header HTTP yang harus dikirim:**
```http
Authorization: Bearer <your_jwt_token_here>
```

---

## 📸 Panduan Upload Gambar (Media Module)

Karena Golang didesain bebas (*unopinionated*), Kopi-Popi menggunakan modul **Media** terpusat untuk segala kebutuhan upload gambar (seperti Profil atau Foto Produk).

**Alur Upload Gambar:**
1. *Frontend* mengirim file fisik gambar (`multipart/form-data`) ke endpoint `POST /uploads`. Parameter `folder` bisa diisi `products` atau `profiles`.
2. *Backend* akan memvalidasi file (maksimal 5MB, format JPG/PNG/WEBP).
3. Jika lolos, file disimpan secara lokal di dalam container Docker di folder `/app/uploads` (di-mapping ke Host `./uploads` lewat Docker Volume agar tidak hilang).
4. *Backend* mengembalikan URL publik (misal: `http://localhost:8080/uploads/products/1234.jpg`).
5. *Frontend* menyimpan URL tersebut, lalu memanggil `POST /products` atau `PATCH /users/me` dengan mengirim data JSON berisi `image_url` atau `profile_picture` menggunakan URL tersebut.

---

## 🔓 Public Endpoints (Tanpa Token)

Endpoint berikut bersifat publik dan dapat diakses tanpa mengirimkan header `Authorization`.

| Fitur / Domain | Method | Endpoint | Use Case |
| :--- | :--- | :--- | :--- |
| **Auth** | `POST` | `/auth/register` | Customer membuat akun baru. |
| **Auth** | `POST` | `/auth/login` | Mendapatkan token JWT. |
| **Auth** | `POST` | `/auth/forgot-password`| Meminta link/token reset password ke email. |
| **Auth** | `POST` | `/auth/reset-password` | Mengganti password menggunakan token. |
| **Branch** | `GET` | `/branches` | Menampilkan cabang publik. Admin bisa request *inactive*. |
| **Catalog** | `GET` | `/categories` | Menampilkan kategori produk. |
| **Catalog** | `GET` | `/products` | Menampilkan daftar produk (tanpa resep). |
| **Catalog** | `GET` | `/products/{id}` | Detail produk. Resep disembunyikan kecuali untuk Admin/Manager. |

---

## 🔒 Protected Endpoints (Menggunakan Token)

Endpoint berikut **WAJIB** menyertakan token JWT yang valid.

| Fitur / Domain | Contoh Endpoint Utama | Akses Ideal |
| :--- | :--- | :--- |
| **User Profile** | `/users/me`, `/users/me/password` | *Semua Role* |
| **Users Mgmt** | `/users`, `/users/managers` | *Admin & Manager* |
| **Branches** | `/branches` (POST, PUT, DELETE) | *Admin* |
| **Catalog (Material)**| `/materials` (CRUD) | *Admin & Manager* |
| **Catalog (Product)**| `/products` (POST, PUT, DELETE) | *Admin* |
| **Media / Uploads** | `/uploads` (`POST`) | *Semua Pengguna Terdaftar* |

---

## 🛠️ Cara Menampilkan Secara Visual (Swagger UI)

Karena spesifikasi JSON dipisah ke beberapa file, cara terbaik untuk membacanya selama tahap _development_:

1. **Menggunakan VS Code (Lokal)**
   - Install ekstensi **OpenAPI (Swagger) Editor** (oleh 42Crunch).
   - Buka file `openapi.json`.
   - Tekan `Shift + Alt + O` untuk menampilkan UI interaktif.

2. **Bundle Menjadi 1 File Utuh (Untuk di Deploy / Redocly)**
   Gunakan Redocly CLI untuk membundel spesifikasi modular menjadi satu file:
   ```bash
   npx @redocly/cli bundle openapi.json -o openapi-bundled.json
   ```
   *File `openapi-bundled.json` ini sudah tersedia di dalam folder ini dan bisa diimpor ke Postman atau Swagger UI.*
