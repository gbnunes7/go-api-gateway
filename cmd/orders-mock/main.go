package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Order struct {
	ID         string `json:"id"`
	TotalPrice int    `json:"total_price"`
	CreatedAt  string `json:"created_at"`
	BillingID  string `json:"billing_id"`
	UserID     string `json:"user_id"`
}

func main() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		orders := []Order{
			{ID: "order-1", TotalPrice: 10000, CreatedAt: "2025-03-08T10:00:00Z", BillingID: "billing-1", UserID: "user-1"},
			{ID: "order-2", TotalPrice: 25000, CreatedAt: "2025-03-08T11:00:00Z", BillingID: "billing-2", UserID: "user-2"},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(orders)
	})

	fmt.Println("Orders mock service running on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
