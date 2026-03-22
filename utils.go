package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[Error] Failed to encode JSON response: %v", err)
	}
}

func GenerateID() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("[Error] Failed to generate random ID: %v", err)
		return "00000000"
	}
	return hex.EncodeToString(bytes)
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}
