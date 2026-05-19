# API Documentation - Multi Outlet POS & Inventory

Dokumentasi API ini ditulis menggunakan standar **OpenAPI 3.0.3**. Untuk menjaga agar dokumen tetap mudah dibaca, dikelola (_maintainable_), dan *scalable*, dokumentasi ini dipisah berdasarkan konsep **Domain-Driven Design**.

---

## 📂 Struktur Direktori

Seluruh file spesifikasi JSON dibagi menjadi komponen-komponen berikut:

| Folder / File | Deskripsi | Contoh Isi / Domain |
| :--- | :--- | :--- |
| **`openapi.json`** | File **Root / Utama**. Tempat mendeklarasikan Info API, Server, Security Global, dan mapping Path. | `"title": "Multi Outlet POS"`, mapping `$ref` ke folder `paths/`. |
| **`paths/`** | Berisi definisi semua endpoint (Method, Params, Request Body). Dipisah per fitur/domain. | `auth.json`, `sales.json`, `catalog.json`, `inventory.json` |
| **`schemas/`** | Berisi definisi **Model/Struct** JSON agar *reusable* dan tidak perlu ditulis berulang-ulang. | `User`, `Product`, `Branch`, `ErrorResponse` |
| **`responses/`** | Berisi definisi **Format Response Standard** aplikasi. | `GenericSuccess`, `GenericError` |

---

## 🔐 Authentication Strategy

Sistem menggunakan global authentication berbasis **JWT (JSON Web Token)**. Secara *default*, semua endpoint membutuhkan autentikasi kecuali jika di-override dengan `"security": []`.

**Format Header HTTP yang harus dikirim:**
```http
Authorization: Bearer <your_jwt_token_here>
```

---

## 🔓 Public Endpoints (Tanpa Token)

Endpoint berikut bersifat publik dan dapat diakses tanpa mengirimkan header `Authorization`. Dirancang agar _customer_ dapat mengeksplorasi outlet, menu, dan melakukan autentikasi sebelum memesan.

| Fitur / Domain | Method | Endpoint | Use Case / Alasan Public |
| :--- | :--- | :--- | :--- |
| **Auth** | `POST` | `/auth/register` | Customer membuat akun baru. |
| **Auth** | `POST` | `/auth/login` | Mendapatkan token JWT untuk akses API. |
| **Auth** | `POST` | `/auth/forgot-password`| Meminta link/token reset password ke email. |
| **Auth** | `POST` | `/auth/reset-password` | Mengganti password menggunakan token valid. |
| **Branches** | `GET` | `/branches` | Customer ingin melihat lokasi cabang terdekat. |
| **Catalog** | `GET` | `/categories` | Menampilkan daftar kategori produk (Kopi, Non-Kopi, dsb). |
| **Catalog** | `GET` | `/products` | Menampilkan katalog produk (mendukung _filter_ & _search_). |
| **Catalog** | `GET` | `/products/{id}` | Melihat detail produk (harga, gambar, deskripsi bahan). |

> **Contoh Request (Public)**:
> ```bash
> curl -X GET "http://localhost:8080/products?search=kopi&limit=10"
> ```

---

## 🔒 Protected Endpoints (Menggunakan Token)

Seluruh endpoint di bawah ini **WAJIB** menyertakan token JWT yang valid. Endpoint ini mencakup transaksi pengguna pribadi, serta operasional toko (hanya untuk staf).

| Fitur / Domain | Contoh Endpoint Utama | Deskripsi & Akses Ideal |
| :--- | :--- | :--- |
| **User Profile** | `/me`, `/me/password` | Menampilkan/Update profil. Akses: *Semua Role (Customer, Kasir, Admin)* |
| **Carts & Sales** | `/carts`, `/transactions` | Menambah keranjang, checkout pesanan. Akses: *Customer & Kasir* |
| **Inventory** | `/materials`, `/inventories` | Manajemen stok gudang/bahan baku di tiap cabang. Akses: *Admin / Manager* |
| **Restock** | `/restock-requests` | Meminta/menyetujui penambahan stok antar cabang. Akses: *Manager / Kasir* |
| **Shifts** | `/shifts/open`, `/shifts/close`| Buka & tutup shift harian untuk pencatatan kas. Akses: *Kasir / Manager* |
| **Admin Panel** | `/users`, `/products` (`POST`) | Mendaftarkan staf baru, update/tambah menu produk. Akses: *Admin* |

> **Contoh Request (Protected)**:
> ```bash
> curl -X GET http://localhost:8080/me \
>   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
> ```

---

## 🛠️ Cara Menampilkan Secara Visual (Swagger UI)

Karena spesifikasi JSON ini dipisah ke beberapa file, cara terbaik untuk membacanya selama tahap _development_:

1. **Menggunakan VS Code (Lokal)**
   - Install ekstensi **OpenAPI (Swagger) Editor** (oleh 42Crunch).
   - Buka file `openapi.json`.
   - Tekan `Shift + Alt + O` (atau klik tombol Swagger Preview di pojok kanan atas) untuk menampilkan UI interaktif.

2. **Bundle Menjadi 1 File Utuh (Untuk di Deploy)**
   Jika library backend kamu membutuhkan 1 file `swagger.json` utuh, kamu bisa me-render sistem modular ini menjadi 1 file menggunakan `Redocly CLI`:
   ```bash
   npx @redocly/cli bundle backend/api/openapi.json -o backend/api/dist.json
   ```
   *File `dist.json` ini nanti bisa diload langsung di server backend Golang.*
