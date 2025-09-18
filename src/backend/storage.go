package backend

// Storage defines pluggable storage operations for Registrar and Matcher.
type Storage interface {
	// registration operations
	AddRegistration(clientID, objectID string)
	RemoveRegistration(clientID, objectID string)
	GetRegistrations(clientID string) []string

	// notification operations
	AddNotification(n Notification)
	GetNotifications(clientID string) []Notification
	FetchAndClearNotifications(clientID string) []Notification

	// object/version operations
	UpdateObjectVersion(objectID string, version int64)
	GetObjectVersion(objectID string) int64
	GetRegistrants(objectID string) []string
}
