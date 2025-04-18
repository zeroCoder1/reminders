package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/google/uuid"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Email required", http.StatusBadRequest)
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

    _, err := DB.Exec(`
        INSERT INTO subscriptions 
        (id, email, type, name, start_date, end_date) 
        VALUES (?, ?, ?, ?, ?, ?)`,
        sub.ID, sub.Email, sub.Type, sub.Name, sub.StartDate, sub.EndDate,
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
        http.Error(w, "ID required", http.StatusBadRequest)
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
