package pushreceiver

// Credentials represents the configuration required to connect to
// Google's messaging service. The structure mirrors the shape used by
// the original JavaScript library.
type Credentials struct {
	GCM  GCMCredentials
	Keys Keys
	// PersistentIDs holds the list of message identifiers already
	// processed by the client. It can be nil when no messages have
	// been received yet.
	PersistentIDs []string
}

// GCMCredentials groups the identifier and security token obtained
// from the GCM checkin and register endpoints.
type GCMCredentials struct {
	AndroidID     string
	SecurityToken string
}

// Keys contains the encryption keys used to decrypt incoming
// notifications.
type Keys struct {
	PrivateKey string
	AuthSecret string
}
