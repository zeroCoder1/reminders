package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()

	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/subscriptions/add", addSubscription)
	http.HandleFunc("/subscriptions/delete", deleteSubscription)
	http.HandleFunc("/subscriptions/list", listSubscriptions)

	log.Println("🚀 Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
