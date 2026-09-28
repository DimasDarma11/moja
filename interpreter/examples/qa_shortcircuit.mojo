// QA: Short-circuit dengan side effect via fungsi
panggil = 0
gawe naikkan() {
    panggil = panggil + 1
    balik bener
}
omong(salah lan naikkan())
omong(panggil)
omong(bener utawa naikkan())
omong(panggil)