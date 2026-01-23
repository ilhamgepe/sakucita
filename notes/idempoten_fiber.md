Ini dia rangkuman **Final Flow** menggunakan **Go Fiber Idempotency Middleware + Redis** dalam format Markdown. Kamu bisa simpan ini sebagai panduan teknis implementasi di projekmu.

---

# 🚀 Implementasi Idempotency: Go Fiber + Redis Storage

Dokumentasi ini menjelaskan alur kerja otomatis dari middleware Go Fiber dalam menangani request duplikat untuk memastikan operasi (seperti transaksi) hanya dieksekusi tepat satu kali (*Exactly-once execution*).

## 1. Spesifikasi Teknis

* **Framework:** Go Fiber (v2)
* **Middleware:** `idempotency`
* **Storage:** Redis (v3)
* **Header Target:** `Idempotency-Key` (Wajib UUID v4 dari Klien)
* **TTL (Lama Cache):** 1 Jam

---

## 2. Alur Kerja (Step-by-Step)

### A. Tahap Request Pertama (New Request)

1. **Intercept:** Middleware menangkap `Idempotency-Key` dari header.
2. **Redis Check:** Middleware mengecek ke Redis. Karena data kosong, middleware membuat entry baru dengan status **`LOCKED`**.
3. **Execution:** Request diteruskan ke **Controller/Handler** kamu.
4. **Capture:** Setelah Controller selesai, middleware menangkap (capture) `Status Code` dan `Response Body`.
5. **Commit:** Middleware memperbarui data di Redis (mengubah status dari `LOCKED` menjadi hasil response tadi) dan menyimpannya selama **1 Jam**.

### B. Tahap Request Berulang (Duplicate Request)

1. **Intercept:** Middleware menemukan `Idempotency-Key` yang sama di Redis.
2. **Pengecekan Status:**
* **Jika Status masih `LOCKED`:** Berarti request pertama masih diproses. Middleware otomatis mengembalikan **`409 Conflict`**.
* **Jika Status sudah `SUCCESS`:** Middleware mengambil data dari Redis dan langsung mengembalikan **Cached Response** ke klien tanpa menyentuh Controller kamu lagi.



### C. Tahap Error Handling

* **Server Error (500):** Jika proses di Controller gagal, middleware secara cerdas **tidak akan menyimpan cache**, sehingga klien diperbolehkan melakukan *retry* dengan kunci yang sama.
* **Payload Berbeda:** (Opsional) Disarankan klien menyertakan hash dari body jika ingin memvalidasi integritas data lebih ketat.

---

## 3. Contoh Konfigurasi Go Fiber

```go
app.Use(idempotency.New(idempotency.Config{
    Next: func(c *fiber.Ctx) bool {
        // Hanya jalankan untuk method penulisan data
        return c.Method() != fiber.MethodPost && c.Method() != fiber.MethodPatch
    },
    Lifetime:      1 * time.Hour,      // Sesuai kebutuhanmu
    KeyHeader:     "Idempotency-Key", // Header yang dipantau
    Storage:       redisStore,        // Wajib diarahkan ke Redis
}))

```

---

## 4. Keuntungan Utama bagi Backend Kamu

1. **Otomatis:** Kamu tidak perlu menulis kode `IF-ELSE` di dalam controller untuk mengecek cache.
2. **Hemat Resource:** Database Postgres tidak akan terbebani oleh request *retry* karena sudah diputus di level middleware.
3. **Keamanan Transaksi:** Menghindari *double payment* atau duplikasi data akibat masalah jaringan antara klien dan server.

---

## 5. Ringkasan Status Code

| Status | Makna |
| --- | --- |
| **200/201** | Transaksi berhasil diproses (atau diambil dari cache). |
| **409 Conflict** | Request sedang berjalan, jangan kirim ulang dulu. |
| **400 Bad Request** | Klien lupa mengirimkan header `Idempotency-Key`. |

---

Dengan flow ini, sistem kamu sudah standar industri (seperti Stripe atau PayPal) dalam menangani idempotensi! Ada lagi bagian yang ingin kamu pertegas?