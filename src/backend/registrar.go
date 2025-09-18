package backend

import "sync"

// Registrar manages client registrations and pending notifications.
type Registrar struct {
	registrations map[string]map[string]bool // clientID -> objectID set
	notifications map[string][]Notification  // clientID -> notifications
	mu            sync.RWMutex
}

func NewRegistrar() *Registrar {
	return &Registrar{
		registrations: make(map[string]map[string]bool),
		notifications: make(map[string][]Notification),
	}
}

// Register a client for an object.
func (r *Registrar) Register(clientID, objectID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.registrations[clientID]; !ok {
		r.registrations[clientID] = make(map[string]bool)
	}
	r.registrations[clientID][objectID] = true
}

// Unregister a client from an object.
func (r *Registrar) Unregister(clientID, objectID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if reg, ok := r.registrations[clientID]; ok {
		delete(reg, objectID)
	}
}

// Add a notification for a client.
func (r *Registrar) AddNotification(n Notification) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifications[n.ClientID] = append(r.notifications[n.ClientID], n)
}

// Get notifications for a client.
func (r *Registrar) GetNotifications(clientID string) []Notification {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.notifications[clientID]
}

// FetchAndClear returns notifications for a client and clears them atomically.
func (r *Registrar) FetchAndClear(clientID string) []Notification {
	r.mu.Lock()
	defer r.mu.Unlock()
	nots := r.notifications[clientID]
	r.notifications[clientID] = nil
	return nots
}
