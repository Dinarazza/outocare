package main

import (
	"autocare/repository"
	"autocare/service"
	"fmt"
	"strconv"
	"strings"
)

// HandlerSequentialSearch menangani pencarian menggunakan Sequential Search
func HandlerSequentialSearch() {
	fmt.Println()
	SubHeader("SEQUENTIAL SEARCH")
	fmt.Println("  Metode: memeriksa setiap data satu per satu dari awal.")
	fmt.Println("  Kompleksitas: O(n)")
	fmt.Println()

	list, err := repository.GetSemuaKendaraan()
	if err != nil || len(list) == 0 {
		InfoTidakAda("Belum ada data kendaraan.")
		TungguEnter()
		return
	}

	platNomor := bacaInput("  Masukkan Plat Nomor yang dicari: ")
	if platNomor == "" {
		InfoGagal("Plat nomor tidak boleh kosong.")
		TungguEnter()
		return
	}

	fmt.Printf("\n  Mencari '%s' menggunakan Sequential Search...\n", strings.ToUpper(platNomor))
	hasil := service.SequentialSearch(list, platNomor)

	if hasil == nil {
		InfoTidakAda(fmt.Sprintf("Kendaraan dengan plat nomor '%s' tidak ditemukan.", strings.ToUpper(platNomor)))
	} else {
		fmt.Println()
		InfoSukses("Kendaraan ditemukan!")
		GarisTipis(60)
		fmt.Printf("  ID            : %d\n", hasil.ID)
		fmt.Printf("  Plat Nomor    : %s\n", hasil.PlatNomor)
		fmt.Printf("  Merek         : %s\n", hasil.Merek)
		fmt.Printf("  Model         : %s\n", hasil.Model)
		fmt.Printf("  Tahun Produksi: %d\n", hasil.TahunProduksi)
		fmt.Printf("  Tipe          : %s\n", hasil.TipeKendaraan)
		servis := repository.GetServisTermakhir(hasil.ID)
		fmt.Printf("  Servis Terakhir: %s\n", servis)
	}
	TungguEnter()
}

// HandlerBinarySearch menangani pencarian menggunakan Binary Search
func HandlerBinarySearch() {
	fmt.Println()
	SubHeader("BINARY SEARCH")
	fmt.Println("  Metode: membagi data menjadi dua bagian dan membandingkan tengah.")
	fmt.Println("  Kompleksitas: O(log n) — SYARAT: data sudah terurut!")
	fmt.Println()

	list, err := repository.GetSemuaKendaraan()
	if err != nil || len(list) == 0 {
		InfoTidakAda("Belum ada data kendaraan.")
		TungguEnter()
		return
	}

	// Binary Search membutuhkan data terurut
	sorted := service.UrutkanByPlatNomor(list)
	fmt.Printf("  Data diurutkan terlebih dahulu (%d kendaraan).\n", len(sorted))

	platNomor := bacaInput("  Masukkan Plat Nomor yang dicari: ")
	if platNomor == "" {
		InfoGagal("Plat nomor tidak boleh kosong.")
		TungguEnter()
		return
	}

	fmt.Printf("\n  Mencari '%s' menggunakan Binary Search...\n", strings.ToUpper(platNomor))
	hasil := service.BinarySearch(sorted, platNomor)

	if hasil == nil {
		InfoTidakAda(fmt.Sprintf("Kendaraan dengan plat nomor '%s' tidak ditemukan.", strings.ToUpper(platNomor)))
	} else {
		fmt.Println()
		InfoSukses("Kendaraan ditemukan!")
		GarisTipis(60)
		fmt.Printf("  ID            : %d\n", hasil.ID)
		fmt.Printf("  Plat Nomor    : %s\n", hasil.PlatNomor)
		fmt.Printf("  Merek         : %s\n", hasil.Merek)
		fmt.Printf("  Model         : %s\n", hasil.Model)
		fmt.Printf("  Tahun Produksi: %d\n", hasil.TahunProduksi)
		fmt.Printf("  Tipe          : %s\n", hasil.TipeKendaraan)
		servis := repository.GetServisTermakhir(hasil.ID)
		fmt.Printf("  Servis Terakhir: %s\n", servis)
	}
	TungguEnter()
}

// HandlerSelectionSortTahun menangani pengurutan berdasarkan tahun produksi
func HandlerSelectionSortTahun() {
	fmt.Println()
	SubHeader("SELECTION SORT - BERDASARKAN TAHUN PRODUKSI")
	fmt.Println("  Algoritma: pilih elemen terkecil/terbesar, tukar dengan posisi saat ini.")
	fmt.Println("  Kompleksitas: O(n²)")
	fmt.Println()

	list, err := repository.GetSemuaKendaraan()
	if err != nil || len(list) == 0 {
		InfoTidakAda("Belum ada data kendaraan.")
		TungguEnter()
		return
	}

	MenuUrutan()
	pilihanStr := bacaInput("")
	pilihan, _ := strconv.Atoi(pilihanStr)
	ascending := pilihan != 2

	arahTeks := "Ascending (Terlama → Terbaru)"
	if !ascending {
		arahTeks = "Descending (Terbaru → Terlama)"
	}
	fmt.Printf("\n  Mengurutkan %d kendaraan secara %s...\n", len(list), arahTeks)

	sorted := service.SelectionSortByTahun(list, ascending)

	fmt.Println()
	fmt.Printf("  Hasil Pengurutan - Tahun Produksi (%s):\n", arahTeks)
	TampilkanKendaraan(sorted)
	TungguEnter()
}

// HandlerInsertionSortServis menangani pengurutan berdasarkan tanggal servis terakhir
func HandlerInsertionSortServis() {
	fmt.Println()
	SubHeader("INSERTION SORT - BERDASARKAN TANGGAL SERVIS TERAKHIR")
	fmt.Println("  Algoritma: sisipkan elemen ke posisi yang tepat dalam bagian terurut.")
	fmt.Println("  Kompleksitas: O(n²) worst case, O(n) best case")
	fmt.Println()

	list, err := repository.GetSemuaKendaraan()
	if err != nil || len(list) == 0 {
		InfoTidakAda("Belum ada data kendaraan.")
		TungguEnter()
		return
	}

	// Buat map tanggal servis terakhir
	servisMap := make(map[int]string)
	for _, k := range list {
		servisMap[k.ID] = repository.GetServisTermakhir(k.ID)
	}

	MenuUrutan()
	pilihanStr := bacaInput("")
	pilihan, _ := strconv.Atoi(pilihanStr)
	ascending := pilihan != 2

	arahTeks := "Ascending (Terlama → Terbaru)"
	if !ascending {
		arahTeks = "Descending (Terbaru → Terlama)"
	}

	fmt.Printf("\n  Mengurutkan %d kendaraan secara %s...\n", len(list), arahTeks)
	sorted := service.InsertionSortByTanggalServis(list, servisMap, ascending)

	fmt.Println()
	fmt.Printf("  Hasil Pengurutan - Tanggal Servis Terakhir (%s):\n", arahTeks)
	TampilkanKendaraan(sorted)
	TungguEnter()
}
