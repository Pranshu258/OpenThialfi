package backend

// Matcher provides object/version operations using a Storage backend.
type Matcher struct {
	store Storage
}

func NewMatcherWithStore(s Storage) *Matcher {
	return &Matcher{store: s}
}

// UpdateObjectVersion updates the version of an object.
func (m *Matcher) UpdateObjectVersion(objectID string, version int64) {
	m.store.UpdateObjectVersion(objectID, version)
}

// RegisterClient registers a client for an object.
func (m *Matcher) RegisterClient(objectID, clientID string) {
	m.store.AddRegistration(clientID, objectID)
}

// UnregisterClient unregisters a client from an object.
func (m *Matcher) UnregisterClient(objectID, clientID string) {
	m.store.RemoveRegistration(clientID, objectID)
}

// GetObjectVersion returns the version of an object.
func (m *Matcher) GetObjectVersion(objectID string) int64 {
	return m.store.GetObjectVersion(objectID)
}

// GetRegistrants returns the clients registered for an object.
func (m *Matcher) GetRegistrants(objectID string) []string {
	return m.store.GetRegistrants(objectID)
}
