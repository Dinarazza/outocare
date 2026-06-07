package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aquasecurity/table"
)

type Pemilik struct {
	ID           int
	NamaPemilik  string
	NomorTelepon string
	Alamat       string
	Active       bool
}

type Kendaraan struct {
	PlatNomor     string
	MerkKendaraan string
	TipeKendaraan string
	TahunProduksi int
	IDPemilik     int
	TanggalInput  time.Time
}

type RiwayatServis struct {
	IDServis       int
	PlatNomor      string
	TanggalServis  time.Time
	JenisKerusakan string
	DetailServis   string
	Biaya          int
}

var DaftarKendaraan [100]Kendaraan
var DaftarPemilik [100]Pemilik
var DaftarServis [200]RiwayatServis

func InputText(prompt string) string {
	var scanner = bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func TampilkanError(err error) bool {
	if err != nil {
		fmt.Println("\nERROR:", err)
		return true
	}
	return false
}

func InputAngka(prompt string) int {
	var input string
	var angka int
	var err error
	var valid bool

	valid = false
	for !valid {
		input = InputText(prompt)
		angka, err = strconv.Atoi(input)
		if err == nil {
			valid = true
		} else {
			TampilkanError(errors.New("Input harus berupa angka."))
		}
	}
	return angka
}

func InputTeksWajib(prompt, pesanError string) string {
	var input string
	var valid bool

	valid = false
	for !valid {
		input = InputText(prompt)
		if input != "" {
			valid = true
		} else {
			TampilkanError(errors.New(pesanError))
		}
	}
	return input
}

func GetNamaBulan(nomorBulan int) string {
	var namaBulan [12]string
	namaBulan[0] = "Januari"
	namaBulan[1] = "Februari"
	namaBulan[2] = "Maret"
	namaBulan[3] = "April"
	namaBulan[4] = "Mei"
	namaBulan[5] = "Juni"
	namaBulan[6] = "Juli"
	namaBulan[7] = "Agustus"
	namaBulan[8] = "September"
	namaBulan[9] = "Oktober"
	namaBulan[10] = "November"
	namaBulan[11] = "Desember"

	if nomorBulan >= 1 && nomorBulan <= 12 {
		return namaBulan[nomorBulan-1]
	}
	return "Tidak Diketahui"
}

func CariIndeksPemilikSequential(id int, jumlahPemilik *int) int {
	var i int
	var indeks int

	indeks = -1
	i = 0
	for i < *jumlahPemilik && indeks == -1 {
		if DaftarPemilik[i].ID == id && DaftarPemilik[i].Active {
			indeks = i
		}
		i++
	}
	return indeks
}

func CariNamaPemilik(id int, jumlahPemilik *int) string {
	var i int
	var nama string

	nama = "(Pemilik tidak ditemukan)"
	i = 0
	for i < *jumlahPemilik {
		if DaftarPemilik[i].ID == id && DaftarPemilik[i].Active {
			nama = DaftarPemilik[i].NamaPemilik
			i = *jumlahPemilik // sudah ketemu, hentikan loop
		} else {
			i++
		}
	}
	return nama
}

func CariIndeksKendaraanSequential(plat string, jumlahKendaraan *int) int {
	var i int
	var indeks int

	indeks = -1
	i = 0
	for i < *jumlahKendaraan && indeks == -1 {
		if strings.ToUpper(DaftarKendaraan[i].PlatNomor) == strings.ToUpper(plat) {
			indeks = i
		}
		i++
	}
	return indeks
}

func CariKendaraanSequential(keyword string, jumlahKendaraan *int, jumlahPemilik *int) {
	var i, ditemukan int
	var t *table.Table
	var k Kendaraan
	var namaPemilik string

	keyword = strings.ToLower(keyword)
	ditemukan = 0

	t = table.New(os.Stdout)
	t.SetRowLines(false)
	t.SetHeaders("No", "Plat Nomor", "Merk", "Tipe", "Tahun Produksi", "Nama Pemilik")

	for i = 0; i < *jumlahKendaraan; i++ {
		if strings.Contains(strings.ToLower(DaftarKendaraan[i].PlatNomor), keyword) {
			k = DaftarKendaraan[i]
			namaPemilik = CariNamaPemilik(k.IDPemilik, jumlahPemilik)
			t.AddRow(
				strconv.Itoa(i+1),
				k.PlatNomor,
				k.MerkKendaraan,
				k.TipeKendaraan,
				strconv.Itoa(k.TahunProduksi),
				namaPemilik,
			)
			ditemukan++
		}
	}

	if ditemukan == 0 {
		TampilkanError(errors.New("Kendaraan tidak ditemukan."))
	} else {
		fmt.Println("\n================= HASIL PENCARIAN =================")
		fmt.Println()
		t.Render()
		fmt.Printf("\nDitemukan %d kendaraan.\n", ditemukan)
	}
}

func CariKendaraanBinary(keyword string, jumlahKendaraan *int, jumlahPemilik *int) {
	var kiri, tengah, kanan int
	var ditemukan bool
	var t *table.Table
	var namaPemilik string
	var platTengah, platKeyword string

	UrutkanKendaraanByPlatSelectionSort(jumlahKendaraan)

	platKeyword = strings.ToUpper(keyword)
	kiri = 0
	kanan = *jumlahKendaraan - 1
	ditemukan = false
	tengah = 0

	for kiri <= kanan && !ditemukan {
		tengah = (kiri + kanan) / 2
		platTengah = strings.ToUpper(DaftarKendaraan[tengah].PlatNomor)

		if platTengah == platKeyword {
			ditemukan = true
		} else if platTengah < platKeyword {
			kiri = tengah + 1
		} else {
			kanan = tengah - 1
		}
	}

	if !ditemukan {
		TampilkanError(errors.New("Kendaraan dengan plat nomor tersebut tidak ditemukan."))
	} else {
		fmt.Println("\n================= HASIL PENCARIAN =================")
		t = table.New(os.Stdout)
		t.SetRowLines(false)
		t.SetHeaders("No", "Plat Nomor", "Merk", "Tipe", "Tahun Produksi", "Nama Pemilik")

		namaPemilik = CariNamaPemilik(DaftarKendaraan[tengah].IDPemilik, jumlahPemilik)
		t.AddRow(
			strconv.Itoa(tengah+1),
			DaftarKendaraan[tengah].PlatNomor,
			DaftarKendaraan[tengah].MerkKendaraan,
			DaftarKendaraan[tengah].TipeKendaraan,
			strconv.Itoa(DaftarKendaraan[tengah].TahunProduksi),
			namaPemilik,
		)
		fmt.Println()
		t.Render()
	}
}

func UrutkanKendaraanByPlatSelectionSort(jumlahKendaraan *int) {
	var i, j, indeksMin int
	var temp Kendaraan

	for i = 0; i < *jumlahKendaraan-1; i++ {
		indeksMin = i
		for j = i + 1; j < *jumlahKendaraan; j++ {
			if strings.ToUpper(DaftarKendaraan[j].PlatNomor) < strings.ToUpper(DaftarKendaraan[indeksMin].PlatNomor) {
				indeksMin = j
			}
		}
		if indeksMin != i {
			temp = DaftarKendaraan[i]
			DaftarKendaraan[i] = DaftarKendaraan[indeksMin]
			DaftarKendaraan[indeksMin] = temp
		}
	}
}

func UrutkanKendaraanByTahunSelectionSort(jumlahKendaraan *int) {
	var i, j, indeksMin int
	var temp Kendaraan

	for i = 0; i < *jumlahKendaraan-1; i++ {
		indeksMin = i
		for j = i + 1; j < *jumlahKendaraan; j++ {
			// Ascending: tahun lebih kecil lebih dulu
			if DaftarKendaraan[j].TahunProduksi < DaftarKendaraan[indeksMin].TahunProduksi {
				indeksMin = j
			}
		}
		if indeksMin != i {
			temp = DaftarKendaraan[i]
			DaftarKendaraan[i] = DaftarKendaraan[indeksMin]
			DaftarKendaraan[indeksMin] = temp
		}
	}
}

func GetTanggalServisTerakhir(plat string, jumlahServis *int) time.Time {
	var i int
	var tanggalTerakhir time.Time

	for i = 0; i < *jumlahServis; i++ {
		if strings.ToUpper(DaftarServis[i].PlatNomor) == strings.ToUpper(plat) {
			if DaftarServis[i].TanggalServis.After(tanggalTerakhir) {
				tanggalTerakhir = DaftarServis[i].TanggalServis
			}
		}
	}
	return tanggalTerakhir
}

func UrutkanKendaraanByTanggalServisInsertionSort(jumlahKendaraan *int, jumlahServis *int) {
	var i, j int
	var kendaraanSaatIni Kendaraan
	var tanggalSaatIni, tanggalJ time.Time
	var harus_geser bool

	for i = 1; i < *jumlahKendaraan; i++ {
		kendaraanSaatIni = DaftarKendaraan[i]
		tanggalSaatIni = GetTanggalServisTerakhir(kendaraanSaatIni.PlatNomor, jumlahServis)
		j = i - 1

		harus_geser = true
		for j >= 0 && harus_geser {
			tanggalJ = GetTanggalServisTerakhir(DaftarKendaraan[j].PlatNomor, jumlahServis)
			harus_geser = tanggalJ.Before(tanggalSaatIni)

			if harus_geser {
				DaftarKendaraan[j+1] = DaftarKendaraan[j]
				j--
			}
		}
		DaftarKendaraan[j+1] = kendaraanSaatIni
	}
}

func TampilkanSemuaKendaraan(jumlahKendaraan *int, jumlahPemilik *int) {
	var i int
	var t *table.Table
	var k Kendaraan
	var namaPemilik, teleponPemilik, alamatPemilik string
	var indeksPemilik int

	if *jumlahKendaraan == 0 {
		fmt.Println("\nINFO: Belum ada data kendaraan.")
		return
	}

	t = table.New(os.Stdout)
	t.SetRowLines(false)
	t.SetHeaders("No", "Plat Nomor", "Merk", "Tipe", "Tahun", "Nama Pemilik", "Telepon", "Alamat")

	for i = 0; i < *jumlahKendaraan; i++ {
		k = DaftarKendaraan[i]
		indeksPemilik = CariIndeksPemilikSequential(k.IDPemilik, jumlahPemilik)

		if indeksPemilik != -1 {
			namaPemilik = DaftarPemilik[indeksPemilik].NamaPemilik
			teleponPemilik = DaftarPemilik[indeksPemilik].NomorTelepon
			alamatPemilik = DaftarPemilik[indeksPemilik].Alamat
		} else {
			namaPemilik = "Tidak diketahui"
			teleponPemilik = "-"
			alamatPemilik = "-"
		}

		t.AddRow(
			strconv.Itoa(i+1),
			k.PlatNomor,
			k.MerkKendaraan,
			k.TipeKendaraan,
			strconv.Itoa(k.TahunProduksi),
			namaPemilik,
			teleponPemilik,
			alamatPemilik,
		)
	}
	fmt.Println()
	t.Render()
}

func TampilkanSemuaPemilik(jumlahPemilik *int) {
	var i int
	var t *table.Table
	var p Pemilik

	if *jumlahPemilik == 0 {
		fmt.Println("\nINFO: Belum ada data pemilik.")
		return
	}

	t = table.New(os.Stdout)
	t.SetRowLines(false)
	t.SetHeaders("ID", "Nama Pemilik", "Nomor Telepon", "Alamat")

	for i = 0; i < *jumlahPemilik; i++ {
		if DaftarPemilik[i].Active {
			p = DaftarPemilik[i]
			t.AddRow(
				strconv.Itoa(p.ID),
				p.NamaPemilik,
				p.NomorTelepon,
				p.Alamat,
			)
		}
	}
	fmt.Println()
	t.Render()
}

func TampilkanSemuaServis(jumlahServis *int) {
	var i int
	var t *table.Table
	var s RiwayatServis

	if *jumlahServis == 0 {
		fmt.Println("\nINFO: Belum ada data riwayat servis.")
		return
	}

	t = table.New(os.Stdout)
	t.SetRowLines(false)
	t.SetHeaders("ID", "Plat Nomor", "Tanggal Servis", "Jenis Kerusakan", "Detail Servis", "Biaya (Rp)")

	for i = 0; i < *jumlahServis; i++ {
		s = DaftarServis[i]
		t.AddRow(
			strconv.Itoa(s.IDServis),
			s.PlatNomor,
			s.TanggalServis.Format("02-01-2006"),
			s.JenisKerusakan,
			s.DetailServis,
			strconv.Itoa(s.Biaya),
		)
	}
	fmt.Println()
	t.Render()
}

func SimpanDataKendaraan(jumlahKendaraan *int) {
	var i int
	var dataAktif [100]Kendaraan
	var jumlahAktif int
	var file *os.File
	var encoder *json.Encoder
	var err error

	jumlahAktif = 0
	for i = 0; i < *jumlahKendaraan; i++ {
		dataAktif[jumlahAktif] = DaftarKendaraan[i]
		jumlahAktif++
	}

	file, err = os.Create("kendaraan.json")
	if TampilkanError(err) {
		return
	}
	defer file.Close()
	encoder = json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(dataAktif[:jumlahAktif])
}

func MuatDataKendaraan(jumlahKendaraan *int) {
	var err error
	var data []Kendaraan
	var file *os.File
	var decoder *json.Decoder
	var i int

	file, err = os.Open("kendaraan.json")
	if err != nil {
		return
	}
	defer file.Close()

	decoder = json.NewDecoder(file)
	err = decoder.Decode(&data)
	if TampilkanError(err) {
		return
	}

	for i = 0; i < len(data) && i < 100; i++ {
		DaftarKendaraan[i] = data[i]
	}
	*jumlahKendaraan = len(data)
}

func SimpanDataPemilik(jumlahPemilik *int) {
	var i int
	var dataAktif [100]Pemilik
	var jumlahAktif int
	var file *os.File
	var encoder *json.Encoder
	var err error

	jumlahAktif = 0
	for i = 0; i < *jumlahPemilik; i++ {
		if DaftarPemilik[i].Active {
			dataAktif[jumlahAktif] = DaftarPemilik[i]
			jumlahAktif++
		}
	}

	file, err = os.Create("pemilik.json")
	if TampilkanError(err) {
		return
	}
	defer file.Close()
	encoder = json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(dataAktif[:jumlahAktif])
}

func MuatDataPemilik(jumlahPemilik *int) {
	var err error
	var data []Pemilik
	var file *os.File
	var decoder *json.Decoder
	var i int

	file, err = os.Open("pemilik.json")
	if err != nil {
		return
	}
	defer file.Close()

	decoder = json.NewDecoder(file)
	err = decoder.Decode(&data)
	if TampilkanError(err) {
		return
	}

	for i = 0; i < len(data) && i < 100; i++ {
		DaftarPemilik[i] = data[i]
	}
	*jumlahPemilik = len(data)
}

func SimpanDataServis(jumlahServis *int) {
	var i int
	var dataAktif [200]RiwayatServis
	var jumlahAktif int
	var file *os.File
	var encoder *json.Encoder
	var err error

	jumlahAktif = 0
	for i = 0; i < *jumlahServis; i++ {
		dataAktif[jumlahAktif] = DaftarServis[i]
		jumlahAktif++
	}

	file, err = os.Create("servis.json")
	if TampilkanError(err) {
		return
	}
	defer file.Close()
	encoder = json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(dataAktif[:jumlahAktif])
}

func MuatDataServis(jumlahServis *int) {
	var err error
	var data []RiwayatServis
	var file *os.File
	var decoder *json.Decoder
	var i int

	file, err = os.Open("servis.json")
	if err != nil {
		return
	}
	defer file.Close()

	decoder = json.NewDecoder(file)
	err = decoder.Decode(&data)
	if TampilkanError(err) {
		return
	}

	for i = 0; i < len(data) && i < 200; i++ {
		DaftarServis[i] = data[i]
	}
	*jumlahServis = len(data)
}

func TambahKendaraan(platNomor, merk, tipe string, tahun, idPemilik int, jumlahKendaraan *int) {
	if *jumlahKendaraan >= 100 {
		fmt.Println("\nERROR: Kapasitas kendaraan penuh.")
		return
	}

	DaftarKendaraan[*jumlahKendaraan].PlatNomor = platNomor
	DaftarKendaraan[*jumlahKendaraan].MerkKendaraan = merk
	DaftarKendaraan[*jumlahKendaraan].TipeKendaraan = tipe
	DaftarKendaraan[*jumlahKendaraan].TahunProduksi = tahun
	DaftarKendaraan[*jumlahKendaraan].IDPemilik = idPemilik
	DaftarKendaraan[*jumlahKendaraan].TanggalInput = time.Now()

	*jumlahKendaraan++
	fmt.Println("\nINFO: Data kendaraan berhasil ditambahkan.")
}

func EditKendaraan(indeks int, platNomor, merk, tipe string, tahun, idPemilik int) {
	var i int
	i = indeks - 1
	DaftarKendaraan[i].PlatNomor = platNomor
	DaftarKendaraan[i].MerkKendaraan = merk
	DaftarKendaraan[i].TipeKendaraan = tipe
	DaftarKendaraan[i].TahunProduksi = tahun
	DaftarKendaraan[i].IDPemilik = idPemilik
	fmt.Println("\nINFO: Data kendaraan berhasil diubah.")
}

func HapusKendaraan(indeks int, jumlahKendaraan *int) {
	var i, posisi int
	posisi = indeks - 1

	for i = posisi; i < *jumlahKendaraan-1; i++ {
		DaftarKendaraan[i] = DaftarKendaraan[i+1]
	}
	DaftarKendaraan[*jumlahKendaraan-1] = Kendaraan{}
	*jumlahKendaraan--
	fmt.Println("\nINFO: Data kendaraan berhasil dihapus.")
}

func TambahPemilik(nama, telepon, alamat string, jumlahPemilik *int) {
	var i, idMaks, idBaru int

	if *jumlahPemilik >= 100 {
		fmt.Println("\nERROR: Kapasitas pemilik penuh.")
		return
	}

	idMaks = 0
	for i = 0; i < *jumlahPemilik; i++ {
		if DaftarPemilik[i].ID > idMaks {
			idMaks = DaftarPemilik[i].ID
		}
	}
	idBaru = idMaks + 1

	DaftarPemilik[*jumlahPemilik].ID = idBaru
	DaftarPemilik[*jumlahPemilik].NamaPemilik = nama
	DaftarPemilik[*jumlahPemilik].NomorTelepon = telepon
	DaftarPemilik[*jumlahPemilik].Alamat = alamat
	DaftarPemilik[*jumlahPemilik].Active = true

	*jumlahPemilik++
	fmt.Printf("\nINFO: Pemilik berhasil ditambahkan dengan ID %d.\n", idBaru)
}

func EditPemilik(id int, nama, telepon, alamat string, jumlahPemilik *int) {
	var indeks int

	indeks = CariIndeksPemilikSequential(id, jumlahPemilik)
	if indeks == -1 {
		TampilkanError(errors.New("Pemilik tidak ditemukan."))
		return
	}

	DaftarPemilik[indeks].NamaPemilik = nama
	DaftarPemilik[indeks].NomorTelepon = telepon
	DaftarPemilik[indeks].Alamat = alamat
	fmt.Println("\nINFO: Data pemilik berhasil diubah.")
}

func HapusPemilik(id int, jumlahKendaraan *int, jumlahPemilik *int) {
	var i, indeks int
	var masihPunyaKendaraan bool

	masihPunyaKendaraan = false
	i = 0
	for i < *jumlahKendaraan && !masihPunyaKendaraan {
		if DaftarKendaraan[i].IDPemilik == id {
			masihPunyaKendaraan = true
		}
		i++
	}

	if masihPunyaKendaraan {
		TampilkanError(errors.New("Pemilik masih memiliki kendaraan. Hapus kendaraan terlebih dahulu."))
		return
	}

	indeks = CariIndeksPemilikSequential(id, jumlahPemilik)
	if indeks == -1 {
		TampilkanError(errors.New("Pemilik tidak ditemukan."))
		return
	}

	DaftarPemilik[indeks].Active = false
	fmt.Println("\nINFO: Data pemilik berhasil dihapus.")
}

func TambahServis(platNomor, jenisKerusakan, detailServis string, biaya int, jumlahServis *int) {
	var i, idMaks, idBaru int

	if *jumlahServis >= 200 {
		fmt.Println("\nERROR: Kapasitas riwayat servis penuh.")
		return
	}

	idMaks = 0
	for i = 0; i < *jumlahServis; i++ {
		if DaftarServis[i].IDServis > idMaks {
			idMaks = DaftarServis[i].IDServis
		}
	}
	idBaru = idMaks + 1

	DaftarServis[*jumlahServis].IDServis = idBaru
	DaftarServis[*jumlahServis].PlatNomor = strings.ToUpper(platNomor)
	DaftarServis[*jumlahServis].TanggalServis = time.Now()
	DaftarServis[*jumlahServis].JenisKerusakan = jenisKerusakan
	DaftarServis[*jumlahServis].DetailServis = detailServis
	DaftarServis[*jumlahServis].Biaya = biaya

	*jumlahServis++
	fmt.Println("\nINFO: Riwayat servis berhasil ditambahkan.")
}

func HapusServis(indeks int, jumlahServis *int) {
	var i, posisi int
	posisi = indeks - 1

	for i = posisi; i < *jumlahServis-1; i++ {
		DaftarServis[i] = DaftarServis[i+1]
	}
	DaftarServis[*jumlahServis-1] = RiwayatServis{}
	*jumlahServis--
	fmt.Println("\nINFO: Riwayat servis berhasil dihapus.")
}

func TampilkanStatistik(jumlahServis *int) {
	var i, j int
	var s RiwayatServis
	var t *table.Table
	var labelBulan string

	var namaBulanArray [100]string
	var jumlahPerBulan [100]int
	var jumlahBulan int

	var namaKerusakanArray [50]string
	var jumlahPerKerusakan [50]int
	var jumlahJenisKerusakan int

	var sudahAda bool
	var tempNama string
	var tempJumlah int

	if *jumlahServis == 0 {
		fmt.Println("\nINFO: Belum ada riwayat servis.")
		return
	}

	fmt.Println("\n===== STATISTIK SERVIS KENDARAAN =====")

	jumlahBulan = 0
	jumlahJenisKerusakan = 0

	for i = 0; i < *jumlahServis; i++ {
		s = DaftarServis[i]
		labelBulan = GetNamaBulan(int(s.TanggalServis.Month())) + " " + strconv.Itoa(s.TanggalServis.Year())

		sudahAda = false
		j = 0
		for j < jumlahBulan && !sudahAda {
			if namaBulanArray[j] == labelBulan {
				jumlahPerBulan[j]++
				sudahAda = true
			}
			j++
		}
		if !sudahAda && jumlahBulan < 100 {
			namaBulanArray[jumlahBulan] = labelBulan
			jumlahPerBulan[jumlahBulan] = 1
			jumlahBulan++
		}

		sudahAda = false
		j = 0
		for j < jumlahJenisKerusakan && !sudahAda {
			if namaKerusakanArray[j] == s.JenisKerusakan {
				jumlahPerKerusakan[j]++
				sudahAda = true
			}
			j++
		}
		if !sudahAda && jumlahJenisKerusakan < 50 {
			namaKerusakanArray[jumlahJenisKerusakan] = s.JenisKerusakan
			jumlahPerKerusakan[jumlahJenisKerusakan] = 1
			jumlahJenisKerusakan++
		}
	}

	fmt.Println("\n--- JUMLAH KENDARAAN DISERVIS PER BULAN ---")
	t = table.New(os.Stdout)
	t.SetRowLines(false)
	t.SetHeaders("Bulan", "Jumlah")
	for i = 0; i < jumlahBulan; i++ {
		t.AddRow(namaBulanArray[i], strconv.Itoa(jumlahPerBulan[i]))
	}
	fmt.Println()
	t.Render()

	for i = 0; i < jumlahJenisKerusakan-1; i++ {
		for j = i + 1; j < jumlahJenisKerusakan; j++ {
			if jumlahPerKerusakan[j] > jumlahPerKerusakan[i] {
				tempJumlah = jumlahPerKerusakan[i]
				jumlahPerKerusakan[i] = jumlahPerKerusakan[j]
				jumlahPerKerusakan[j] = tempJumlah

				tempNama = namaKerusakanArray[i]
				namaKerusakanArray[i] = namaKerusakanArray[j]
				namaKerusakanArray[j] = tempNama
			}
		}
	}

	fmt.Println("\n--- KATEGORI KERUSAKAN YANG PALING SERING ---")
	t = table.New(os.Stdout)
	t.SetRowLines(false)
	t.SetHeaders("No", "Jenis Kerusakan", "Jumlah")
	for i = 0; i < jumlahJenisKerusakan; i++ {
		t.AddRow(strconv.Itoa(i+1), namaKerusakanArray[i], strconv.Itoa(jumlahPerKerusakan[i]))
	}
	fmt.Println()
	t.Render()
}

func KelolaKendaraan(jumlahKendaraan *int, jumlahPemilik *int) {
	var pilihan int
	var platNomor, merk, tipe, konfirmasi string
	var tahun, idPemilik, indeksKendaraan int
	var tahunInput, idInput, indexInput string
	var err error
	var idx int

	fmt.Println("\n\n============= KELOLA DATA KENDARAAN =============")

	for {
		fmt.Println("\nMenu kelola:")
		fmt.Println("1. Lihat semua kendaraan")
		fmt.Println("2. Tambah kendaraan")
		fmt.Println("3. Edit kendaraan")
		fmt.Println("4. Hapus kendaraan")
		fmt.Println("0. Kembali ke menu utama")

		pilihan = InputAngka("\n? Pilih menu\n> ")

		switch pilihan {
		case 1:
			fmt.Println("\n=========== DAFTAR SEMUA KENDARAAN ===========")
			TampilkanSemuaKendaraan(jumlahKendaraan, jumlahPemilik)

		case 2:
			fmt.Println("\n============== TAMBAH KENDARAAN ==============")
			if *jumlahPemilik == 0 {
				fmt.Println("\nERROR: Belum ada data pemilik. Silakan tambah pemilik terlebih dahulu.")
			} else {
				platNomor = strings.ToUpper(InputTeksWajib("\n? Plat Nomor (contoh: B 1234 AB)\n> ", "Plat nomor tidak boleh kosong."))
				merk = InputTeksWajib("\n? Merk Kendaraan (contoh: Toyota)\n> ", "Merk kendaraan tidak boleh kosong.")
				tipe = InputTeksWajib("\n? Tipe Kendaraan (contoh: Avanza)\n> ", "Tipe kendaraan tidak boleh kosong.")
				tahun = 0
				for tahun < 1900 || tahun > time.Now().Year() {
					tahun = InputAngka("\n? Tahun Produksi (contoh: 2020)\n> ")
					if tahun < 1900 || tahun > time.Now().Year() {
						TampilkanError(errors.New("Tahun produksi tidak valid. Masukkan antara 1900 sampai " + strconv.Itoa(time.Now().Year()) + "."))
					}
				}

				fmt.Println("\nPilih pemilik:")
				TampilkanSemuaPemilik(jumlahPemilik)

				idPemilik = 0
				for CariIndeksPemilikSequential(idPemilik, jumlahPemilik) == -1 {
					idPemilik = InputAngka("\n? Masukkan ID Pemilik\n> ")
					if CariIndeksPemilikSequential(idPemilik, jumlahPemilik) == -1 {
						TampilkanError(errors.New("ID pemilik tidak valid."))
					}
				}

				TambahKendaraan(platNomor, merk, tipe, tahun, idPemilik, jumlahKendaraan)
				SimpanDataKendaraan(jumlahKendaraan)
			}

		case 3:
			fmt.Println("\n============== EDIT KENDARAAN ==============")
			if *jumlahKendaraan == 0 {
				fmt.Println("\nINFO: Belum ada data kendaraan.")
			} else {
				TampilkanSemuaKendaraan(jumlahKendaraan, jumlahPemilik)

				indeksKendaraan = 0
				for indeksKendaraan < 1 || indeksKendaraan > *jumlahKendaraan {
					indexInput = InputText("\n? Nomor kendaraan yang akan diedit\n> ")
					indeksKendaraan, err = strconv.Atoi(indexInput)
					if err != nil || indeksKendaraan < 1 || indeksKendaraan > *jumlahKendaraan {
						TampilkanError(errors.New("Nomor tidak valid."))
						indeksKendaraan = 0
					}
				}

				idx = indeksKendaraan - 1
				fmt.Println("\nINFO: Tekan Enter jika tidak ingin mengubah field tersebut.")

				platNomor = InputText("\n? Plat Nomor baru\n> ")
				if platNomor == "" {
					platNomor = DaftarKendaraan[idx].PlatNomor
				} else {
					platNomor = strings.ToUpper(platNomor)
				}

				merk = InputText("\n? Merk Kendaraan baru\n> ")
				if merk == "" {
					merk = DaftarKendaraan[idx].MerkKendaraan
				}

				tipe = InputText("\n? Tipe Kendaraan baru\n> ")
				if tipe == "" {
					tipe = DaftarKendaraan[idx].TipeKendaraan
				}

				tahunInput = InputText("\n? Tahun Produksi baru\n> ")
				if tahunInput == "" {
					tahun = DaftarKendaraan[idx].TahunProduksi
				} else {
					tahun, err = strconv.Atoi(tahunInput)
					for err != nil || tahun < 1900 || tahun > time.Now().Year() {
						TampilkanError(errors.New("Tahun produksi tidak valid."))
						tahunInput = InputText("\n? Tahun Produksi baru\n> ")
						tahun, err = strconv.Atoi(tahunInput)
					}
				}

				fmt.Println("\n--- Pilih Pemilik Baru (tekan Enter untuk tetap) ---")
				TampilkanSemuaPemilik(jumlahPemilik)

				idInput = InputText("\n? ID Pemilik baru\n> ")
				if idInput == "" {
					idPemilik = DaftarKendaraan[idx].IDPemilik
				} else {
					idPemilik, err = strconv.Atoi(idInput)
					for err != nil || CariIndeksPemilikSequential(idPemilik, jumlahPemilik) == -1 {
						TampilkanError(errors.New("ID pemilik tidak valid."))
						idInput = InputText("\n? ID Pemilik baru\n> ")
						if idInput == "" {
							idPemilik = DaftarKendaraan[idx].IDPemilik
							err = nil
						} else {
							idPemilik, err = strconv.Atoi(idInput)
						}
					}
				}

				EditKendaraan(indeksKendaraan, platNomor, merk, tipe, tahun, idPemilik)
				SimpanDataKendaraan(jumlahKendaraan)
			}

		case 4:
			fmt.Println("\n============== HAPUS KENDARAAN ==============")
			if *jumlahKendaraan == 0 {
				fmt.Println("\nINFO: Belum ada data kendaraan.")
			} else {
				TampilkanSemuaKendaraan(jumlahKendaraan, jumlahPemilik)

				for indeksKendaraan < 1 || indeksKendaraan > *jumlahKendaraan {
					indexInput = InputText("\n? Nomor kendaraan yang akan dihapus\n> ")
					indeksKendaraan, err = strconv.Atoi(indexInput)
					if err != nil || indeksKendaraan < 1 || indeksKendaraan > *jumlahKendaraan {
						TampilkanError(errors.New("Nomor tidak valid."))
						indeksKendaraan = 0
					}
				}

				konfirmasi = InputText(fmt.Sprintf("\n? Yakin hapus kendaraan %s? (y/n)\n> ",
					DaftarKendaraan[indeksKendaraan-1].PlatNomor))
				if strings.ToLower(konfirmasi) == "y" {
					HapusKendaraan(indeksKendaraan, jumlahKendaraan)
					SimpanDataKendaraan(jumlahKendaraan)
				} else {
					fmt.Println("\nINFO: Penghapusan dibatalkan.")
				}
			}
		case 0:
			return
		default:
			TampilkanError(errors.New("Pilihan tidak valid."))
		}
	}
}

func KelolaPemilik(jumlahKendaraan *int, jumlahPemilik *int) {
	var pilihan, id, indeks int
	var nama, telepon, alamat, konfirmasi string

	fmt.Println("\n\n============= KELOLA DATA PEMILIK =============")

	for {
		fmt.Println("\nMenu kelola:")
		fmt.Println("1. Lihat semua pemilik")
		fmt.Println("2. Tambah pemilik")
		fmt.Println("3. Edit pemilik")
		fmt.Println("4. Hapus pemilik")
		fmt.Println("0. Kembali ke menu utama")

		pilihan = InputAngka("\n? Pilih menu\n> ")

		switch pilihan {
		case 1:
			fmt.Println("\n============ DAFTAR SEMUA PEMILIK ============")
			if *jumlahPemilik == 0 {
				fmt.Println("\nINFO: Belum ada data pemilik.")
			} else {
				TampilkanSemuaPemilik(jumlahPemilik)
			}
		case 2:
			fmt.Println("\n============== TAMBAH PEMILIK ==============")
			nama = InputTeksWajib("\n? Nama Pemilik\n> ", "Nama pemilik tidak boleh kosong.")
			telepon = InputTeksWajib("\n? Nomor Telepon\n> ", "Nomor telepon tidak boleh kosong.")
			alamat = InputTeksWajib("\n? Alamat\n> ", "Alamat tidak boleh kosong.")
			TambahPemilik(nama, telepon, alamat, jumlahPemilik)
			SimpanDataPemilik(jumlahPemilik)
		case 3:
			fmt.Println("\n============== EDIT PEMILIK ==============")
			if *jumlahPemilik == 0 {
				fmt.Println("\nINFO: Belum ada data pemilik.")
			} else {
				TampilkanSemuaPemilik(jumlahPemilik)
				id = InputAngka("\n? ID Pemilik yang akan diedit\n> ")
				indeks = CariIndeksPemilikSequential(id, jumlahPemilik)

				if indeks == -1 {
					TampilkanError(errors.New("Pemilik tidak ditemukan."))
				} else {
					fmt.Println("\nINFO: Tekan Enter jika tidak ingin mengubah field tersebut.")

					nama = InputText("\n? Nama Pemilik baru\n> ")
					if nama == "" {
						nama = DaftarPemilik[indeks].NamaPemilik
					}
					telepon = InputText("\n? Nomor Telepon baru\n> ")
					if telepon == "" {
						telepon = DaftarPemilik[indeks].NomorTelepon
					}
					alamat = InputText("\n? Alamat baru\n> ")
					if alamat == "" {
						alamat = DaftarPemilik[indeks].Alamat
					}

					EditPemilik(id, nama, telepon, alamat, jumlahPemilik)
					SimpanDataPemilik(jumlahPemilik)
				}
			}
		case 4:
			fmt.Println("\n============== HAPUS PEMILIK ==============")
			if *jumlahPemilik == 0 {
				fmt.Println("\nINFO: Belum ada data pemilik.")
			} else {
				TampilkanSemuaPemilik(jumlahPemilik)
				id = InputAngka("\n? ID Pemilik yang akan dihapus\n> ")
				konfirmasi = InputText(fmt.Sprintf("\n? Yakin hapus pemilik ID %d? (y/n)\n> ", id))
				if strings.ToLower(konfirmasi) == "y" {
					HapusPemilik(id, jumlahKendaraan, jumlahPemilik)
					SimpanDataPemilik(jumlahPemilik)
				} else {
					fmt.Println("\nINFO: Hapus pemilik dibatalkan.")
				}
			}
		case 0:
			return
		default:
			TampilkanError(errors.New("Pilihan tidak valid."))
		}
	}
}

func KelolaServis(jumlahKendaraan *int, jumlahServis *int) {
	var pilihan, biaya, indeksServis int
	var platNomor, jenisKerusakan, detailServis, indexInput, konfirmasi string
	var err error

	fmt.Println("\n\n============= KELOLA RIWAYAT SERVIS =============")

	for {
		fmt.Println("\nMenu kelola:")
		fmt.Println("1. Lihat semua riwayat servis")
		fmt.Println("2. Tambah riwayat servis")
		fmt.Println("3. Hapus riwayat servis")
		fmt.Println("0. Kembali ke menu utama")

		pilihan = InputAngka("\n? Pilih menu\n> ")

		switch pilihan {
		case 1:
			fmt.Println("\n========== DAFTAR SEMUA RIWAYAT SERVIS ==========")
			if *jumlahServis == 0 {
				fmt.Println("\nINFO: Belum ada riwayat servis.")
			} else {
				TampilkanSemuaServis(jumlahServis)
			}
		case 2:
			fmt.Println("\n============ TAMBAH RIWAYAT SERVIS ============")
			if *jumlahKendaraan == 0 {
				fmt.Println("\nERROR: Belum ada data kendaraan.")
			} else {
				platNomor = ""
				for CariIndeksKendaraanSequential(platNomor, jumlahKendaraan) == -1 {
					platNomor = InputTeksWajib("\n? Plat Nomor Kendaraan\n> ", "Plat nomor tidak boleh kosong.")
					if CariIndeksKendaraanSequential(platNomor, jumlahKendaraan) == -1 {
						TampilkanError(errors.New("Kendaraan dengan plat nomor tersebut tidak ditemukan."))
					}
				}

				jenisKerusakan = InputTeksWajib("\n? Jenis Kerusakan (contoh: Mesin, Rem, Kelistrikan)\n> ", "Jenis kerusakan tidak boleh kosong.")
				detailServis = InputTeksWajib("\n? Detail Servis\n> ", "Detail servis tidak boleh kosong.")

				biaya = -1
				for biaya < 0 {
					biaya = InputAngka("\n? Biaya Servis (Rp)\n> ")
					if biaya < 0 {
						TampilkanError(errors.New("Biaya servis tidak boleh negatif."))
					}
				}

				TambahServis(platNomor, jenisKerusakan, detailServis, biaya, jumlahServis)
				SimpanDataServis(jumlahServis)
			}
		case 3:
			fmt.Println("\n============ HAPUS RIWAYAT SERVIS ============")
			if *jumlahServis == 0 {
				fmt.Println("\nINFO: Belum ada riwayat servis.")
			} else {
				TampilkanSemuaServis(jumlahServis)

				for indeksServis < 1 || indeksServis > *jumlahServis {
					indexInput = InputText("\n? Nomor riwayat servis yang akan dihapus\n> ")
					indeksServis, err = strconv.Atoi(indexInput)
					if err != nil || indeksServis < 1 || indeksServis > *jumlahServis {
						TampilkanError(errors.New("Nomor tidak valid."))
						indeksServis = 0
					}
				}

				konfirmasi = InputText(fmt.Sprintf("\n? Yakin hapus riwayat servis ID %d? (y/n)\n> ",
					DaftarServis[indeksServis-1].IDServis))
				if strings.ToLower(konfirmasi) == "y" {
					HapusServis(indeksServis, jumlahServis)
					SimpanDataServis(jumlahServis)
				} else {
					fmt.Println("\nINFO: Penghapusan dibatalkan.")
				}
			}
		case 0:
			return
		default:
			TampilkanError(errors.New("Pilihan tidak valid."))
		}
	}
}

func MenuCariKendaraan(jumlahKendaraan *int, jumlahPemilik *int) {
	var pilihan int
	var keyword string
	var selesai bool

	fmt.Println("\n\n============= CARI DATA KENDARAAN =============")

	selesai = false
	for !selesai {
		fmt.Println("\nMetode pencarian:")
		fmt.Println("1. Sequential Search (pencarian berdasarkan sebagian plat nomor)")
		fmt.Println("2. Binary Search (pencarian berdasarkan plat nomor tepat)")
		fmt.Println("0. Kembali ke menu utama")

		pilihan = InputAngka("\n? Pilih metode\n> ")

		switch pilihan {
		case 1:
			fmt.Println("\n====== CARI PLAT NOMOR (SEQUENTIAL SEARCH) ======")
			if *jumlahKendaraan == 0 {
				fmt.Println("\nINFO: Belum ada data kendaraan.")
			} else {
				TampilkanSemuaKendaraan(jumlahKendaraan, jumlahPemilik)
				keyword = InputText("\n? Masukkan plat nomor yang dicari\n> ")
				CariKendaraanSequential(keyword, jumlahKendaraan, jumlahPemilik)
			}

		case 2:
			fmt.Println("\n======== CARI PLAT NOMOR (BINARY SEARCH) ========")
			if *jumlahKendaraan == 0 {
				fmt.Println("\nINFO: Belum ada data kendaraan.")
			} else {
				TampilkanSemuaKendaraan(jumlahKendaraan, jumlahPemilik)
				keyword = InputText("\n? Masukkan plat nomor yang dicari (harus persis)\n> ")
				CariKendaraanBinary(keyword, jumlahKendaraan, jumlahPemilik)
			}
		case 0:
			return
		default:
			TampilkanError(errors.New("Pilihan tidak valid."))
		}
	}
}

func MenuUrutkanKendaraan(jumlahKendaraan *int, jumlahPemilik *int, jumlahServis *int) {
	var pilihan int

	fmt.Println("\n\n============= URUTKAN DATA KENDARAAN =============")

	for {
		fmt.Println("\nMenu pengurutan:")
		fmt.Println("1. Tahun Produksi - Ascending (Selection Sort)")
		fmt.Println("2. Tanggal Servis Terakhir - Descending (Insertion Sort)")
		fmt.Println("0. Kembali ke menu utama")

		pilihan = InputAngka("\n? Pilih menu\n> ")

		switch pilihan {
		case 1:
			fmt.Println("\n==== URUTKAN TAHUN PRODUKSI ASCENDING ====")
			if *jumlahKendaraan == 0 {
				fmt.Println("\nINFO: Belum ada data kendaraan.")
			} else {
				UrutkanKendaraanByTahunSelectionSort(jumlahKendaraan)
				TampilkanSemuaKendaraan(jumlahKendaraan, jumlahPemilik)
			}
		case 2:
			fmt.Println("\n==== URUTKAN TANGGAL SERVIS (DESCENDING) ====")
			if *jumlahKendaraan == 0 {
				fmt.Println("\nINFO: Belum ada data kendaraan.")
			} else {
				UrutkanKendaraanByTanggalServisInsertionSort(jumlahKendaraan, jumlahServis)
				TampilkanSemuaKendaraan(jumlahKendaraan, jumlahPemilik)
			}
		case 0:
			return
		default:
			TampilkanError(errors.New("Pilihan tidak valid."))
		}
	}
}

func main() {
	var jumlahKendaraan, jumlahPemilik, jumlahServis, pilihan int

	fmt.Println("\n==== AUTOCARE - MANAJEMEN SERVIS KENDARAAN ====")

	MuatDataKendaraan(&jumlahKendaraan)
	MuatDataPemilik(&jumlahPemilik)
	MuatDataServis(&jumlahServis)

	for {
		fmt.Println("\n\nMenu utama:")
		fmt.Println("1. Kelola Data Kendaraan")
		fmt.Println("2. Kelola Data Pemilik")
		fmt.Println("3. Kelola Riwayat Servis")
		fmt.Println("4. Cari Kendaraan (Sequential & Binary Search)")
		fmt.Println("5. Urutkan Kendaraan (Selection & Insertion Sort)")
		fmt.Println("6. Statistik Servis Kendaraan")
		fmt.Println("0. Keluar")

		pilihan = InputAngka("\n? Pilih menu\n> ")

		switch pilihan {
		case 1:
			KelolaKendaraan(&jumlahKendaraan, &jumlahPemilik)
		case 2:
			KelolaPemilik(&jumlahKendaraan, &jumlahPemilik)
		case 3:
			KelolaServis(&jumlahKendaraan, &jumlahServis)
		case 4:
			MenuCariKendaraan(&jumlahKendaraan, &jumlahPemilik)
		case 5:
			MenuUrutkanKendaraan(&jumlahKendaraan, &jumlahPemilik, &jumlahServis)
		case 6:
			TampilkanStatistik(&jumlahServis)
		case 0:
			fmt.Println("\n================== SAMPAI JUMPA! ==================")
			os.Exit(0)
		default:
			TampilkanError(errors.New("Pilihan tidak valid."))
		}
	}
}
