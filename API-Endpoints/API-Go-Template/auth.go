package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func generateToken() (string, string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", "", err
	}

	token := hex.EncodeToString(bytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	return token, tokenHash, nil
}

func authenticateRequest(r *http.Request) (int, error) {
	header := r.Header.Get("Authorization")

	if !strings.HasPrefix(header, "User ") {
		return 0, errors.New("missing authorization token")
	}

	token := strings.TrimPrefix(header, "User ")

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	var userID int

	err := db.QueryRow(`
		SELECT user_id
		FROM api_tokens
		WHERE token_hash = ?
		  AND expires_at > NOW()
	`, tokenHash).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("invalid or expired token")
		}

		return 0, err
	}

	return userID, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var request LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var (
		userID       int
		passwordHash string
	)

	err = db.QueryRow(`
		SELECT id, password_hash
		FROM authorized_users
		WHERE username = ?
	`, request.Username).Scan(&userID, &passwordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(request.Password),
	)

	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, tokenHash, err := generateToken()

	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	expires := time.Now().Add(24 * time.Hour)

	_, err = db.Exec(`
		INSERT INTO api_tokens (user_id, token_hash, expires_at)
		VALUES (?, ?, ?)
	`, userID, tokenHash, expires)

	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"token":       token,
		"expires_at":  expires,
	}

	json.NewEncoder(w).Encode(response)
}
