package main

import (
	"encoding/json"
	"net/http"
)

func athletesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAthletes(w, r)

	case http.MethodPost:
		createAthlete(w, r)

	case http.MethodPut:
		updateAthlete(w, r)

	case http.MethodDelete:
		deleteAthlete(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAthletes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT athleten_id, athleten_vorname, athleten_name
		FROM athletes
		ORDER BY athleten_id
	`)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Athlete struct {
		AthletenID      uint   `json:"athleten_id"`
		AthletenVorname string `json:"athleten_vorname"`
		AthletenName    string `json:"athleten_name"`
	}

	var athletes []Athlete

	for rows.Next() {
		var athlete Athlete

		err := rows.Scan(
			&athlete.AthletenID,
			&athlete.AthletenVorname,
			&athlete.AthletenName,
		)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		athletes = append(athletes, athlete)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(athletes)
}

func createAthlete(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Athlete struct {
		AthletenID      uint   `json:"athleten_id"`
		AthletenVorname string `json:"athleten_vorname"`
		AthletenName    string `json:"athleten_name"`
	}

	var athlete Athlete

	err = json.NewDecoder(r.Body).Decode(&athlete)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		INSERT INTO athletes (
			athleten_id,
			athleten_vorname,
			athleten_name
		)
		VALUES (?, ?, ?)
	`,
		athlete.AthletenID,
		athlete.AthletenVorname,
		athlete.AthletenName,
	)

	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"athlete": athlete,
	})
}

func updateAthlete(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Athlete struct {
		AthletenID      uint   `json:"athleten_id"`
		AthletenVorname string `json:"athleten_vorname"`
		AthletenName    string `json:"athleten_name"`
	}

	var athlete Athlete

	err = json.NewDecoder(r.Body).Decode(&athlete)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE athletes
		SET athleten_vorname = ?, athleten_name = ?
		WHERE athleten_id = ?
	`,
		athlete.AthletenVorname,
		athlete.AthletenName,
		athlete.AthletenID,
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
		http.Error(w, "Athlete not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"athlete": athlete,
	})
}

func deleteAthlete(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Athlete struct {
		AthletenID uint `json:"athleten_id"`
	}

	var athlete Athlete

	err = json.NewDecoder(r.Body).Decode(&athlete)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		DELETE FROM athletes
		WHERE athleten_id = ?
	`, athlete.AthletenID)

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
		http.Error(w, "Athlete not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "Athlete deleted",
		"athleten_id": athlete.AthletenID,
	})
}
