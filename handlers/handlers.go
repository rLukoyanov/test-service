package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"test-service/downstream"
)

type request struct {
	Downstream []downstream.Target `json:"downstream"`
	Payload    json.RawMessage     `json:"payload"`
}

func proxyHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] called", name)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"error":"read body: `+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req request
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, `{"error":"invalid json: `+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		if len(req.Downstream) == 0 {
			http.Error(w, `{"error":"no downstream targets in body"}`, http.StatusBadRequest)
			return
		}

		payload := []byte(req.Payload)
		results := downstream.CallAll(req.Downstream, payload)

		resp := map[string]any{
			"handler": name,
			"results": results,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

var One   = proxyHandler("one")
var Two   = proxyHandler("two")
var Three = proxyHandler("three")
var Four  = proxyHandler("four")
