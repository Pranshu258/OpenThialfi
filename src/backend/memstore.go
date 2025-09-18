package backend

import "sync"

// MemStore is a simple in-memory storage for testing and development.
type MemStore struct {
	mu          sync.RWMutex
	regs        map[string]map[string]bool // clientID -> objectID set
	notifs      map[string][]Notification  // clientID -> notifications
	objects     map[string]int64           // objectID -> version
	registrants map[string]map[string]bool // objectID -> clientID set
}

func NewMemStore() *MemStore {
	return &MemStore{
		regs:        make(map[string]map[string]bool),
		notifs:      make(map[string][]Notification),
		objects:     make(map[string]int64),
		registrants: make(map[string]map[string]bool),
	}
}

func (m *MemStore) AddRegistration(clientID, objectID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.regs[clientID]; !ok {
		m.regs[clientID] = make(map[string]bool)
	}
	m.regs[clientID][objectID] = true
	if _, ok := m.registrants[objectID]; !ok {
		m.registrants[objectID] = make(map[string]bool)
	}
	m.registrants[objectID][clientID] = true
}

func (m *MemStore) RemoveRegistration(clientID, objectID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if reg, ok := m.regs[clientID]; ok {
		delete(reg, objectID)
	}
	if reg2, ok := m.registrants[objectID]; ok {
		delete(reg2, clientID)
	}
}

func (m *MemStore) GetRegistrations(clientID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []string{}
	if reg, ok := m.regs[clientID]; ok {
		for oid := range reg {
			out = append(out, oid)
		}
	}
	return out
}

func (m *MemStore) AddNotification(n Notification) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifs[n.ClientID] = append(m.notifs[n.ClientID], n)
}

func (m *MemStore) GetNotifications(clientID string) []Notification {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.notifs[clientID]
}

func (m *MemStore) FetchAndClearNotifications(clientID string) []Notification {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := m.notifs[clientID]
	m.notifs[clientID] = nil
	return n
}

func (m *MemStore) UpdateObjectVersion(objectID string, version int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.objects[objectID] = version
}

func (m *MemStore) GetObjectVersion(objectID string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.objects[objectID]
}

func (m *MemStore) GetRegistrants(objectID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []string{}
	if reg, ok := m.registrants[objectID]; ok {
		for cid := range reg {
			out = append(out, cid)
		}
	}
	return out
}
