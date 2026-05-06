package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "expense-api/docs" // <-- GANTI "expense-api" sesuai module name di go.mod kamu

	httpSwagger "github.com/swaggo/http-swagger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type ErrorResponse struct {
	Error string `json:"error"`
}

type ResponPesan struct {
	Pesan string `json:"pesan"`
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

type Input struct {
	Jumlah    float64 `json:"jumlah" example:"50000"`
	Deskripsi string  `json:"deskripsi" example:"Makan siang"`
	Kategori  string  `json:"kategori" example:"makanan"`
	Tanggal   string  `json:"tanggal" example:"2025-01-15"`
}
type data struct {
	gorm.Model
	Jumlah    float64   `json:"jumlah"`
	Deskripsi string    `json:"deskripsi"`
	Kategori  string    `json:"kategori"`
	Tanggal   time.Time `json:"tanggal"`
}

type getData struct {
	ID        int       `json:"id" example:"1"`
	Jumlah    float64   `json:"jumlah" example:"50000"`
	Deskripsi string    `json:"deskripsi" example:"Makan siang"`
	Kategori  string    `json:"kategori" example:"makanan"`
	Tanggal   time.Time `json:"tanggal"`
}

type Total struct {
	Total float64 `json:"total" example:"250000"`
}

// handlerInput godoc
// @Summary      Tambah pengeluaran baru
// @Description  Menambahkan data pengeluaran baru ke database
// @Description  Format tanggal harus YYYY-MM-DD
// @Tags         Pengeluaran
// @Accept       json
// @Produce      json
// @Param        request  body      Input        true  "Data pengeluaran"
// @Success      200      {object}  ResponPesan
// @Failure      400      {object}  ErrorResponse
// @Failure      405      {object}  ErrorResponse
// @Router       /tambah [post]
func handlerInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, "method harus POST", 405)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, "format JSON tidak valid", 400)
		return
	}

	if input.Jumlah == 0 {
		sendError(w, "jumlah tidak valid", 400)
		return
	}
	if input.Deskripsi == "" {
		sendError(w, "mohon isi deskripsi", 400)
		return
	}
	if input.Kategori == "" {
		sendError(w, "mohon isi kategori", 400)
		return
	}
	if input.Tanggal == "" {
		sendError(w, "mohon isi tanggal", 400)
		return
	}

	tanggal, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		sendError(w, "format tanggal harus YYYY-MM-DD", 400)
		return
	}

	db.Create(&data{
		Jumlah:    input.Jumlah,
		Deskripsi: input.Deskripsi,
		Kategori:  input.Kategori,
		Tanggal:   tanggal,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponPesan{
		Pesan: "data berhasil dibuat",
	})
}

// handlerPengeluaran godoc
// @Summary      Lihat semua pengeluaran
// @Description  Menampilkan seluruh data pengeluaran dari database
// @Tags         Pengeluaran
// @Produce      json
// @Success      200  {array}   getData
// @Failure      400  {object}  ErrorResponse
// @Router       /pengeluaran [get]
func handlerPengeluaran(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, "method harus GET", 400)
		return
	}
	var datas []data
	db.Find(&datas)

	var hasil []getData
	for _, v := range datas {
		hasil = append(hasil, getData{
			ID:        int(v.ID),
			Jumlah:    v.Jumlah,
			Deskripsi: v.Deskripsi,
			Kategori:  v.Kategori,
			Tanggal:   v.Tanggal,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}

// handlerTotal godoc
// @Summary      Lihat total pengeluaran
// @Description  Menghitung dan menampilkan total seluruh pengeluaran
// @Tags         Pengeluaran
// @Produce      json
// @Success      200  {object}  Total
// @Failure      405  {object}  ErrorResponse
// @Router       /total [get]
func handlerTotal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, "method harus GET", 405)
		return
	}
	var total float64
	db.Model(&data{}).Select("SUM(jumlah)").Scan(&total)
	temp := Total{
		Total: total,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temp)
}

// handlerFilter godoc
// @Summary      Filter pengeluaran berdasarkan tanggal
// @Description  Menampilkan pengeluaran dalam rentang tanggal tertentu
// @Description  Format tanggal: YYYY-MM-DD
// @Tags         Pengeluaran
// @Produce      json
// @Param        dari    query   string  true  "Tanggal awal (YYYY-MM-DD)"
// @Param        sampai  query   string  true  "Tanggal akhir (YYYY-MM-DD)"
// @Success      200     {array}   getData
// @Failure      405     {object}  ErrorResponse
// @Router       /filter [get]
func handlerFilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, "method harus GET", 405)
		return
	}

	dari := r.URL.Query().Get("dari")
	sampai := r.URL.Query().Get("sampai")

	var data []data
	db.Where("tanggal BETWEEN ? AND ?", dari, sampai).Find(&data)

	var hasil []getData
	for _, v := range data {
		hasil = append(hasil, getData{
			ID:        int(v.ID),
			Jumlah:    v.Jumlah,
			Deskripsi: v.Deskripsi,
			Kategori:  v.Kategori,
			Tanggal:   v.Tanggal,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}

// handlerUpdate godoc
// @Summary      Update pengeluaran berdasarkan ID
// @Description  Mengubah data pengeluaran yang sudah ada
// @Description  Format tanggal harus YYYY-MM-DD
// @Tags         Pengeluaran
// @Accept       json
// @Produce      json
// @Param        id       query   int    true  "ID pengeluaran"
// @Param        request  body    Input  true  "Data pengeluaran baru"
// @Success      200      {object}  ResponPesan
// @Failure      400      {object}  ErrorResponse
// @Failure      405      {object}  ErrorResponse
// @Router       /update [put]
func handlerUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		sendError(w, "method harus PUT", 405)
		return
	}

	strId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		sendError(w, "format id tidak valid", 400)
		return
	}
	if id == 0 {
		sendError(w, "mohon masukkan id dengan benar", 400)
		return
	}

	var input Input
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, "format JSON tidak valid", 400)
		return
	}

	if input.Jumlah == 0 {
		sendError(w, "jumlah tidak valid", 400)
		return
	}
	if input.Deskripsi == "" {
		sendError(w, "mohon isi deskripsi", 400)
		return
	}
	if input.Kategori == "" {
		sendError(w, "mohon isi kategori", 400)
		return
	}
	if input.Tanggal == "" {
		sendError(w, "mohon isi tanggal", 400)
		return
	}

	tanggal, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		sendError(w, "format tanggal harus YYYY-MM-DD", 400)
		return
	}

	db.Model(&data{}).Where("id = ?", id).Updates(map[string]any{
		"jumlah":    input.Jumlah,
		"deskripsi": input.Deskripsi,
		"kategori":  input.Kategori,
		"tanggal":   tanggal,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponPesan{
		Pesan: "data berhasil diupdate",
	})
}

// handlerHapus godoc
// @Summary      Hapus pengeluaran berdasarkan ID
// @Description  Menghapus data pengeluaran dari database
// @Tags         Pengeluaran
// @Produce      json
// @Param        id  query  int  true  "ID pengeluaran"
// @Success      200  {object}  ResponPesan
// @Failure      400  {object}  ErrorResponse
// @Failure      405  {object}  ErrorResponse
// @Router       /hapus [delete]
func handlerHapus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		sendError(w, "method harus DELETE", 405)
		return
	}

	strId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		sendError(w, "ID tidak valid", 400)
		return
	}

	if id == 0 {
		sendError(w, "ID tidak ada/tidak terdaftar", 400)
		return
	}

	db.Delete(&data{}, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponPesan{
		Pesan: "data berhasil dihapus",
	})
}

// @title           Expense API
// @version         1.0
// @description     REST API untuk tracking pengeluaran harian
// @description     Fitur: CRUD pengeluaran, total, filter by tanggal

// @contact.name    Jason
// @contact.url     https://github.com/Tarquished

// @host            expense-api-production-4f40.up.railway.app
// @schemes 		https
// @BasePath        /

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=test162534 dbname=expense_tracker port=5432 sslmode=disable"
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Gagal konek ke database", err)
		return
	}

	db.AutoMigrate(&data{})

	http.HandleFunc("/tambah", handlerInput)
	http.HandleFunc("/pengeluaran", handlerPengeluaran)
	http.HandleFunc("/total", handlerTotal)
	http.HandleFunc("/filter", handlerFilter)
	http.HandleFunc("/update", handlerUpdate)
	http.HandleFunc("/hapus", handlerHapus)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server jalan di port", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
