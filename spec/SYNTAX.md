# MOJA Syntax Specification — v0.4

**Status**: Final Draft — v0.4  
**Tanggal**: September 2026  
**Bahasa Implementasi Interpreter**: Go 1.26+

---

## 1. Pengantar

MojoLang adalah bahasa pemrograman edukatif berbasis kosakata Jawa Mojokerto. Spesifikasi ini mendefinisikan sintaks, semantik, dan perilaku interpreter versi 0.1 (MVP).

---

## 2. Konvensi Notasi Grammar

| Notasi | Arti |
|---|---|
| `A → B C` | A terdiri dari B diikuti C |
| `A → B \| C` | A adalah B atau C |
| `A → B*` | Nol atau lebih B |
| `A → B+` | Satu atau lebih B |
| `A → B?` | Nol atau satu B |
| `"literal"` | Token literal persis seperti ditulis |
| `UPPERCASE` | Token terminal dari lexer |

---

## 3. Karakter, Encoding, dan Terminasi Pernyataan

- **Encoding**: UTF-8.
- **Baris**: Dipisahkan oleh `\n` (LF) atau `\r\n` (CRLF).
- **Spasi & Tab**: Diabaikan di antara token.

### 3.1 Terminasi Pernyataan (KEPUTUSAN: AMB-01)

**Newline (`\n`) berfungsi sebagai terminator pernyataan** (seperti Go).

- Setiap pernyataan harus berada pada baris terpisah.
- Dua pernyataan di satu baris adalah **syntax error**.
- Pengecualian: Baris yang diakhiri dengan `{` boleh dilanjutkan ke baris berikutnya (baris dalam blok).

```
// Valid
x = 10
y = 20

// Syntax Error: dua statement di satu baris
x = 10 y = 20
```

Newline **di dalam** ekspresi yang belum selesai diabaikan:
```
// Valid — ekspresi panjang yang terbungkus tanda kurung
hasil = (10 +
         20)
```

---

## 4. Komentar

```
komentar → "//" karakter_apapun_sampai_\n
```

Komentar satu baris. Diabaikan oleh lexer. Komentar multi-baris tidak didukung di v0.1.

---

## 5. Literal dan Tipe Data

### 5.1 Number

```
NUMBER → DIGIT+ ("." DIGIT+)?
DIGIT  → "0".."9"
```

- Internal: `float64`.
- Bilangan negatif dibentuk oleh unary `-`.
- **Valid**: `0`, `42`, `3.14`
- **Tidak valid**: `.5`, `1.`

**Format cetak Number:**
- Jika hasilnya bilangan bulat (tidak ada fraksi), cetak tanpa titik desimal: `5` bukan `5.0`
- Jika hasilnya desimal, cetak dengan presisi alami: `3.3333333333`
- Implementasi: `strconv.FormatFloat(val, 'f', -1, 64)`

### 5.2 String

```
STRING     → '"' karakter* '"'
karakter   → <karakter apapun kecuali '"' dan '\'>
           | escape_sequence
escape_seq → '\"' | '\\' | '\n' | '\t'
```

- Satu baris saja (no multiline string literals).
- Escape sequence: `\"`, `\\`, `\n`, `\t`.

### 5.3 Boolean

```
BOOLEAN → "bener" | "salah"
```

- `bener` = true, `salah` = false.
- **Keyword terpesan**, bukan identifier.
- Case-sensitive: `Bener`, `BENER` adalah identifier biasa, bukan Boolean.

### 5.4 Array / List (MOJA v0.3)

```
ARRAY_LITERAL → "[" (expression ("," expression)*)? "]"
```

- Array adalah koleksi terurut dari nol atau lebih elemen.
- Elemen array dapat berupa tipe data apapun (Number, String, Boolean, atau Array bersarang).
- **Array Kosong**: `[]`
- **Array Bersatu**: `[10, 20, 30]`
- **Array Campuran**: `[10, "halo", bener]`
- **Array Bersarang**: `[[1, 2], [3, 4]]`

**Format cetak Array (`omong`):**
- Diinspeksi dalam tanda kurung siku dengan elemen dipisahkan koma dan spasi: `[10, 20, 30]`
- String di dalam array dicetak dengan tanda petik ganda: `["Budi", "Sari"]`

### 5.5 Akses & Mutasi Indeks (MOJA v0.3)

```
INDEX_EXPR   → primary "[" expression "]"
INDEX_ASSIGN → INDEX_EXPR "=" expression
```

- Indeks berbasis nol (0-based indexing).
- Indeks harus berupa tipe **Number** dengan nilai **bilangan bulat (integer)** non-negatif.
- Akses di luar rentang (`idx < 0` atau `idx >= len`): menghasilkan **Runtime Error**.
- Mutasi elemen: `angka[1] = 99` memperbarui nilai elemen pada indeks tertentu.

### 5.6 Fungsi Bawaan `panjang()` (MOJA v0.3)

```
panjang(array | string) → Number
```

- Mengembalikan jumlah elemen dalam Array atau jumlah karakter dalam String.
- `panjang([])` → `0`
- `panjang([10, 20, 30])` → `3`
- `panjang("MOJA")` → `4`

---

## 6. Identifier

```
IDENTIFIER → huruf (huruf | digit | "_")*
huruf      → 'a'..'z' | 'A'..'Z' | '_'
digit      → '0'..'9'
```

**Keyword terpesan** (tidak boleh dipakai sebagai identifier):

| Keyword | Fungsi |
|---|---|
| `omong` | Output |
| `yen` | Kondisi (if) |
| `liyane` | Alternatif (else) |
| `baleni` | Perulangan (loop) |
| `saka` | Awal perulangan (from) |
| `tekan` | Akhir perulangan (to) |
| `gawe` | Deklarasi fungsi (function) |
| `balik` | Pengembalian nilai (return) |
| `bener` | Boolean true |
| `salah` | Boolean false |
| `lan` | Operator logika AND |
| `utawa` | Operator logika OR |
| `ora` | Operator logika NOT (unary) |

---

## 7. Operator

### 7.1 Operator Aritmatika

| Operator | Operasi | Aturan Tipe |
|---|---|---|
| `+` | Penjumlahan / Penggabungan | Number+Number → Number; jika salah satu String → String concat |
| `-` | Pengurangan | Number only |
| `*` | Perkalian | Number only |
| `/` | Pembagian | Number only |

**Aturan String Coercion untuk `+`:**
- Jika salah satu operand String, operand lain dikonversi ke teks.
- Number → tanpa trailing zero: `20` → `"20"`, `3.14` → `"3.14"`
- Boolean → `"bener"` atau `"salah"`
- `"-"`, `"*"`, `"/"` dengan tipe bukan Number → runtime error.

### 7.2 Operator Perbandingan

| Operator | Arti | Tipe Hasil |
|---|---|---|
| `==` | Sama dengan | Boolean |
| `!=` | Tidak sama dengan | Boolean |
| `>` | Lebih besar | Boolean (Number only) |
| `<` | Lebih kecil | Boolean (Number only) |
| `>=` | Lebih besar/sama | Boolean (Number only) |
| `<=` | Lebih kecil/sama | Boolean (Number only) |

- `==` dan `!=` berlaku untuk semua tipe. Perbandingan beda tipe selalu `salah`.
- `>`, `<`, `>=`, `<=` hanya untuk Number → runtime error jika tipe lain.

### 7.3 Operator Logika (MOJA v0.4)

| Operator | Jenis | Arti | Contoh |
|---|---|---|---|
| `lan` | Infix | AND logika | `bener lan salah` → `salah` |
| `utawa` | Infix | OR logika | `salah utawa bener` → `bener` |
| `ora` | Prefix unary | NOT logika | `ora bener` → `salah` |

**Aturan tipe:**
- Semua operand **harus bertipe Boolean** (`bener` / `salah`).
- Non-boolean menghasilkan **runtime error**:
  - `10 lan bener` → `Operand kiri operator 'lan' harus bertipe Boolean`
  - `ora 42` → `Operator 'ora' hanya mendukung tipe Boolean`
- Hasil selalu **Boolean**.

**Short-circuit evaluation:**
- `salah lan ekspresi` → mengembalikan `salah` **tanpa** mengevaluasi `ekspresi`.
- `bener utawa ekspresi` → mengembalikan `bener` **tanpa** mengevaluasi `ekspresi`.

```
// Contoh
a = bener
b = salah

omong(a lan b)    // salah
omong(a utawa b)  // bener
omong(ora a)      // salah

// Kombinasi dengan yen
umur = 20
duwe_ktp = bener
yen umur >= 17 lan duwe_ktp {
    omong("Syarat terpenuhi")
}
```

### 7.4 Precedence (MOJA v0.4)

| Level | Operator | Asosiatif |
|---|---|---|
| 8 (tertinggi) | `[...]` akses indeks | Kiri |
| 7 | pemanggilan fungsi `f(...)` | Kiri |
| 6 | Unary `-`, `ora` | Kanan |
| 5 | `*`, `/` | Kiri |
| 4 | `+`, `-` | Kiri |
| 3 | `==`, `!=`, `>`, `<`, `>=`, `<=` | Kiri |
| 2 | `lan` (AND) | Kiri |
| 1 (terendah) | `utawa` (OR) | Kiri |

**Contoh precedence:**
```
ora bener lan salah        // (ora bener) lan salah  → salah
bener utawa salah lan salah // bener utawa (salah lan salah) → bener
(bener utawa salah) lan salah // false (kurung mengubah urutan)
```

---

## 8. Grammar Lengkap

```
program        → statement* EOF

statement      → omong_stmt
               | assign_stmt
               | index_assign_stmt
               | yen_stmt
               | baleni_stmt
               | gawe_stmt
               | return_stmt
               | expr_stmt
               | NEWLINE

omong_stmt     → "omong" "(" expression ")" NEWLINE

assign_stmt    → IDENTIFIER "=" expression NEWLINE

index_assign_stmt → INDEX_EXPR "=" expression NEWLINE

yen_stmt       → "yen" expression block ("liyane" block)?

baleni_stmt    → "baleni" IDENTIFIER "saka" expression "tekan" expression block

gawe_stmt      → "gawe" IDENTIFIER "(" params? ")" block

return_stmt    → "balik" expression NEWLINE

expr_stmt      → expression NEWLINE

block          → "{" NEWLINE statement* "}"

params         → IDENTIFIER ("," IDENTIFIER)*

expression     → logical_or
logical_or     → logical_and ("utawa" logical_and)*
logical_and    → comparison ("lan" comparison)*
comparison     → addition (("==" | "!=" | ">" | "<" | ">=" | "<=") addition)*
addition       → multiplication (("+" | "-") multiplication)*
multiplication → unary (("*" | "/") unary)*
unary          → "-" unary | "ora" unary | primary
primary        → NUMBER | STRING | "bener" | "salah"
               | IDENTIFIER | call_expr | index_expr | "(" expression ")"

call_expr      → IDENTIFIER "(" arguments? ")"
index_expr     → primary "[" expression "]"
arguments      → expression ("," expression)*
```

---

## 9. Semantik Pernyataan

### 9.1 `omong`

Mencetak representasi teks ekspresi ke stdout, diikuti `\n`. Menerima tepat satu argumen. Bukan fungsi — tidak bisa disimpan ke variabel.

### 9.2 Assignment (KEPUTUSAN: AMB-03)

**Semantik "walk up the scope chain":**

1. Saat assignment `x = nilai`, interpreter mencari nama `x` dari scope saat ini ke atas menuju global scope.
2. Jika ditemukan di scope manapun, **variabel tersebut diperbarui**.
3. Jika tidak ditemukan di scope manapun, variabel baru dibuat di **scope saat ini**.

```
// Pola accumulator bekerja dengan benar
gawe hitung() {
    total = 0               // buat 'total' di scope fungsi
    baleni i saka 1 tekan 3 {
        total = total + i   // 'total' ditemukan di scope fungsi → diperbarui
    }
    balik total             // mengembalikan 6
}
```

```
// Variabel global dapat dimodifikasi dari dalam fungsi
x = 10
gawe ubah() {
    x = 99    // 'x' ditemukan di global scope → diperbarui
}
ubah()
omong(x)      // Mencetak: 99
```

> **Catatan**: Ini berbeda dari pure block-scoped semantics. Dipilih karena lebih intuitif untuk target pengguna pemula MojoLang.

### 9.3 Kondisi `yen` / `liyane`

- Ekspresi kondisi **harus** mengevaluasi ke Boolean. Jika tidak → runtime error.
- `liyane` opsional.
- Tidak ada `yen liyane` (else-if) di v0.1.

```
yen umur >= 17 {
    omong("Wis cukup umur")
} liyane {
    omong("Durung cukup umur")
}
```

### 9.4 Perulangan `baleni` (KEPUTUSAN: AMB-07, AMB-08)

```
baleni IDENTIFIER saka expr_start tekan expr_end block
```

- `expr_start` dan `expr_end` dievaluasi **sekali** sebelum loop.
- Keduanya harus Number → runtime error jika tidak.
- Batas **inklusif**: `saka 1 tekan 5` → iterasi untuk `1, 2, 3, 4, 5`.
- Step selalu `+1`.
- Jika `saka > tekan`: **zero iterations**, tidak error (AMB-07 resolved).
- Variabel iterator tersedia hanya di dalam blok loop (dibuat di scope blok).
- **Modifikasi variabel iterator dari dalam loop = runtime error** (AMB-08 resolved):
  ```
  [MojoLang Error] Baris N: Variabel iterator 'i' tidak dapat diubah di dalam loop
  ```

### 9.5 Fungsi `gawe` (KEPUTUSAN: AMB-04)

- **Tidak ada hoisting**: Fungsi hanya bisa dipanggil **setelah** baris deklarasinya.
- Pemanggilan fungsi sebelum deklarasi → runtime error saat pemanggilan.
- Rekursi didukung (fungsi dapat memanggil dirinya sendiri karena fungsi sudah ada di scope saat blok dievaluasi).
- Jumlah argumen harus sama persis dengan parameter → runtime error jika tidak cocok.
- Batas kedalaman rekursi: **1000 level**. Jika terlampaui:
  ```
  [MojoLang Error] Baris N: Kedalaman rekursi melampaui batas (1000)
  ```

### 9.6 Return `balik` (KEPUTUSAN: AMB-02)

- `balik` **harus** diikuti ekspresi di baris yang sama (AMB-02 resolved).
- `balik` sendirian di satu baris = **syntax error**:
  ```
  [MojoLang Error] Baris N: 'balik' harus diikuti ekspresi
  ```
- `balik` di luar fungsi = runtime error:
  ```
  [MojoLang Error] Baris N: 'balik' hanya dapat digunakan di dalam fungsi
  ```
- **Keyword return adalah `balik`, bukan `bali`.**
- Fungsi yang berakhir tanpa `balik` dan hasilnya digunakan dalam ekspresi → runtime error.

---

## 10. Scope dan Environment

**Model**: Lexical scoping dengan linked environment chain.

1. Program dimulai dengan satu global environment.
2. Setiap blok `{ }` membuat environment baru dengan parent = environment saat ini.
3. **Lookup variabel** (baca): Cari dari scope saat ini → naik ke parent → sampai global.
4. **Assignment variabel** (tulis): Cari nama dari scope saat ini → naik → jika ditemukan, perbarui di sana. Jika tidak ditemukan, buat baru di scope saat ini.
5. Setelah blok selesai, environment blok dibuang.
6. Fungsi membawa closure environment (environment di mana fungsi didefinisikan sebagai parent).

---

## 11. Error Handling

### Format Pesan Error

```
[MojoLang Error] Baris N: <deskripsi dalam Bahasa Indonesia>
```

### Daftar Error Utama

| Kategori | Kondisi | Pesan |
|---|---|---|
| Syntax | Token tidak dikenali | `Token tidak dikenali: '<karakter>'` |
| Syntax | `balik` tanpa ekspresi | `'balik' harus diikuti ekspresi` |
| Syntax | String tidak ditutup | `String tidak ditutup pada baris N` |
| Runtime | Variabel belum ada | `Variabel '<nama>' belum dideklarasikan` |
| Runtime | Tipe salah untuk `-`/`*`/`/` | `Operator '<op>' hanya mendukung tipe Number` |
| Runtime | Tipe salah untuk `>`/`<`/`>=`/`<=` | `Operator '<op>' hanya mendukung tipe Number` |
| Runtime | Kondisi bukan Boolean | `Kondisi 'yen' harus bernilai Boolean (bener/salah)` |
| Runtime | `balik` di luar fungsi | `'balik' hanya dapat digunakan di dalam fungsi` |
| Runtime | Argumen tidak cocok | `Fungsi '<nama>' butuh N argumen, diberikan M` |
| Runtime | Pembagian dengan nol | `Pembagian dengan nol tidak diperbolehkan` |
| Runtime | Rekursi terlalu dalam | `Kedalaman rekursi melampaui batas (1000)` |
| Runtime | Ubah iterator loop | `Variabel iterator '<nama>' tidak dapat diubah di dalam loop` |
| Runtime | Fungsi tidak ada | `Fungsi '<nama>' belum dideklarasikan` |
| Runtime | Operand `lan`/`utawa` non-Boolean | `Operand kiri/kanan operator 'lan'/'utawa' harus bertipe Boolean` |
| Runtime | Operand `ora` non-Boolean | `Operator 'ora' hanya mendukung tipe Boolean (bener/salah)` |

### Perilaku saat Error

- Hentikan eksekusi pada error pertama.
- Cetak pesan error ke **stderr**.
- Exit code `1`.
- Exit code `0` jika program selesai normal.

---

## 12. Contoh Program Lengkap

```
// halo.mojo
omong("Sugeng rawuh ing Mojokerto!")
```

```
// variabel.mojo
jeneng = "Dimaz"
umur = 20
omong("Jeneng: " + jeneng)
omong("Umur: " + umur)
```

```
// kondisi.mojo
umur = 16
yen umur >= 17 {
    omong("Wis cukup umur")
} liyane {
    omong("Durung cukup umur")
}
```

```
// perulangan.mojo
baleni i saka 1 tekan 5 {
    omong("Iterasi ke: " + i)
}
```

```
// fungsi.mojo
gawe tambah(a, b) {
    balik a + b
}

hasil = tambah(10, 20)
omong("Hasil tambah: " + hasil)
```

```
// fungsi_rekursif.mojo
gawe faktorial(n) {
    yen n <= 1 {
        balik 1
    } liyane {
        balik n * faktorial(n - 1)
    }
}

omong("5! = " + faktorial(5))
```

```
// accumulator.mojo — membuktikan walk-up scope assignment
gawe jumlahkan(maks) {
    total = 0
    baleni i saka 1 tekan maks {
        total = total + i
    }
    balik total
}

omong("Jumlah 1..10 = " + jumlahkan(10))
```

---

## 13. Di Luar Scope v0.1

- REPL interaktif
- Error handling (`coba`/`tangkap`)
- Impor antar-modul
- `yen liyane` (else-if chain)
- `baleni` dengan langkah kustom dan kondisi berhenti
- Fungsi built-in tambahan (slice, hapus, dsb.)
- Map / Dictionary
- Compiler ke bytecode

---

## Riwayat Perubahan

| Versi | Tanggal | Keterangan |
|---|---|---|
| 0.1-draft | September 2026 | Draft awal |
| 0.1-final | September 2026 | AMB-01..08 diselesaikan, grammar dilengkapi |
| 0.3 | September 2026 | Array/List: literal, akses indeks, mutasi, builtin `panjang()` |
| 0.4 | September 2026 | Operator logika: `lan` (AND), `utawa` (OR), `ora` (NOT) dengan short-circuit |
