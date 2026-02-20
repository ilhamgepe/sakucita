Pertanyaan kamu bagus. Dan ini penting supaya nanti tidak refactor besar.

Jawaban singkatnya dulu:

> ❗ Divisor 10.000 **tidak berubah** walaupun multi currency.
> ❗ Yang berubah hanya **minor unit currency**, bukan divisor persen.

Sekarang kita bedah pelan-pelan.

---

# 🔥 1️⃣ Yang Tidak Berubah: Divisor 10.000

Divisor 10.000 itu untuk **persentase (basis point)**.

Itu tidak ada hubungannya dengan IDR atau USD.

Karena:

```
1% = 100 bps
100% = 10.000 bps
```

Mau IDR, USD, JPY, tetap sama.

Jadi:

```
fee = amount * bps / 10000
```

Tetap berlaku untuk semua currency.

---

# 🧱 2️⃣ Yang Berubah: Minor Unit Currency

Yang berbeda antar currency adalah:

| Currency | Minor Unit |
| -------- | ---------- |
| IDR      | 1          |
| USD      | 100        |
| EUR      | 100        |
| JPY      | 1          |
| KWD      | 1000       |

Contoh:

* USD 12.34 → simpan 1234
* IDR 12.34 → tidak valid (IDR tidak punya cent)
* JPY 500 → simpan 500

Jadi yang berubah adalah cara kamu **menyimpan amount**, bukan cara hitung persen.

---

# 🧱 3️⃣ Apa yang Harus Ditambahkan di Table?

Kalau mau multi currency proper, tambahkan:

## ✅ A. Tambah currency column

Di transactions:

```sql
currency CHAR(3) NOT NULL
```

Contoh:

```
IDR
USD
JPY
```

---

## ✅ B. Buat table currencies

Ini best practice.

```sql
CREATE TABLE currencies (
  code CHAR(3) PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  minor_unit SMALLINT NOT NULL
);
```

Contoh isi:

```sql
INSERT INTO currencies (code, name, minor_unit) VALUES
('IDR', 'Indonesian Rupiah', 1),
('USD', 'US Dollar', 100),
('JPY', 'Japanese Yen', 1);
```

---

# 🧠 4️⃣ Cara Hitung Kalau Multi Currency

Misalnya USD:

User bayar:

```
$12.34
```

Simpan:

```
amount = 1234
currency = 'USD'
```

Fee 2%:

```
1234 * 200 / 10000 = 24
```

Artinya 24 cent ($0.24)

Tidak perlu ubah divisor.
Tidak perlu logic tambahan.

---

# 🎯 Jadi Apa Yang Perlu Kamu Tambahkan?

Minimal:

### Di transactions

```
currency CHAR(3)
```

### Buat table currencies

```
code
minor_unit
```

Dan mungkin:

### Di payment_channels

Kalau fee beda per currency:

```
currency CHAR(3)
```

---

# ⚠️ Yang Tidak Perlu Diubah

❌ Tidak perlu ubah divisor
❌ Tidak perlu ubah basis point
❌ Tidak perlu ubah logic hitung fee

---

# 🧠 Mental Model yang Benar

* Minor unit → soal penyimpanan amount
* Basis point → soal penyimpanan persen
* Dua sistem ini tidak saling mempengaruhi

---

# 🔥 Arsitektur Fintech yang Sehat

Semua uang → simpan dalam minor unit (BIGINT)
Semua persen → simpan dalam bps (BIGINT)
Semua hitung → integer math

Multi currency tetap stabil tanpa refactor besar.

---

Kalau kamu mau, aku bisa gambarkan flow lengkap multi currency dari request API sampai masuk ledger supaya makin kebayang production-grade nya.
