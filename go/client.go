package pushreceiver

import (
	"errors"
)

// Client maintains a connection to the FCM/GCM service. The
// implementation here is intentionally lightweight and focuses on the
// public API exposed by the original JavaScript library. Network
// communication and message decryption are left as future work.
type Client struct {
	creds         Credentials
	persistentIDs []string
	// TODO: add network connection fields when implementing the wire
	// protocol.
}

// Listen validates the provided credentials and returns a new Client
// instance. The function mirrors the behaviour of the JavaScript
// library's `listen` helper.
func Listen(creds Credentials, notificationCallback func(Notification)) (*Client, error) {
	// notificationCallback is currently unused; it is kept for API
	// compatibility and will be hooked once the networking portion is
	// implemented.
	_ = notificationCallback
	if creds.GCM.AndroidID == "" {
		return nil, errors.New("missing gcm.androidId in credentials")
	}
	if creds.GCM.SecurityToken == "" {
		return nil, errors.New("missing gcm.securityToken in credentials")
	}
	if creds.Keys.PrivateKey == "" {
		return nil, errors.New("missing keys.privateKey in credentials")
	}
	if creds.Keys.AuthSecret == "" {
		return nil, errors.New("missing keys.authSecret in credentials")
	}

	c := &Client{
		creds:         creds,
		persistentIDs: append([]string(nil), creds.PersistentIDs...),
	}

	// In a full implementation the client would establish a TLS
	// connection to mtalk.google.com and start streaming incoming
	// messages. For now we simply return the configured client so the
	// caller can manage its lifecycle.
	return c, nil
}

// Destroy terminates the client. The current implementation only acts
// as a placeholder and exists to match the JavaScript API.
func (c *Client) Destroy() {
	// TODO: close network resources once implemented.
	c.persistentIDs = nil
}

// Notification represents a decrypted data message stanza. The fields
// are intentionally minimal; applications are expected to unmarshal
// the payload according to their needs.
type Notification struct {
	Data map[string]string
}
