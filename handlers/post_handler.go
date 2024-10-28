package handlers

import (
	"belajar-golang/database"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Post struct {
	ID    int    `json:"id"`
	Title  string `json:"title"`
	Desciption string `json:"desciption"`
}

// mengambil semua pengguna
func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, desciption FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Desciption); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// menambahkan pengguna baru
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("INSERT INTO posts (title, desciption) VALUES (?, ?)", post.Title, post.Desciption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	post.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// memperbarui data pengguna berdasarkan ID
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	vars := mux.Vars(r)
	id := vars["id"]

	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update data pengguna berdasarkan ID
	_, err := database.DB.Exec("UPDATE posts SET title = ?, desciption = ? WHERE id = ?", post.Title, post.Desciption, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})
}

// menghapus pengguna berdasarkan ID
func DeletePost(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Hapus pengguna berdasarkan ID
	_, err := database.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}
