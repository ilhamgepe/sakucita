

---

# 🛡️ Dokumentasi Implementasi Idempotency Key

Dokumentasi ini mengacu pada **IETF Draft: The Idempotency-Key HTTP Header Field** menggunakan **Middleware** dan **Redis** sebagai storage.

## 1. Spesifikasi Header

* **Header Name:** `Idempotency-Key`
* **Format:** `String` (Direkomendasikan **UUID v4**)
* **Storage:** Redis (TTL: 1 Jam)

---

## 2. Alur Logika (Flowchart)

### A. Tahap Inisialisasi (Request Masuk)

1. **Terima Request:** Middleware membaca header `Idempotency-Key` dan menghitung `Hash` dari Request Body.
2. **Atomic Lock (Redis SET NX):** Middleware mencoba menyimpan key ke Redis menggunakan perintah `SET <key> <data> NX EX 300` (TTL awal 5 menit untuk fase locking).
* **Jika Berhasil (OK):** Berarti ini request pertama. Set status menjadi `IN_PROGRESS` dan lanjut ke Logic/Controller.
* **Jika Gagal (Key Exist):** Lanjut ke tahap **Pengecekan Duplikat**.



### B. Tahap Pengecekan Duplikat

1. **Ambil Data dari Redis:**
2. **Validasi Hash:** Bandingkan hash request saat ini dengan hash yang tersimpan di Redis.
* **Hash Berbeda:** Return `422 Unprocessable Content` (Key sama digunakan untuk payload berbeda).


3. **Validasi Status:**
* **Status `IN_PROGRESS`:** Return `409 Conflict` (Transaksi sebelumnya masih berjalan).
* **Status `COMPLETED`:** Return **Cached Response** (Status code & Body yang disimpan di Redis).



### C. Tahap Finalisasi (Setelah Logic Selesai)

1. **Update Redis:** Setelah Controller selesai, update data di Redis:
* Status: `COMPLETED`
* Response Body: `<isi_response>`
* HTTP Status Code: `<status_code>`
* TTL: Perpanjang menjadi `3600` (1 Jam).


2. **Error Handling:** Jika Controller menghasilkan error internal (500), **hapus key** dari Redis agar klien bisa mencoba kembali secara bersih.

---

## 3. Struktur Data Redis (JSON)

Key: `idempotency:key:{UUID}`

```json
{
  "request_hash": "sha256_hash_string",
  "status": "IN_PROGRESS | COMPLETED",
  "response_code": 201,
  "response_body": {
    "transaction_id": "ABC-123",
    "amount": 50000
  }
}

```

---

## 4. Tabel Respon HTTP

| Status Code | Kondisi | Aksi Klien |
| --- | --- | --- |
| **2xx / 201** | Berhasil (Baru atau dari Cache) | Selesai. |
| **400 Bad Request** | Header `Idempotency-Key` tidak ada | Perbaiki request, tambahkan header. |
| **409 Conflict** | Request sedang diproses (Locked) | Tunggu beberapa saat, lalu retry. |
| **422 Unprocessable Content** | Key sama, tapi Body Request berbeda | Gunakan key baru atau perbaiki Body. |

---

## 5. Keuntungan Flow Ini

1. **Race Condition Safe:** Menggunakan `SET NX` (Atomic) memastikan tidak ada dua proses yang berjalan bersamaan untuk satu key.
2. **Efisiensi:** Menghemat resource database (Postgres) karena hasil transaksi sukses langsung diambil dari Redis.
3. **Resilient:** Dengan TTL awal 5 menit (saat locking), jika server crash, key akan otomatis terhapus sehingga tidak terkunci selamanya.

---

**Next Step:** Apakah kamu ingin saya buatkan contoh implementasi *pseudo-code* untuk Middleware ini dalam bahasa pemrograman tertentu (seperti Node.js/Go/Python)?