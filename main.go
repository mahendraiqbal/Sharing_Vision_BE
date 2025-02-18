package main

import (
	"log"
	"net/http"

	"article-api/handlers"
	"article-api/db"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Inisialisasi koneksi database
	db.InitDB("root:password@tcp(127.0.0.1:3306)/article")
	defer db.CloseDB()

	// Inisialisasi router
	r := mux.NewRouter()

	// Routing
	r.HandleFunc("/article", handlers.CreateArticle).Methods("POST") // POST /article
	r.HandleFunc("/articles", handlers.GetArticles).Methods("GET")  // GET /articles?limit=<limit>&offset=<offset>
	r.HandleFunc("/article/{id}", handlers.GetArticleByID).Methods("GET") // GET /article/<id>
	r.HandleFunc("/article/{id}", handlers.UpdateArticle).Methods("PUT") // PUT /article/<id>
	r.HandleFunc("/article/{id}", handlers.DeleteArticle).Methods("DELETE") // DELETE /article/<id>

	// Konfigurasi CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Izinkan akses dari Next.js
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true, // Aktifkan debug untuk melihat log CORS
	})

	// Gunakan middleware CORS
	handler := c.Handler(r)

	// Jalankan server
	log.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
