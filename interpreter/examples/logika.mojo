// Contoh operator logika MOJA v0.4
// lan = AND, utawa = OR, ora = NOT

a = bener
b = salah

// Operator dasar
omong("=== Operator Dasar ===")
omong(a lan b)
omong(a utawa b)
omong(ora a)
omong(ora b)

// Tabel kebenaran AND (lan)
omong("=== Tabel Kebenaran AND (lan) ===")
omong(bener lan bener)
omong(bener lan salah)
omong(salah lan bener)
omong(salah lan salah)

// Tabel kebenaran OR (utawa)
omong("=== Tabel Kebenaran OR (utawa) ===")
omong(bener utawa bener)
omong(bener utawa salah)
omong(salah utawa bener)
omong(salah utawa salah)

// Precedence: lan lebih ketat dari utawa
omong("=== Precedence ===")
omong(bener utawa salah lan salah)
omong(ora bener lan salah)
omong((bener utawa salah) lan salah)

// Kombinasi dengan perbandingan dan yen
omong("=== Kombinasi dengan Yen ===")
umur = 20
duwe_ktp = bener
yen umur >= 17 lan duwe_ktp {
    omong("Syarat terpenuhi")
}

panas = bener
udan = salah
yen ora udan lan panas {
    omong("Bisa piknik")
}

// Short-circuit
omong("=== Short-circuit ===")
omong(salah lan bener)
omong(bener utawa salah)
