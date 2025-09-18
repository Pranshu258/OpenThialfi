package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPublishAndFetch(t *testing.T) {
	initServices()
	setupHandlers()
	// use httptest server
	srv := httptest.NewServer(http.DefaultServeMux)
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
	body, _ := ioutil.ReadAll(resp.Body)
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
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &nots)
	if len(nots) != 0 {
		t.Fatalf("expected 0 notifications after clear, got %d", len(nots))
	}
}
