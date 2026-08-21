package main

import (
	"encoding/json"
	"net/http"
)

func countryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCountries(w, r)

	case http.MethodPost:
		createCountry(w, r)

	case http.MethodPut:
		updateCountry(w, r)

	case http.MethodDelete:
		deleteCountry(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getCountries(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT land_id, land_code, land_name
		FROM country
		ORDER BY land_id
	`)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Country struct {
		LandID   uint16 `json:"land_id"`
		LandCode string `json:"land_code"`
		LandName string `json:"land_name"`
	}

	var countries []Country

	for rows.Next() {
		var c Country

		err := rows.Scan(
			&c.LandID,
			&c.LandCode,
			&c.LandName,
		)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		countries = append(countries, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func createCountry(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Country struct {
		LandID   uint16 `json:"land_id"`
		LandCode string `json:"land_code"`
		LandName string `json:"land_name"`
	}

	var country Country

	err = json.NewDecoder(r.Body).Decode(&country)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		INSERT INTO country (land_id, land_code, land_name)
		VALUES (?, ?, ?)
	`,
		country.LandID,
		country.LandCode,
		country.LandName,
	)

	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"country": country,
	})
}

func updateCountry(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Country struct {
		LandID   uint16 `json:"land_id"`
		LandCode string `json:"land_code"`
		LandName string `json:"land_name"`
	}

	var country Country

	err = json.NewDecoder(r.Body).Decode(&country)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE country
		SET land_code = ?, land_name = ?
		WHERE land_id = ?
	`,
		country.LandCode,
		country.LandName,
		country.LandID,
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
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"country": country,
	})
}

func deleteCountry(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Country struct {
		LandID uint16 `json:"land_id"`
	}

	var country Country

	err = json.NewDecoder(r.Body).Decode(&country)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		DELETE FROM country
		WHERE land_id = ?
	`, country.LandID)

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
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Country deleted",
		"land_id": country.LandID,
	})
}


