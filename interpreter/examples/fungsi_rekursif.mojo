// fungsi_rekursif.mojo
gawe faktorial(n) {
    yen n <= 1 {
        balik 1
    } liyane {
        balik n * faktorial(n - 1)
    }
}

omong("5! = " + faktorial(5))
