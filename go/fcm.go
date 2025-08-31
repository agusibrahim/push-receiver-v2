package pushreceiver

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	fisURL      = "https://firebaseinstallations.googleapis.com/v1/"
	fcmRegURL   = "https://fcmregistrations.googleapis.com/v1/"
	fcmEndpoint = "https://fcm.googleapis.com/fcm/send"
)

// FirebaseConfig contains Firebase project configuration.
type FirebaseConfig struct {
	APIKey    string
	AppID     string
	ProjectID string
	VapidKey  string
}

// InstallFCM registers the device installation with Firebase.
func InstallFCM(cfg FirebaseConfig) (*Installation, error) {
	url := fmt.Sprintf("%sprojects/%s/installations", fisURL, cfg.ProjectID)
	body := map[string]any{
		"appId":       cfg.AppID,
		"authVersion": "FIS_v2",
		"fid":         generateFID(),
		"sdkVersion":  "w:0.6.4",
	}
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", cfg.APIKey)
	req.Header.Set("x-firebase-client", base64.StdEncoding.EncodeToString([]byte(`{"heartbeats":[],"version":2}`)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("install failed: %s", resp.Status)
	}
	var inst Installation
	if err := json.Unmarshal(data, &inst); err != nil {
		return nil, err
	}
	return &inst, nil
}

// Installation represents response from InstallFCM.
type Installation struct {
	Name      string `json:"name"`
	Fid       string `json:"fid"`
	AuthToken struct {
		Token string `json:"token"`
	} `json:"authToken"`
}

// RegisterFCM registers FCM using credentials from InstallFCM and GCM token.
func RegisterFCM(cfg FirebaseConfig, authToken, token string) (*FCMRegistration, error) {
	keys, err := createKeys()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%sprojects/%s/registrations", fcmRegURL, cfg.ProjectID)
	body := map[string]any{
		"web": map[string]any{
			"applicationPubKey": cfg.VapidKey,
			"auth":              urlSafe(keys.AuthSecret),
			"endpoint":          fmt.Sprintf("%s/%s", fcmEndpoint, token),
			"p256dh":            urlSafe(keys.PublicKey),
		},
	}
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", cfg.APIKey)
	req.Header.Set("x-goog-firebase-installations-auth", authToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("register fcm failed: %s", resp.Status)
	}
	var reg FCMRegistration
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	reg.Keys = keys
	return &reg, nil
}

// FCMRegistration contains keys and server response.
type FCMRegistration struct {
	Name  string `json:"name"`
	Token string `json:"token"`
	Keys  KeyPair
}

// KeyPair holds generated keys for message decryption.
type KeyPair struct {
	PrivateKey string
	PublicKey  string
	AuthSecret string
}

func createKeys() (KeyPair, error) {
	curve := elliptic.P256()
	priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return KeyPair{}, err
	}
	pub := elliptic.Marshal(curve, x, y)
	auth := make([]byte, 16)
	if _, err := rand.Read(auth); err != nil {
		return KeyPair{}, err
	}
	return KeyPair{
		PrivateKey: base64.StdEncoding.EncodeToString(priv),
		PublicKey:  base64.StdEncoding.EncodeToString(pub),
		AuthSecret: base64.StdEncoding.EncodeToString(auth),
	}, nil
}

func urlSafe(b64 string) string {
	s := strings.ReplaceAll(b64, "+", "-")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.TrimRight(s, "=")
	return s
}

func generateFID() string {
	b := make([]byte, 17)
	rand.Read(b)
	b[0] = 0x70 + (b[0] % 16)
	return base64.StdEncoding.EncodeToString(b)
}
