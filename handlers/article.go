package handlers

import (
	"article-api/models"
	"article-api/repositories"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"errors"

	"github.com/gorilla/mux"
)

// Validasi data artikel
func validateArticle(article models.Article) error {
	// Validasi Title
	if article.Title == "" {
		return errors.New("Title harus diisi")
	}
	if len(article.Title) < 20 {
		return errors.New("Title minimal 20 karakter")
	}

	// Validasi Content
	if article.Content == "" {
		return errors.New("Content harus diisi")
	}
	if len(article.Content) < 200 {
		return errors.New("Content minimal 200 karakter")
	}

	// Validasi Category
	if article.Category == "" {
		return errors.New("Category harus diisi")
	}
	if len(article.Category) < 3 {
		return errors.New("Category minimal 3 karakter")
	}

	// Validasi Status
	validStatus := map[string]bool{
		"publish": true,
		"draft":   true,
		"thrash":  true,
	}
	if article.Status == "" {
		return errors.New("Status harus diisi")
	}
	if !validStatus[strings.ToLower(article.Status)] {
		return errors.New("Status harus salah satu dari: publish, draft, thrash")
	}

	return nil
}

// CreateArticle menangani request POST /article
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method tidak diizinkan:", r.Method)
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body ke struct Article
	var article models.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		log.Println("Gagal membaca request body:", err)
		http.Error(w, "Gagal membaca request body", http.StatusBadRequest)
		return
	}

	// Validasi data artikel
	if err := validateArticle(article); err != nil {
		log.Println("Validasi gagal:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simpan artikel ke database
	id, err := repositories.CreateArticle(article)
	if err != nil {
		log.Println("Gagal menyimpan artikel:", err)
		http.Error(w, "Gagal menyimpan artikel", http.StatusInternalServerError)
		return
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"id":      id,
		"message": "Artikel berhasil dibuat",
	}
	json.NewEncoder(w).Encode(response)
}

// GetArticles menangani request GET /articles?limit=<limit>&offset=<offset>
func GetArticles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method tidak diizinkan:", r.Method)
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// Ambil nilai limit dan offset dari URL
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Konversi limit dan offset ke integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println("Limit tidak valid:", err)
		http.Error(w, "Limit tidak valid", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Println("Offset tidak valid:", err)
		http.Error(w, "Offset tidak valid", http.StatusBadRequest)
		return
	}

	// Ambil daftar artikel dari database
	articles, err := repositories.GetArticles(limit, offset)
	if err != nil {
		log.Println("Gagal mengambil artikel:", err)
		http.Error(w, "Gagal mengambil artikel", http.StatusInternalServerError)
		return
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// GetArticleByID menangani request GET /article/<id>
func GetArticleByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method tidak diizinkan:", r.Method)
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// Ambil nilai id dari URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Konversi id ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("ID tidak valid:", err)
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Ambil detail artikel dari database
	article, err := repositories.GetArticleByID(id)
	if err != nil {
		log.Println("Gagal mengambil artikel:", err)
		http.Error(w, "Gagal mengambil artikel", http.StatusInternalServerError)
		return
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// UpdateArticle menangani request PUT /article/<id>
func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Println("Method tidak diizinkan:", r.Method)
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// Ambil nilai id dari URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Konversi id ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("ID tidak valid:", err)
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Decode request body ke struct Article
	var article models.Article
	err = json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		log.Println("Gagal membaca request body:", err)
		http.Error(w, "Gagal membaca request body", http.StatusBadRequest)
		return
	}

	// Validasi data artikel
	if err := validateArticle(article); err != nil {
		log.Println("Validasi gagal:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update artikel di database
	err = repositories.UpdateArticle(id, article)
	if err != nil {
		log.Println("Gagal mengubah artikel:", err)
		http.Error(w, "Gagal mengubah artikel", http.StatusInternalServerError)
		return
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Artikel berhasil diubah",
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteArticle menangani request DELETE /article/<id>
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		log.Println("Method tidak diizinkan:", r.Method)
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// Ambil nilai id dari URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Konversi id ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("ID tidak valid:", err)
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Hapus artikel dari database
	err = repositories.DeleteArticle(id)
	if err != nil {
		log.Println("Gagal menghapus artikel:", err)
		http.Error(w, "Gagal menghapus artikel", http.StatusInternalServerError)
		return
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Artikel berhasil dihapus",
	}
	json.NewEncoder(w).Encode(response)
}
