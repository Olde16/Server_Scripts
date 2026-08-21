package main

import (
	"encoding/json"
	"net/http"
)

func medalsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMedals(w, r)

	case http.MethodPost:
		createMedal(w, r)

	case http.MethodPut:
		updateMedal(w, r)

	case http.MethodDelete:
		deleteMedal(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getMedals(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT medaillen_id, art_medaille
		FROM medals
		ORDER BY medaillen_id
	`)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Medal struct {
		MedaillenID uint16 `json:"medaillen_id"`
		ArtMedaille string `json:"art_medaille"`
	}

	var medals []Medal

	for rows.Next() {
		var medal Medal

		err := rows.Scan(
			&medal.MedaillenID,
			&medal.ArtMedaille,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		medals = append(medals, medal)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medals)
}

func createMedal(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Medal struct {
		MedaillenID uint16 `json:"medaillen_id"`
		ArtMedaille string `json:"art_medaille"`
	}

	var medal Medal

	err = json.NewDecoder(r.Body).Decode(&medal)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		INSERT INTO medals (medaillen_id, art_medaille)
		VALUES (?, ?)
	`,
		medal.MedaillenID,
		medal.ArtMedaille,
	)

	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"medal":  medal,
	})
}

func updateMedal(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Medal struct {
		MedaillenID uint16 `json:"medaillen_id"`
		ArtMedaille string `json:"art_medaille"`
	}

	var medal Medal

	err = json.NewDecoder(r.Body).Decode(&medal)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE medals
		SET art_medaille = ?
		WHERE medaillen_id = ?
	`,
		medal.ArtMedaille,
		medal.MedaillenID,
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
		http.Error(w, "Medal not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"medal":  medal,
	})
}

func deleteMedal(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Medal struct {
		MedaillenID uint16 `json:"medaillen_id"`
	}

	var medal Medal

	err = json.NewDecoder(r.Body).Decode(&medal)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		DELETE FROM medals
		WHERE medaillen_id = ?
	`, medal.MedaillenID)

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
		http.Error(w, "Medal not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"message":      "Medal deleted",
		"medaillen_id": medal.MedaillenID,
	})
}
