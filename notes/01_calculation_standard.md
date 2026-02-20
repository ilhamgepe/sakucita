Bagus. Ini memang perlu didokumentasikan supaya semua developer tidak bikin kesalahan fatal di layer uang.

Di bawah ini adalah versi catatan yang bisa langsung kamu taruh di repo (README atau docs/finance.md).

---

# 💰 Financial Calculation Standard (Sakucita)

Dokumen ini menjelaskan standar perhitungan uang dan fee di sistem.

⚠️ Aturan ini wajib diikuti.
⚠️ Tidak boleh menggunakan float untuk uang.

---

# 1️⃣ Currency Standard

Saat ini sistem hanya mendukung:

```
IDR (Indonesian Rupiah)
```

Karena IDR tidak memiliki desimal:

```
Minor unit = 1
```

Artinya:

* Rp 10.000 disimpan sebagai `10000`
* Semua amount disimpan dalam `BIGINT`
* Tidak pernah menggunakan NUMERIC atau FLOAT

---

# 2️⃣ Database Rules

Semua kolom uang wajib:

```
BIGINT
```

Contoh:

```
amount
gateway_fee_amount
platform_fee_amount
fee_amount
net_amount
```

Tidak boleh:

```
FLOAT
NUMERIC
DECIMAL
```

---

# 3️⃣ Percentage Storage (Basis Point)

Persentase tidak disimpan sebagai desimal.

Sistem menggunakan **Basis Point (BPS)**.

Definisi:

```
1% = 100 bps
100% = 10000 bps
```

Konversi:

| Percent | Disimpan di DB |
| ------- | -------------- |
| 0.7%    | 70             |
| 1.5%    | 150            |
| 2%      | 200            |
| 2.25%   | 225            |

Rumus konversi:

```
bps = percent × 100
```

---

# 4️⃣ Cara Hitung Percentage Fee

Semua perhitungan menggunakan integer math.

Formula:

```
fee = amount * bps / 10000
```

Contoh:

```
amount = 10000
gateway_fee_percentage = 70
```

Perhitungan:

```
10000 * 70 / 10000 = 70
```

Hasil 70 rupiah.

Tidak ada float.
Tidak ada pembulatan manual.
Division otomatis floor (deterministic).

---

# 5️⃣ Cara Hitung Total Fee

Jika ada fixed + percentage:

```
percent_part = amount * bps / 10000
total_fee = percent_part + fixed
```

Contoh:

```
amount = 10000
platform_fee_percentage = 200
platform_fee_fixed = 750
```

Perhitungan:

```
percent_part = 10000 * 200 / 10000 = 200
total_fee = 200 + 750 = 950
```

---

# 6️⃣ Urutan Perhitungan Transaksi

Standar perhitungan:

1. Hitung gateway fee
2. Kurangi dari amount
3. Hitung platform fee dari sisa
4. Kurangi lagi untuk mendapatkan net amount

Contoh lengkap:

```
amount = 10000
gateway = 0.7% (70)
platform = 750 + 2% (200)
```

Step:

```
pg_fee = 10000 * 70 / 10000 = 70
after_pg = 10000 - 70 = 9930

platform_percent = 9930 * 200 / 10000 = 198
platform_fee = 198 + 750 = 948

net = 9930 - 948 = 8982
```

Validasi:

```
fee_amount + net_amount = amount
```

---

# 7️⃣ Larangan Keras

❌ Tidak boleh menggunakan float untuk uang
❌ Tidak boleh menyimpan persen sebagai 0.7 atau 0.02
❌ Tidak boleh menyimpan hasil persen sebelum dihitung ke rupiah
❌ Tidak boleh ada selisih 1 rupiah tanpa penjelasan

---

# 8️⃣ Helper Function (Go Standard)

```go
const PercentDivisor int64 = 10000

func CalculateFee(amount int64, fixed int64, bps int64) int64 {
    percentPart := amount * bps / PercentDivisor
    return percentPart + fixed
}
```

---

# 9️⃣ Future Proofing

Jika suatu saat support USD:

* Tetap BIGINT
* USD disimpan dalam cent
* Minor unit USD = 100
* Formula tetap sama
* Tidak perlu ubah arsitektur

---

# 🔒 Prinsip Utama

* Semua uang = integer
* Semua persen = basis point
* Semua kalkulasi = deterministic
* Ledger harus selalu balance

---

Kalau kamu mau, aku bisa bantu buatkan versi lebih formal untuk internal RFC atau ADR (Architecture Decision Record) supaya ini jadi keputusan arsitektur resmi.
