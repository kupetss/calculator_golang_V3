package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"calculator_golangV3/config/calculator"
	"calculator_golangV3/config/handlers"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func StartServer() {
	calculator.Initialize()
	handlers.Initialize()
	err := os.MkdirAll("database", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	file, err := os.OpenFile("database/results.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	file.Close()

	log.Println("Server is starting...")
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/calculate", handlers.HandleComputation).Methods("POST")
	router.HandleFunc("/api/v1/expressions/{id}", handlers.HandleFetch).Methods("GET")
	router.HandleFunc("/api/v1/expressions", handlers.HandleFetchAll).Methods("GET")
	router.HandleFunc("/internal/task", handlers.HandleTaskOrchestration).Methods("POST")

	http.ListenAndServe(":8080", enableCORS(router))
}

func main() {
	StartServer()
}
