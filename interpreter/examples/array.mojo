// array.mojo — Contoh penggunaan Array ing MOJA v0.3
angka = [10, 20, 30, 40]

omong("Elemen pisanan (indeks 0):")
omong(angka[0])

omong("Jumlah elemen (panjang):")
omong(panjang(angka))

gawe jumlah(a, b) {
    balik a + b
}

omong("Jumlah elemen 0 + elemen 1:")
omong(jumlah(angka[0], angka[1]))

// Mutasi elemen array
angka[1] = 99
omong("Setelah mutasi indeks 1:")
omong(angka)

// Nested array (matriks)
matriks = [[1, 2], [3, 4]]
omong("Matriks indeks [0][1]:")
omong(matriks[0][1])
