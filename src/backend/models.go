package backend

// Object represents a versioned shared object.
type Object struct {
	ID      string
	Version int64
}

// Client represents a client in the system.
type Client struct {
	ID   string
	Name string
}

// Registration links a client to an object.
type Registration struct {
	ClientID string
	ObjectID string
}

// Notification represents a notification to be sent to a client.
type Notification struct {
	ClientID string
	ObjectID string
	Version  int64
	Unknown  bool // true if version is unknown
}
