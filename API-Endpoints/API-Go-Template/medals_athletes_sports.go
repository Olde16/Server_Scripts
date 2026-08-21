package main

import (
	"encoding/json"
	"net/http"
)

type MedalAthleteSport struct {
	LandID          uint16 `json:"land_id"`
	AthletenID      uint   `json:"athleten_id"`
	MedaillenID     uint16 `json:"medaillen_id"`
	SportartID      uint   `json:"sportart_id"`
	AnzahlMedaille  uint   `json:"anzahl_medaillen"`
}

func medalsAthletesSportsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMedalsAthletesSports(w, r)

	case http.MethodPost:
		createMedalsAthletesSports(w, r)

	case http.MethodPut:
		updateMedalsAthletesSports(w, r)

	case http.MethodDelete:
		deleteMedalsAthletesSports(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getMedalsAthletesSports(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			land_id,
			athleten_id,
			medaillen_id,
			sportart_id,
			anzahl_medaillen
		FROM medals_athletes_sports
	`

	var args []interface{}

	landID := r.URL.Query().Get("land_id")
	athleteID := r.URL.Query().Get("athleten_id")

	if landID != "" {
		query += " WHERE land_id = ?"
		args = append(args, landID)
	}

	if athleteID != "" {
		if landID != "" {
			query += " AND athleten_id = ?"
		} else {
			query += " WHERE athleten_id = ?"
		}

		args = append(args, athleteID)
	}

	query += `
		ORDER BY land_id, athleten_id, medaillen_id, sportart_id
	`

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []MedalAthleteSport

	for rows.Next() {
		var result MedalAthleteSport

		err := rows.Scan(
			&result.LandID,
			&result.AthletenID,
			&result.MedaillenID,
			&result.SportartID,
			&result.AnzahlMedaille,
		)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		results = append(results, result)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func createMedalsAthletesSports(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var data MedalAthleteSport

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		INSERT INTO medals_athletes_sports (
			land_id,
			athleten_id,
			medaillen_id,
			sportart_id,
			anzahl_medaillen
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		data.LandID,
		data.AthletenID,
		data.MedaillenID,
		data.SportartID,
		data.AnzahlMedaille,
	)

	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func updateMedalsAthletesSports(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var data MedalAthleteSport

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE medals_athletes_sports
		SET anzahl_medaillen = ?
		WHERE land_id = ?
		  AND athleten_id = ?
		  AND medaillen_id = ?
		  AND sportart_id = ?
	`,
		data.AnzahlMedaille,
		data.LandID,
		data.AthletenID,
		data.MedaillenID,
		data.SportartID,
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
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func deleteMedalsAthletesSports(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var data MedalAthleteSport

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		DELETE FROM medals_athletes_sports
		WHERE land_id = ?
		  AND athleten_id = ?
		  AND medaillen_id = ?
		  AND sportart_id = ?
	`,
		data.LandID,
		data.AthletenID,
		data.MedaillenID,
		data.SportartID,
	)

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
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Record deleted",
	})
}


