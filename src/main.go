package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"openthialfi/src/backend"
)

var (
	registrar *backend.Registrar
	matcher   *backend.Matcher
	once      sync.Once
)

func initServices() {
	once.Do(func() {
		store := backend.NewMemStore()
		registrar = backend.NewRegistrarWithStore(store)
		matcher = backend.NewMatcherWithStore(store)
	})
}

// register handler: POST /register {"client_id":"c1","object_id":"o1"}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID string `json:"client_id"`
		ObjectID string `json:"object_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	registrar.Register(req.ClientID, req.ObjectID)
	matcher.RegisterClient(req.ObjectID, req.ClientID)
	w.WriteHeader(http.StatusOK)
}

// unregister handler: POST /unregister {"client_id":"c1","object_id":"o1"}
func unregisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID string `json:"client_id"`
		ObjectID string `json:"object_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	registrar.Unregister(req.ClientID, req.ObjectID)
	matcher.UnregisterClient(req.ObjectID, req.ClientID)
	w.WriteHeader(http.StatusOK)
}

// publish handler: POST /publish {"object_id":"o1","version":42}
func publishHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ObjectID string `json:"object_id"`
		Version  int64  `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	matcher.UpdateObjectVersion(req.ObjectID, req.Version)
	// notify registrants
	for _, clientID := range matcher.GetRegistrants(req.ObjectID) {
		n := backend.Notification{ClientID: clientID, ObjectID: req.ObjectID, Version: req.Version, Unknown: false}
		registrar.AddNotification(n)
	}
	w.WriteHeader(http.StatusOK)
}

// fetch notifications: GET /notifications?client_id=c1
func fetchHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "missing client_id", http.StatusBadRequest)
		return
	}
	notifs := registrar.FetchAndClear(clientID)
	json.NewEncoder(w).Encode(notifs)
}

func setupHandlers() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/unregister", unregisterHandler)
	http.HandleFunc("/publish", publishHandler)
	http.HandleFunc("/notifications", fetchHandler)
}

func main() {
	initServices()
	setupHandlers()
	fmt.Println("Thialfi-like Notification Service Backend Started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
