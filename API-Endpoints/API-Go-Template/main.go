package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Success: true,
		Message: "Olympics API is running!",
	}

	json.NewEncoder(w).Encode(response)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var count int

	err := db.QueryRow(
		"SELECT COUNT(*) FROM authorized_users",
	).Scan(&count)

	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"users":   count,
	}

	json.NewEncoder(w).Encode(response)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "You are authenticated!",
		"user_id": userID,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
    initDatabase()

    http.HandleFunc("/api/hello", helloHandler)
    http.HandleFunc("/api/users", usersHandler)
    http.HandleFunc("/api/login", loginHandler)
    http.HandleFunc("/api/protected", protectedHandler)
    http.HandleFunc("/api/country", countryHandler)
    http.HandleFunc("/api/athletes", athletesHandler)
    http.HandleFunc("/api/medals", medalsHandler)
    http.HandleFunc("/api/sports", sportsHandler)
    http.HandleFunc("/api/medals-athletes-sports", medalsAthletesSportsHandler)

    log.Println("Olympics API listening on :8080")

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
