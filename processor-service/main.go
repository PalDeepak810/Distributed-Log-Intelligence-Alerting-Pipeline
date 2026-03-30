package main

import (
	"encoding/json"
	"log"
	"net/http"

	"processor/consumer"
	"processor/monitor"
)

func main() {

	topic := "logs-topic"
	workerCount := 5

	// Start Kafka consumer
	go consumer.StartConsumer(topic, workerCount)

	// Routes
	http.HandleFunc("/logs", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(monitor.GetLogs())
	}))

	http.HandleFunc("/alerts", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(monitor.GetAlerts())
	}))

	http.HandleFunc("/metrics", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(monitor.GetMetrics())
	}))

	log.Println("Processor service started. API listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// FIXED CORS (works for Vercel + production)
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Allow all origins (safe for your project/demo)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow headers & methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
