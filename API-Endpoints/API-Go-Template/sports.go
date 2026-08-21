package main

import (
	"encoding/json"
	"net/http"
)

func sportsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSports(w, r)

	case http.MethodPost:
		createSport(w, r)

	case http.MethodPut:
		updateSport(w, r)

	case http.MethodDelete:
		deleteSport(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getSports(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT sportart_id, sportart
		FROM sports
		ORDER BY sportart_id
	`)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Sport struct {
		SportartID uint   `json:"sportart_id"`
		Sportart   string `json:"sportart"`
	}

	var sports []Sport

	for rows.Next() {
		var sport Sport

		err := rows.Scan(
			&sport.SportartID,
			&sport.Sportart,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		sports = append(sports, sport)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sports)
}

func createSport(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Sport struct {
		SportartID uint   `json:"sportart_id"`
		Sportart   string `json:"sportart"`
	}

	var sport Sport

	err = json.NewDecoder(r.Body).Decode(&sport)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		INSERT INTO sports (sportart_id, sportart)
		VALUES (?, ?)
	`,
		sport.SportartID,
		sport.Sportart,
	)

	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"sport":   sport,
	})
}

func updateSport(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Sport struct {
		SportartID uint   `json:"sportart_id"`
		Sportart   string `json:"sportart"`
	}

	var sport Sport

	err = json.NewDecoder(r.Body).Decode(&sport)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE sports
		SET sportart = ?
		WHERE sportart_id = ?
	`,
		sport.Sportart,
		sport.SportartID,
	)

	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Sport not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"sport":   sport,
	})
}

func deleteSport(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Sport struct {
		SportartID uint `json:"sportart_id"`
	}

	var sport Sport

	err = json.NewDecoder(r.Body).Decode(&sport)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		DELETE FROM sports
		WHERE sportart_id = ?
	`, sport.SportartID)

	if err != nil {
		http.Error(w, "Database delete failed", http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Sport not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"message":     "Sport deleted",
		"sportart_id": sport.SportartID,
	})
}
