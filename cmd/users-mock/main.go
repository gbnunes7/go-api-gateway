package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	OrderID string `json:"order_id"`
}

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if delayStr := r.URL.Query().Get("delay"); delayStr != "" {
			if delay, err := strconv.Atoi(delayStr); err == nil && delay > 0 {
				time.Sleep(time.Duration(delay) * time.Second)
			}
		}

		users := []User{
			{ID: "user-1", Name: "John Doe", Email: "john@example.com", OrderID: "order-1"},
			{ID: "user-2", Name: "Jane Doe", Email: "jane@example.com", OrderID: "order-2"},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(users)
	})

	fmt.Println("Users mock service running on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
