package main

import (
	"belajar-golang/database"
	"belajar-golang/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Koneksi ke database
	database.Connect()

	// Buat router
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/login", handlers.Login).Methods("POST")
    r.HandleFunc("/api/v1/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/v1/logout", handlers.Logout).Methods("POST")

	// Endpoint untuk CRUD post (dengan autentikasi)
	r.Handle("/api/v1/posts", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetAllPosts))).Methods("GET")
	r.Handle("/api/v1/posts", handlers.AuthMiddleware(http.HandlerFunc(handlers.CreatePost))).Methods("POST")
	r.Handle("/api/v1/posts/{id}", handlers.AuthMiddleware(http.HandlerFunc(handlers.UpdatePost))).Methods("PUT")
	r.Handle("/api/v1/posts/{id}", handlers.AuthMiddleware(http.HandlerFunc(handlers.DeletePost))).Methods("DELETE")

	// Jalankan server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
