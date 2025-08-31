package pushreceiver

import (
	"crypto/rand"
	"fmt"
)

// Register combines GCM and FCM registration and returns credentials.
func Register(cfg FirebaseConfig) (*Credentials, error) {
	appID := fmt.Sprintf("wp:receiver.push.com#%s", uuid())
	gcmCreds, err := RegisterGCM(appID)
	if err != nil {
		return nil, err
	}
	inst, err := InstallFCM(cfg)
	if err != nil {
		return nil, err
	}
	fcmReg, err := RegisterFCM(cfg, inst.Fid, inst.AuthToken.Token, gcmCreds.Token)
	if err != nil {
		return nil, err
	}
	return &Credentials{
		GCM:  *gcmCreds,
		Keys: fcmReg.Keys,
		FCM:  fcmReg,
	}, nil
}

// Credentials bundles all credentials needed by the client.
type Credentials struct {
	GCM  GCMCredentials
	FCM  *FCMRegistration
	Keys KeyPair
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
