package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = DB.Exec(`INSERT INTO users (email, password_hash) VALUES (?, ?)`, user.Email, string(hash))
	if err != nil {
		log.Println("Signup insert error:", err)
		http.Error(w, "Signup failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Signup successful"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	row := DB.QueryRow(`SELECT password_hash FROM users WHERE email = ?`, creds.Email)

	var hashedPassword string
	if err := row.Scan(&hashedPassword); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(creds.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func addSubscription(w http.ResponseWriter, r *http.Request) {
	email, err := validateJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var sub Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	sub.ID = uuid.NewString()
	sub.Email = email // enforce ownership

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
	email, err := validateJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Subscription ID is required", http.StatusBadRequest)
		return
	}

	_, err = DB.Exec(`ALTER TABLE subscriptions DELETE WHERE id = ? AND email = ?`, id, email)
	if err != nil {
		log.Println("Delete error:", err)
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Deleted"))
}

func listSubscriptions(w http.ResponseWriter, r *http.Request) {
	email, err := validateJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
