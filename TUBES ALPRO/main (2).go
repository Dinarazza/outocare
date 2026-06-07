package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ===================== STRUCT DEFINITIONS =====================

type Pemilik struct {
	ID      string
	Nama    string
	Telepon string
	Alamat  string
}

type Kendaraan struct {
	ID           string
	PlatNomor    string
	Merek        string
	Model        string
	TahunProduksi int
	PemilikID    string
	TanggalServisTerakhir string
}

type RiwayatServis struct {
	ID            string
	KendaraanID   string
	TanggalServis string
	JenisKerusakan string
	Deskripsi     string
	Biaya         float64
	Teknisi       string
}

// ===================== DATA STORAGE =====================

var (
	daftarPemilik   []Pemilik
	daftarKendaraan []Kendaraan
	daftarServis    []RiwayatServis
	reader          = bufio.NewReader(os.Stdin)
)

// ===================== HELPER FUNCTIONS =====================

func input(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func pause() {
	input("\nTekan Enter untuk melanjutkan...")
}

func generateID(prefix string, count int) string {
	return fmt.Sprintf("%s%03d", prefix, count+1)
}

func printSeparator() {
	fmt.Println(strings.Repeat("=", 60))
}

func printLine() {
	fmt.Println(strings.Repeat("-", 60))
}

func cariPemilikByID(id string) *Pemilik {
	for i := range daftarPemilik {
		if daftarPemilik[i].ID == id {
			return &daftarPemilik[i]
		}
	}
	return nil
}

func cariKendaraanByID(id string) *Kendaraan {
	for i := range daftarKendaraan {
		if daftarKendaraan[i].ID == id {
			return &daftarKendaraan[i]
		}
	}
	return nil
}

// ===================== ALGORITMA PENCARIAN =====================

// Sequential Search berdasarkan plat nomor
func sequentialSearch(platNomor string) *Kendaraan {
	for i := range daftarKendaraan {
		if strings.EqualFold(daftarKendaraan[i].PlatNomor, platNomor) {
			return &daftarKendaraan[i]
		}
	}
	return nil
}

// Binary Search berdasarkan plat nomor (data harus terurut)
func binarySearch(platNomor string) *Kendaraan {
	// Buat salinan terurut
	sorted := make([]Kendaraan, len(daftarKendaraan))
	copy(sorted, daftarKendaraan)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToUpper(sorted[i].PlatNomor) < strings.ToUpper(sorted[j].PlatNomor)
	})

	low, high := 0, len(sorted)-1
	target := strings.ToUpper(platNomor)

	for low <= high {
		mid := (low + high) / 2
		current := strings.ToUpper(sorted[mid].PlatNomor)
		if current == target {
			// Cari di data asli
			return sequentialSearch(platNomor)
		} else if current < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return nil
}

// ===================== ALGORITMA PENGURUTAN =====================

// Selection Sort berdasarkan tahun produksi (ascending)
func selectionSortByTahun(data []Kendaraan) []Kendaraan {
	result := make([]Kendaraan, len(data))
	copy(result, data)
	n := len(result)

	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if result[j].TahunProduksi < result[minIdx].TahunProduksi {
				minIdx = j
			}
		}
		result[i], result[minIdx] = result[minIdx], result[i]
	}
	return result
}

// Insertion Sort berdasarkan tanggal servis terakhir (ascending)
func insertionSortByTanggalServis(data []Kendaraan) []Kendaraan {
	result := make([]Kendaraan, len(data))
	copy(result, data)
	n := len(result)

	for i := 1; i < n; i++ {
		key := result[i]
		j := i - 1
		for j >= 0 && result[j].TanggalServisTerakhir > key.TanggalServisTerakhir {
			result[j+1] = result[j]
			j--
		}
		result[j+1] = key
	}
	return result
}

// ===================== MENU PEMILIK =====================

func menuPemilik() {
	for {
		clearScreen()
		printSeparator()
		fmt.Println("         MANAJEMEN DATA PEMILIK")
		printSeparator()
		fmt.Println("  1. Tampilkan Semua Pemilik")
		fmt.Println("  2. Tambah Pemilik")
		fmt.Println("  3. Ubah Pemilik")
		fmt.Println("  4. Hapus Pemilik")
		fmt.Println("  0. Kembali")
		printLine()
		pilihan := input("Pilih menu: ")

		switch pilihan {
		case "1":
			tampilkanSemuaPemilik()
		case "2":
			tambahPemilik()
		case "3":
			ubahPemilik()
		case "4":
			hapusPemilik()
		case "0":
			return
		default:
			fmt.Println("Pilihan tidak valid!")
			pause()
		}
	}
}

func tampilkanSemuaPemilik() {
	clearScreen()
	printSeparator()
	fmt.Println("           DAFTAR PEMILIK KENDARAAN")
	printSeparator()

	if len(daftarPemilik) == 0 {
		fmt.Println("  Belum ada data pemilik.")
	} else {
		fmt.Printf("  %-6s %-20s %-15s %-20s\n", "ID", "Nama", "Telepon", "Alamat")
		printLine()
		for _, p := range daftarPemilik {
			fmt.Printf("  %-6s %-20s %-15s %-20s\n", p.ID, p.Nama, p.Telepon, p.Alamat)
		}
	}
	pause()
}

func tambahPemilik() {
	clearScreen()
	printSeparator()
	fmt.Println("           TAMBAH PEMILIK BARU")
	printSeparator()

	pemilik := Pemilik{
		ID:      generateID("P", len(daftarPemilik)),
		Nama:    input("  Nama        : "),
		Telepon: input("  Telepon     : "),
		Alamat:  input("  Alamat      : "),
	}

	daftarPemilik = append(daftarPemilik, pemilik)
	fmt.Printf("\n  ✓ Pemilik berhasil ditambahkan dengan ID: %s\n", pemilik.ID)
	pause()
}

func ubahPemilik() {
	clearScreen()
	printSeparator()
	fmt.Println("           UBAH DATA PEMILIK")
	printSeparator()

	id := input("  Masukkan ID Pemilik: ")
	pemilik := cariPemilikByID(id)

	if pemilik == nil {
		fmt.Println("  ✗ Pemilik tidak ditemukan!")
		pause()
		return
	}

	fmt.Printf("\n  Data saat ini: %s | %s | %s\n", pemilik.Nama, pemilik.Telepon, pemilik.Alamat)
	fmt.Println("  (Kosongkan untuk tidak mengubah)")
	printLine()

	if nama := input("  Nama baru     : "); nama != "" {
		pemilik.Nama = nama
	}
	if telp := input("  Telepon baru  : "); telp != "" {
		pemilik.Telepon = telp
	}
	if alamat := input("  Alamat baru   : "); alamat != "" {
		pemilik.Alamat = alamat
	}

	fmt.Println("\n  ✓ Data pemilik berhasil diperbarui!")
	pause()
}

func hapusPemilik() {
	clearScreen()
	printSeparator()
	fmt.Println("           HAPUS PEMILIK")
	printSeparator()

	id := input("  Masukkan ID Pemilik: ")

	for i, p := range daftarPemilik {
		if p.ID == id {
			konfirmasi := input(fmt.Sprintf("  Hapus pemilik '%s'? (y/n): ", p.Nama))
			if strings.ToLower(konfirmasi) == "y" {
				daftarPemilik = append(daftarPemilik[:i], daftarPemilik[i+1:]...)
				fmt.Println("  ✓ Pemilik berhasil dihapus!")
			} else {
				fmt.Println("  Penghapusan dibatalkan.")
			}
			pause()
			return
		}
	}

	fmt.Println("  ✗ Pemilik tidak ditemukan!")
	pause()
}

// ===================== MENU KENDARAAN =====================

func menuKendaraan() {
	for {
		clearScreen()
		printSeparator()
		fmt.Println("         MANAJEMEN DATA KENDARAAN")
		printSeparator()
		fmt.Println("  1. Tampilkan Semua Kendaraan")
		fmt.Println("  2. Tambah Kendaraan")
		fmt.Println("  3. Ubah Kendaraan")
		fmt.Println("  4. Hapus Kendaraan")
		fmt.Println("  0. Kembali")
		printLine()
		pilihan := input("Pilih menu: ")

		switch pilihan {
		case "1":
			tampilkanSemuaKendaraan()
		case "2":
			tambahKendaraan()
		case "3":
			ubahKendaraan()
		case "4":
			hapusKendaraan()
		case "0":
			return
		default:
			fmt.Println("Pilihan tidak valid!")
			pause()
		}
	}
}

func tampilkanSemuaKendaraan() {
	clearScreen()
	printSeparator()
	fmt.Println("           DAFTAR KENDARAAN")
	printSeparator()

	if len(daftarKendaraan) == 0 {
		fmt.Println("  Belum ada data kendaraan.")
	} else {
		fmt.Printf("  %-6s %-12s %-10s %-10s %-6s %-10s %-12s\n",
			"ID", "Plat", "Merek", "Model", "Tahun", "Pemilik", "Servis Terakhir")
		printLine()
		for _, k := range daftarKendaraan {
			pemilik := cariPemilikByID(k.PemilikID)
			namaPemilik := "-"
			if pemilik != nil {
				namaPemilik = pemilik.Nama
			}
			servisStr := k.TanggalServisTerakhir
			if servisStr == "" {
				servisStr = "Belum Servis"
			}
			fmt.Printf("  %-6s %-12s %-10s %-10s %-6d %-10s %-12s\n",
				k.ID, k.PlatNomor, k.Merek, k.Model, k.TahunProduksi, namaPemilik, servisStr)
		}
	}
	pause()
}

func tambahKendaraan() {
	clearScreen()
	printSeparator()
	fmt.Println("           TAMBAH KENDARAAN BARU")
	printSeparator()

	tahunStr := input("  Tahun Produksi: ")
	tahun, err := strconv.Atoi(tahunStr)
	if err != nil {
		fmt.Println("  ✗ Tahun tidak valid!")
		pause()
		return
	}

	pemilikID := input("  ID Pemilik    : ")
	if cariPemilikByID(pemilikID) == nil {
		fmt.Println("  ✗ ID Pemilik tidak ditemukan!")
		pause()
		return
	}

	kendaraan := Kendaraan{
		ID:            generateID("K", len(daftarKendaraan)),
		PlatNomor:     strings.ToUpper(input("  Plat Nomor   : ")),
		Merek:         input("  Merek        : "),
		Model:         input("  Model        : "),
		TahunProduksi: tahun,
		PemilikID:     pemilikID,
		TanggalServisTerakhir: "",
	}

	daftarKendaraan = append(daftarKendaraan, kendaraan)
	fmt.Printf("\n  ✓ Kendaraan berhasil ditambahkan dengan ID: %s\n", kendaraan.ID)
	pause()
}

func ubahKendaraan() {
	clearScreen()
	printSeparator()
	fmt.Println("           UBAH DATA KENDARAAN")
	printSeparator()

	id := input("  Masukkan ID Kendaraan: ")
	kendaraan := cariKendaraanByID(id)

	if kendaraan == nil {
		fmt.Println("  ✗ Kendaraan tidak ditemukan!")
		pause()
		return
	}

	fmt.Printf("\n  Data saat ini: %s | %s | %s | %d\n",
		kendaraan.PlatNomor, kendaraan.Merek, kendaraan.Model, kendaraan.TahunProduksi)
	fmt.Println("  (Kosongkan untuk tidak mengubah)")
	printLine()

	if plat := input("  Plat Nomor baru  : "); plat != "" {
		kendaraan.PlatNomor = strings.ToUpper(plat)
	}
	if merek := input("  Merek baru       : "); merek != "" {
		kendaraan.Merek = merek
	}
	if model := input("  Model baru       : "); model != "" {
		kendaraan.Model = model
	}
	if tahunStr := input("  Tahun baru       : "); tahunStr != "" {
		if tahun, err := strconv.Atoi(tahunStr); err == nil {
			kendaraan.TahunProduksi = tahun
		}
	}
	if pemilikID := input("  ID Pemilik baru  : "); pemilikID != "" {
		if cariPemilikByID(pemilikID) != nil {
			kendaraan.PemilikID = pemilikID
		} else {
			fmt.Println("  ⚠ ID Pemilik tidak ditemukan, tidak diubah.")
		}
	}

	fmt.Println("\n  ✓ Data kendaraan berhasil diperbarui!")
	pause()
}

func hapusKendaraan() {
	clearScreen()
	printSeparator()
	fmt.Println("           HAPUS KENDARAAN")
	printSeparator()

	id := input("  Masukkan ID Kendaraan: ")

	for i, k := range daftarKendaraan {
		if k.ID == id {
			konfirmasi := input(fmt.Sprintf("  Hapus kendaraan '%s - %s'? (y/n): ", k.Merek, k.PlatNomor))
			if strings.ToLower(konfirmasi) == "y" {
				daftarKendaraan = append(daftarKendaraan[:i], daftarKendaraan[i+1:]...)
				fmt.Println("  ✓ Kendaraan berhasil dihapus!")
			} else {
				fmt.Println("  Penghapusan dibatalkan.")
			}
			pause()
			return
		}
	}

	fmt.Println("  ✗ Kendaraan tidak ditemukan!")
	pause()
}

// ===================== MENU RIWAYAT SERVIS =====================

func menuServis() {
	for {
		clearScreen()
		printSeparator()
		fmt.Println("         MANAJEMEN RIWAYAT SERVIS")
		printSeparator()
		fmt.Println("  1. Tampilkan Semua Riwayat Servis")
		fmt.Println("  2. Riwayat Servis per Kendaraan")
		fmt.Println("  3. Tambah Riwayat Servis")
		fmt.Println("  4. Hapus Riwayat Servis")
		fmt.Println("  0. Kembali")
		printLine()
		pilihan := input("Pilih menu: ")

		switch pilihan {
		case "1":
			tampilkanSemuaServis()
		case "2":
			riwayatServisPerKendaraan()
		case "3":
			tambahRiwayatServis()
		case "4":
			hapusRiwayatServis()
		case "0":
			return
		default:
			fmt.Println("Pilihan tidak valid!")
			pause()
		}
	}
}

func tampilkanSemuaServis() {
	clearScreen()
	printSeparator()
	fmt.Println("           SEMUA RIWAYAT SERVIS")
	printSeparator()

	if len(daftarServis) == 0 {
		fmt.Println("  Belum ada riwayat servis.")
	} else {
		fmt.Printf("  %-6s %-8s %-12s %-20s %-12s\n",
			"ID", "KenID", "Tanggal", "Jenis Kerusakan", "Biaya")
		printLine()
		for _, s := range daftarServis {
			fmt.Printf("  %-6s %-8s %-12s %-20s Rp%-12.0f\n",
				s.ID, s.KendaraanID, s.TanggalServis, s.JenisKerusakan, s.Biaya)
		}
	}
	pause()
}

func riwayatServisPerKendaraan() {
	clearScreen()
	printSeparator()
	fmt.Println("        RIWAYAT SERVIS PER KENDARAAN")
	printSeparator()

	id := input("  Masukkan ID Kendaraan: ")
	kendaraan := cariKendaraanByID(id)

	if kendaraan == nil {
		fmt.Println("  ✗ Kendaraan tidak ditemukan!")
		pause()
		return
	}

	fmt.Printf("\n  Kendaraan: %s %s (%s)\n", kendaraan.Merek, kendaraan.Model, kendaraan.PlatNomor)
	printLine()

	found := false
	total := 0.0
	for _, s := range daftarServis {
		if s.KendaraanID == id {
			fmt.Printf("  [%s] %s | %s\n", s.ID, s.TanggalServis, s.JenisKerusakan)
			fmt.Printf("         Deskripsi: %s\n", s.Deskripsi)
			fmt.Printf("         Teknisi  : %s | Biaya: Rp%.0f\n", s.Teknisi, s.Biaya)
			printLine()
			total += s.Biaya
			found = true
		}
	}

	if !found {
		fmt.Println("  Belum ada riwayat servis untuk kendaraan ini.")
	} else {
		fmt.Printf("  Total Biaya Servis: Rp%.0f\n", total)
	}
	pause()
}

func tambahRiwayatServis() {
	clearScreen()
	printSeparator()
	fmt.Println("           TAMBAH RIWAYAT SERVIS")
	printSeparator()

	kendaraanID := input("  ID Kendaraan     : ")
	kendaraan := cariKendaraanByID(kendaraanID)
	if kendaraan == nil {
		fmt.Println("  ✗ Kendaraan tidak ditemukan!")
		pause()
		return
	}

	biayaStr := input("  Biaya Servis (Rp): ")
	biaya, err := strconv.ParseFloat(biayaStr, 64)
	if err != nil {
		biaya = 0
	}

	tanggal := time.Now().Format("2006-01-02")
	tgl := input(fmt.Sprintf("  Tanggal (%s): ", tanggal))
	if tgl == "" {
		tgl = tanggal
	}

	servis := RiwayatServis{
		ID:             generateID("S", len(daftarServis)),
		KendaraanID:    kendaraanID,
		TanggalServis:  tgl,
		JenisKerusakan: input("  Jenis Kerusakan  : "),
		Deskripsi:      input("  Deskripsi        : "),
		Biaya:          biaya,
		Teknisi:        input("  Nama Teknisi     : "),
	}

	daftarServis = append(daftarServis, servis)

	// Update tanggal servis terakhir pada kendaraan
	kendaraan.TanggalServisTerakhir = tgl

	fmt.Printf("\n  ✓ Riwayat servis berhasil ditambahkan dengan ID: %s\n", servis.ID)
	pause()
}

func hapusRiwayatServis() {
	clearScreen()
	printSeparator()
	fmt.Println("           HAPUS RIWAYAT SERVIS")
	printSeparator()

	id := input("  Masukkan ID Servis: ")

	for i, s := range daftarServis {
		if s.ID == id {
			konfirmasi := input(fmt.Sprintf("  Hapus riwayat servis '%s' tanggal %s? (y/n): ", s.JenisKerusakan, s.TanggalServis))
			if strings.ToLower(konfirmasi) == "y" {
				daftarServis = append(daftarServis[:i], daftarServis[i+1:]...)
				fmt.Println("  ✓ Riwayat servis berhasil dihapus!")
			} else {
				fmt.Println("  Penghapusan dibatalkan.")
			}
			pause()
			return
		}
	}

	fmt.Println("  ✗ Riwayat servis tidak ditemukan!")
	pause()
}

// ===================== MENU PENCARIAN =====================

func menuPencarian() {
	clearScreen()
	printSeparator()
	fmt.Println("         PENCARIAN KENDARAAN")
	printSeparator()
	fmt.Println("  1. Sequential Search (berdasarkan plat nomor)")
	fmt.Println("  2. Binary Search (berdasarkan plat nomor)")
	fmt.Println("  0. Kembali")
	printLine()
	pilihan := input("Pilih metode pencarian: ")

	if pilihan == "0" {
		return
	}

	if pilihan != "1" && pilihan != "2" {
		fmt.Println("Pilihan tidak valid!")
		pause()
		return
	}

	plat := input("\n  Masukkan Plat Nomor: ")
	var hasil *Kendaraan
	var metode string

	if pilihan == "1" {
		metode = "Sequential Search"
		hasil = sequentialSearch(plat)
	} else {
		metode = "Binary Search"
		hasil = binarySearch(plat)
	}

	fmt.Printf("\n  Metode: %s\n", metode)
	printLine()

	if hasil == nil {
		fmt.Printf("  ✗ Kendaraan dengan plat '%s' tidak ditemukan.\n", plat)
	} else {
		pemilik := cariPemilikByID(hasil.PemilikID)
		namaPemilik := "Tidak diketahui"
		if pemilik != nil {
			namaPemilik = pemilik.Nama
		}
		fmt.Println("  ✓ Kendaraan ditemukan!")
		printLine()
		fmt.Printf("  ID            : %s\n", hasil.ID)
		fmt.Printf("  Plat Nomor    : %s\n", hasil.PlatNomor)
		fmt.Printf("  Merek/Model   : %s %s\n", hasil.Merek, hasil.Model)
		fmt.Printf("  Tahun Produksi: %d\n", hasil.TahunProduksi)
		fmt.Printf("  Pemilik       : %s\n", namaPemilik)
		if hasil.TanggalServisTerakhir != "" {
			fmt.Printf("  Servis Terakhir: %s\n", hasil.TanggalServisTerakhir)
		} else {
			fmt.Printf("  Servis Terakhir: Belum pernah servis\n")
		}
	}
	pause()
}

// ===================== MENU PENGURUTAN =====================

func menuPengurutan() {
	clearScreen()
	printSeparator()
	fmt.Println("         PENGURUTAN KENDARAAN")
	printSeparator()
	fmt.Println("  1. Selection Sort - Tahun Produksi (Ascending)")
	fmt.Println("  2. Insertion Sort - Tanggal Servis Terakhir (Ascending)")
	fmt.Println("  0. Kembali")
	printLine()
	pilihan := input("Pilih metode pengurutan: ")

	var sorted []Kendaraan
	var judulSort string

	switch pilihan {
	case "1":
		sorted = selectionSortByTahun(daftarKendaraan)
		judulSort = "Selection Sort - Tahun Produksi"
	case "2":
		sorted = insertionSortByTanggalServis(daftarKendaraan)
		judulSort = "Insertion Sort - Tanggal Servis"
	case "0":
		return
	default:
		fmt.Println("Pilihan tidak valid!")
		pause()
		return
	}

	clearScreen()
	printSeparator()
	fmt.Printf("  Hasil: %s\n", judulSort)
	printSeparator()
	fmt.Printf("  %-6s %-12s %-10s %-10s %-6s %-12s\n",
		"ID", "Plat", "Merek", "Model", "Tahun", "Servis Terakhir")
	printLine()

	for _, k := range sorted {
		servisStr := k.TanggalServisTerakhir
		if servisStr == "" {
			servisStr = "Belum Servis"
		}
		fmt.Printf("  %-6s %-12s %-10s %-10s %-6d %-12s\n",
			k.ID, k.PlatNomor, k.Merek, k.Model, k.TahunProduksi, servisStr)
	}
	pause()
}

// ===================== STATISTIK =====================

func menuStatistik() {
	clearScreen()
	printSeparator()
	fmt.Println("         STATISTIK SERVIS KENDARAAN")
	printSeparator()

	// Statistik jumlah kendaraan servis per bulan
	fmt.Println("\n  [1] KENDARAAN DISERVIS PER BULAN")
	printLine()

	servisPerBulan := make(map[string]int)
	for _, s := range daftarServis {
		if len(s.TanggalServis) >= 7 {
			bulan := s.TanggalServis[:7] // format YYYY-MM
			servisPerBulan[bulan]++
		}
	}

	if len(servisPerBulan) == 0 {
		fmt.Println("  Belum ada data servis.")
	} else {
		// Urutkan bulan
		bulanList := make([]string, 0, len(servisPerBulan))
		for k := range servisPerBulan {
			bulanList = append(bulanList, k)
		}
		sort.Strings(bulanList)

		maxVal := 0
		for _, v := range servisPerBulan {
			if v > maxVal {
				maxVal = v
			}
		}

		for _, bulan := range bulanList {
			jumlah := servisPerBulan[bulan]
			bar := strings.Repeat("█", jumlah*20/maxVal)
			fmt.Printf("  %s | %-20s %d unit\n", bulan, bar, jumlah)
		}
	}

	// Statistik jenis kerusakan paling sering
	fmt.Println("\n  [2] JENIS KERUSAKAN PALING SERING")
	printLine()

	kerusakanCount := make(map[string]int)
	for _, s := range daftarServis {
		kerusakanCount[s.JenisKerusakan]++
	}

	if len(kerusakanCount) == 0 {
		fmt.Println("  Belum ada data kerusakan.")
	} else {
		// Urutkan berdasarkan jumlah
		type KvPair struct {
			Key   string
			Value int
		}
		pairs := make([]KvPair, 0, len(kerusakanCount))
		for k, v := range kerusakanCount {
			pairs = append(pairs, KvPair{k, v})
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Value > pairs[j].Value
		})

		for rank, p := range pairs {
			marker := "  "
			if rank == 0 {
				marker = "🥇"
			} else if rank == 1 {
				marker = "🥈"
			} else if rank == 2 {
				marker = "🥉"
			}
			fmt.Printf(" %s %-25s : %d kali\n", marker, p.Key, p.Value)
		}
	}

	// Total statistik
	fmt.Println("\n  [3] RINGKASAN UMUM")
	printLine()
	fmt.Printf("  Total Pemilik    : %d\n", len(daftarPemilik))
	fmt.Printf("  Total Kendaraan  : %d\n", len(daftarKendaraan))
	fmt.Printf("  Total Servis     : %d\n", len(daftarServis))

	totalBiaya := 0.0
	for _, s := range daftarServis {
		totalBiaya += s.Biaya
	}
	fmt.Printf("  Total Pendapatan : Rp%.0f\n", totalBiaya)

	pause()
}

// ===================== DATA DEMO =====================

func loadDataDemo() {
	// Data pemilik
	daftarPemilik = []Pemilik{
		{ID: "P001", Nama: "Budi Santoso", Telepon: "081234567890", Alamat: "Jl. Merdeka No. 1"},
		{ID: "P002", Nama: "Siti Rahayu", Telepon: "082345678901", Alamat: "Jl. Sudirman No. 5"},
		{ID: "P003", Nama: "Ahmad Fauzi", Telepon: "083456789012", Alamat: "Jl. Gatot Subroto No. 10"},
	}

	// Data kendaraan
	daftarKendaraan = []Kendaraan{
		{ID: "K001", PlatNomor: "R1234AB", Merek: "Toyota", Model: "Avanza", TahunProduksi: 2020, PemilikID: "P001", TanggalServisTerakhir: "2024-03-15"},
		{ID: "K002", PlatNomor: "R5678CD", Merek: "Honda", Model: "Beat", TahunProduksi: 2019, PemilikID: "P001", TanggalServisTerakhir: "2024-01-20"},
		{ID: "K003", PlatNomor: "H2468EF", Merek: "Yamaha", Model: "NMAX", TahunProduksi: 2022, PemilikID: "P002", TanggalServisTerakhir: "2024-02-10"},
		{ID: "K004", PlatNomor: "AA9999GH", Merek: "Suzuki", Model: "Ertiga", TahunProduksi: 2018, PemilikID: "P003", TanggalServisTerakhir: "2023-12-05"},
		{ID: "K005", PlatNomor: "B1111IJ", Merek: "Daihatsu", Model: "Xenia", TahunProduksi: 2021, PemilikID: "P002", TanggalServisTerakhir: ""},
	}

	// Data riwayat servis
	daftarServis = []RiwayatServis{
		{ID: "S001", KendaraanID: "K001", TanggalServis: "2024-03-15", JenisKerusakan: "Ganti Oli", Deskripsi: "Ganti oli mesin 4L", Biaya: 250000, Teknisi: "Andi"},
		{ID: "S002", KendaraanID: "K001", TanggalServis: "2023-11-10", JenisKerusakan: "Rem Blong", Deskripsi: "Ganti kampas rem depan belakang", Biaya: 450000, Teknisi: "Budi"},
		{ID: "S003", KendaraanID: "K002", TanggalServis: "2024-01-20", JenisKerusakan: "Ganti Oli", Deskripsi: "Ganti oli mesin + filter", Biaya: 150000, Teknisi: "Andi"},
		{ID: "S004", KendaraanID: "K003", TanggalServis: "2024-02-10", JenisKerusakan: "Tune Up", Deskripsi: "Tune up berkala 10.000 km", Biaya: 350000, Teknisi: "Candra"},
		{ID: "S005", KendaraanID: "K004", TanggalServis: "2023-12-05", JenisKerusakan: "Ganti Ban", Deskripsi: "Ganti 4 ban luar", Biaya: 2000000, Teknisi: "Budi"},
		{ID: "S006", KendaraanID: "K001", TanggalServis: "2024-03-20", JenisKerusakan: "Tune Up", Deskripsi: "Tune up + ganti busi", Biaya: 400000, Teknisi: "Andi"},
	}

	fmt.Println("  ✓ Data demo berhasil dimuat!")
}

// ===================== MAIN MENU =====================

func main() {
	clearScreen()
	printSeparator()
	fmt.Println("  ██████╗ ██╗   ██╗████████╗ ██████╗  ██████╗ █████╗ ██████╗ ███████╗")
	fmt.Println("  ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝")
	fmt.Println("  ███████║██║   ██║   ██║   ██║   ██║██║     ███████║██████╔╝█████╗  ")
	fmt.Println("  ██╔══██║██║   ██║   ██║   ██║   ██║██║     ██╔══██║██╔══██╗██╔══╝  ")
	fmt.Println("  ██║  ██║╚██████╔╝   ██║   ╚██████╔╝╚██████╗██║  ██║██║  ██║███████╗")
	fmt.Println("  ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝  ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝")
	printSeparator()
	fmt.Println("  Aplikasi Manajemen dan Riwayat Servis Kendaraan")
	fmt.Println("  Tugas Besar Algoritma Pemrograman 2 - Telkom University")
	printSeparator()

	muat := input("\n  Muat data demo? (y/n): ")
	if strings.ToLower(muat) == "y" {
		loadDataDemo()
	}

	for {
		clearScreen()
		printSeparator()
		fmt.Println("               MENU UTAMA - AutoCare")
		printSeparator()
		fmt.Println("  1. Manajemen Pemilik")
		fmt.Println("  2. Manajemen Kendaraan")
		fmt.Println("  3. Riwayat Servis")
		fmt.Println("  4. Pencarian Kendaraan (Sequential & Binary Search)")
		fmt.Println("  5. Pengurutan Kendaraan (Selection & Insertion Sort)")
		fmt.Println("  6. Statistik Servis")
		fmt.Println("  0. Keluar")
		printSeparator()
		pilihan := input("Pilih menu: ")

		switch pilihan {
		case "1":
			menuPemilik()
		case "2":
			menuKendaraan()
		case "3":
			menuServis()
		case "4":
			menuPencarian()
		case "5":
			menuPengurutan()
		case "6":
			menuStatistik()
		case "0":
			clearScreen()
			printSeparator()
			fmt.Println("  Terima kasih telah menggunakan AutoCare!")
			fmt.Println("  Sampai jumpa! 👋")
			printSeparator()
			os.Exit(0)
		default:
			fmt.Println("  Pilihan tidak valid!")
			pause()
		}
	}
}
