package main

import "time"

type Subscription struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date,omitempty"`
	Currency  string    `json:"currency"`
	Amount    float32   `json:"amount"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
