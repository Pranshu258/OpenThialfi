package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"openthialfi/src/backend"
)

func TestPublishAndFetch(t *testing.T) {
	// create local registrar and matcher for a self-contained test
	store := backend.NewMemStore()
	registrar := backend.NewRegistrarWithStore(store)
	matcher := backend.NewMatcherWithStore(store)

	// register handler
	mux := http.NewServeMux()
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
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
	})
	mux.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ObjectID string `json:"object_id"`
			Version  int64  `json:"version"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		matcher.UpdateObjectVersion(req.ObjectID, req.Version)
		for _, clientID := range matcher.GetRegistrants(req.ObjectID) {
			n := backend.Notification{ClientID: clientID, ObjectID: req.ObjectID, Version: req.Version, Unknown: false}
			registrar.AddNotification(n)
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			http.Error(w, "missing client_id", http.StatusBadRequest)
			return
		}
		nots := registrar.FetchAndClear(clientID)
		json.NewEncoder(w).Encode(nots)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	// register client c1 for object o1
	regBody := map[string]string{"client_id": "c1", "object_id": "o1"}
	b, _ := json.Marshal(regBody)
	resp, err := http.Post(srv.URL+"/register", "application/json", bytes.NewReader(b))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("register failed: %v, status: %v", err, resp.StatusCode)
	}

	// publish object
	pubBody := map[string]interface{}{"object_id": "o1", "version": 42}
	b, _ = json.Marshal(pubBody)
	resp, err = http.Post(srv.URL+"/publish", "application/json", bytes.NewReader(b))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("publish failed: %v, status: %v", err, resp.StatusCode)
	}

	// small sleep to allow handler to enqueue
	time.Sleep(10 * time.Millisecond)

	// fetch notifications
	resp, err = http.Get(srv.URL + "/notifications?client_id=c1")
	if err != nil {
		t.Fatalf("fetch failed: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	var nots []map[string]interface{}
	json.Unmarshal(body, &nots)
	if len(nots) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(nots))
	}
	if int(nots[0]["Version"].(float64)) != 42 {
		t.Fatalf("unexpected version: %v", nots[0]["Version"])
	}

	// fetch again should return empty and clear happened
	resp, err = http.Get(srv.URL + "/notifications?client_id=c1")
	if err != nil {
		t.Fatalf("fetch2 failed: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &nots)
	if len(nots) != 0 {
		t.Fatalf("expected 0 notifications after clear, got %d", len(nots))
	}
}
