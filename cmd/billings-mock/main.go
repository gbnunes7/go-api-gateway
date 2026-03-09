package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Billing struct {
	ID          string `json:"id"`
	PaymentType string `json:"payment_type"`
	PaidValue   int    `json:"paid_value"`
	OrderID     string `json:"order_id"`
}

func main() {
	http.HandleFunc("/billings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		billings := []Billing{
			{ID: "billing-1", PaymentType: "credit", PaidValue: 10000, OrderID: "order-1"},
			{ID: "billing-2", PaymentType: "pix", PaidValue: 25000, OrderID: "order-2"},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(billings)
	})

	fmt.Println("Billings mock service running on :8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatal(err)
	}
}
