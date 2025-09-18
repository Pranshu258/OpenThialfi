package backend

import "sync"

// Matcher manages object versions and registered clients.
type Matcher struct {
	objects     map[string]int64           // objectID -> version
	registrants map[string]map[string]bool // objectID -> clientID set
	mu          sync.RWMutex
}

func NewMatcher() *Matcher {
	return &Matcher{
		objects:     make(map[string]int64),
		registrants: make(map[string]map[string]bool),
	}
}

// UpdateObjectVersion updates the version of an object.
func (m *Matcher) UpdateObjectVersion(objectID string, version int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.objects[objectID] = version
}

// RegisterClient registers a client for an object.
func (m *Matcher) RegisterClient(objectID, clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.registrants[objectID]; !ok {
		m.registrants[objectID] = make(map[string]bool)
	}
	m.registrants[objectID][clientID] = true
}

// UnregisterClient unregisters a client from an object.
func (m *Matcher) UnregisterClient(objectID, clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if reg, ok := m.registrants[objectID]; ok {
		delete(reg, clientID)
	}
}

// GetObjectVersion returns the version of an object.
func (m *Matcher) GetObjectVersion(objectID string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.objects[objectID]
}

// GetRegistrants returns the clients registered for an object.
func (m *Matcher) GetRegistrants(objectID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	clients := []string{}
	if reg, ok := m.registrants[objectID]; ok {
		for clientID := range reg {
			clients = append(clients, clientID)
		}
	}
	return clients
}
