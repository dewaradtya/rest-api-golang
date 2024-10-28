package handlers

import (
	"belajar-golang/database"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var tokenBlacklist = make(map[string]struct{})

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var jwtKey = []byte("secret_key")

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Handler untuk registrasi
func Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	// Simpan user ke database
	_, err = database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", registerRequest.Name, registerRequest.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Handler untuk login
func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user User
	err := database.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", loginRequest.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Buat token JWT
	expirationTime := time.Now().Add(60 * time.Minute) // Token berlaku selama 60 menit
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Ambil token dari header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Tambahkan token ke dalam blacklist
	tokenBlacklist[tokenString] = struct{}{}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Periksa apakah token ada di blacklist
		if _, found := tokenBlacklist[tokenString]; found {
			http.Error(w, "Token has been logged out", http.StatusUnauthorized)
			return
		}

		// Parse token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token valid, lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

// Contoh handler yang hanya bisa diakses jika sudah login
func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "This is a protected route"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
