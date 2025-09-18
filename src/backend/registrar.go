package backend

// Registrar manages client registrations and pending notifications using a Storage backend.
type Registrar struct {
	store Storage
}

func NewRegistrarWithStore(s Storage) *Registrar {
	return &Registrar{store: s}
}

// Register a client for an object.
func (r *Registrar) Register(clientID, objectID string) {
	r.store.AddRegistration(clientID, objectID)
}

// Unregister a client from an object.
func (r *Registrar) Unregister(clientID, objectID string) {
	r.store.RemoveRegistration(clientID, objectID)
}

// Add a notification for a client.
func (r *Registrar) AddNotification(n Notification) {
	r.store.AddNotification(n)
}

// Get notifications for a client.
func (r *Registrar) GetNotifications(clientID string) []Notification {
	return r.store.GetNotifications(clientID)
}

// FetchAndClear returns notifications for a client and clears them atomically.
func (r *Registrar) FetchAndClear(clientID string) []Notification {
	return r.store.FetchAndClearNotifications(clientID)
}
