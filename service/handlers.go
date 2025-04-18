package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Logged in as: " + email))
}

func addSubscription(w http.ResponseWriter, r *http.Request) {
	var sub Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	sub.ID = uuid.NewString()

	startDate, err := time.Parse("2006-01-02", sub.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		return
	}

	var endDate sql.NullTime
	if sub.EndDate != "" {
		t, err := time.Parse("2006-01-02", sub.EndDate)
		if err != nil {
			http.Error(w, "Invalid end_date format", http.StatusBadRequest)
			return
		}
		endDate = sql.NullTime{Time: t, Valid: true}
	}

	_, err = DB.Exec(`
        INSERT INTO subscriptions 
        (id, email, type, name, start_date, end_date, currency, amount) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		sub.ID, sub.Email, sub.Type, sub.Name, startDate, endDate, sub.Currency, sub.Amount,
	)
	if err != nil {
		log.Println("Insert error:", err)
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func deleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Subscription ID is required", http.StatusBadRequest)
		return
	}

	_, err := DB.Exec(`ALTER TABLE subscriptions DELETE WHERE id = ?`, id)
	if err != nil {
		log.Println("Delete error:", err)
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Deleted"))
}

func listSubscriptions(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	rows, err := DB.Query(`
        SELECT id, email, type, name, start_date, end_date, currency, amount, created_at 
        FROM subscriptions 
        WHERE email = ?
        ORDER BY start_date DESC`, email)
	if err != nil {
		log.Println("List error:", err)
		http.Error(w, "Failed to query subscriptions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subscriptions []Subscription
	for rows.Next() {
		var sub Subscription
		var startDate, endDate sql.NullTime
		err := rows.Scan(&sub.ID, &sub.Email, &sub.Type, &sub.Name, &startDate, &endDate, &sub.Currency, &sub.Amount, &sub.CreatedAt)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		sub.StartDate = startDate.Time.Format("2006-01-02")
		if endDate.Valid {
			sub.EndDate = endDate.Time.Format("2006-01-02")
		}
		subscriptions = append(subscriptions, sub)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}
