package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

type Input struct {
	Jumlah    float64 `json:"jumlah"`
	Deskripsi string  `json:"deskripsi"`
	Kategori  string  `json:"kategori"`
	Tanggal   string  `json:"tanggal"`
}
type data struct {
	gorm.Model
	Jumlah    float64   `json:"jumlah"`
	Deskripsi string    `json:"deskripsi"`
	Kategori  string    `json:"kategori"`
	Tanggal   time.Time `json:"tanggal"`
}

type getData struct {
	ID        int       `json:"id"`
	Jumlah    float64   `json:"jumlah"`
	Deskripsi string    `json:"deskripsi"`
	Kategori  string    `json:"kategori"`
	Tanggal   time.Time `json:"tanggal"`
}

type Total struct {
	Total float64 `json:"total"`
}

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

	hasil := map[string]any{
		"pesan": "data berhasil dibuat",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}

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
	hasil := map[string]any{
		"pesan": "data berhasil diupdate",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}

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
	hasil := map[string]any{
		"pesan": "data berhasil dihapus",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}

func main() {
	dsn := "host=localhost user=postgres password=test162534 dbname=expense_tracker port=5432"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Gagal konek ke database", err)
		return
	}
	http.HandleFunc("/tambah", handlerInput)
	http.HandleFunc("/pengeluaran", handlerPengeluaran)
	http.HandleFunc("/total", handlerTotal)
	http.HandleFunc("/filter", handlerFilter)
	http.HandleFunc("/update", handlerUpdate)
	http.HandleFunc("/hapus", handlerHapus)
	db.AutoMigrate(&data{})

	fmt.Println("Server jalan di http://localhost:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
