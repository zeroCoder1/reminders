package main

import (
	"log"
	"net/http"
)

// Recovery middleware to catch panics and log them
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	InitDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/signup", signupHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/subscriptions/add", addSubscription)
	mux.HandleFunc("/subscriptions/delete", deleteSubscription)
	mux.HandleFunc("/subscriptions/list", listSubscriptions)

	log.Println("🚀 Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", recoveryMiddleware(mux)))
}
