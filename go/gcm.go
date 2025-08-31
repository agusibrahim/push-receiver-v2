package pushreceiver

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"

	p "pushreceiver/internal/proto"
)

const (
	checkinURL  = "https://android.clients.google.com/checkin"
	registerURL = "https://android.clients.google.com/c2dm/register3"
)

// CheckIn performs the Android checkin request.
func CheckIn(androidID, securityToken uint64) (uint64, uint64, error) {
	req := buildCheckinRequest(androidID, securityToken)
	httpReq, err := http.NewRequest("POST", checkinURL, bytes.NewReader(req))
	if err != nil {
		return 0, 0, err
	}
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	if resp.StatusCode != 200 {
		return 0, 0, fmt.Errorf("checkin failed: %s", resp.Status)
	}
	return parseCheckinResponse(body)
}

// RegisterGCM performs the GCM register request for the given appID and returns
// the registration token along with androidID and securityToken.
func RegisterGCM(appID string) (*GCMCredentials, error) {
	aid, token, err := CheckIn(0, 0)
	if err != nil {
		return nil, err
	}
	data := "app=org.chromium.linux&X-subtype=" + appID + "&device=" + fmt.Sprintf("%d", aid) + "&sender=" + serverKey()
	httpReq, err := http.NewRequest("POST", registerURL, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Authorization", fmt.Sprintf("AidLogin %d:%d", aid, token))
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("register failed: %s", resp.Status)
	}
	if !bytes.Contains(body, []byte("=")) {
		return nil, errors.New(string(body))
	}
	parts := bytes.SplitN(body, []byte("="), 2)
	return &GCMCredentials{
		Token:         string(parts[1]),
		AndroidID:     fmt.Sprintf("%d", aid),
		SecurityToken: fmt.Sprintf("%d", token),
	}, nil
}

// GCMCredentials represents credentials returned by RegisterGCM.
type GCMCredentials struct {
	Token         string
	AndroidID     string
	SecurityToken string
}

func buildCheckinRequest(androidID, securityToken uint64) []byte {
	chromeBuild := []byte{}
	chromeBuild = append(chromeBuild, p.VarintField(1, 2)...)
	chromeBuild = append(chromeBuild, p.BytesField(2, []byte("63.0.3234.0"))...)
	chromeBuild = append(chromeBuild, p.VarintField(3, 1)...)

	checkin := []byte{}
	checkin = append(checkin, p.VarintField(12, 3)...)
	checkin = append(checkin, p.BytesField(13, chromeBuild)...)

	req := []byte{}
	if androidID != 0 {
		req = append(req, p.VarintField(2, androidID)...)
	}
	req = append(req, p.BytesField(4, checkin)...)
	if securityToken != 0 {
		req = append(req, p.Fixed64Field(13, securityToken)...)
	}
	req = append(req, p.VarintField(14, 3)...)
	req = append(req, p.VarintField(22, 0)...)
	return req
}

func parseCheckinResponse(b []byte) (androidID, securityToken uint64, err error) {
	for i := 0; i < len(b); {
		key, n := decodeVarint(b[i:])
		if n == 0 {
			break
		}
		i += n
		field := int(key >> 3)
		wire := int(key & 7)
		switch field {
		case 7, 8: // fixed64
			if wire != 1 || i+8 > len(b) {
				return 0, 0, fmt.Errorf("invalid wire type")
			}
			v := binary.LittleEndian.Uint64(b[i : i+8])
			i += 8
			if field == 7 {
				androidID = v
			} else {
				securityToken = v
			}
		default:
			// skip other fields
			var skip int
			switch wire {
			case 0:
				_, skip = decodeVarint(b[i:])
			case 1:
				skip = 8
			case 2:
				l, m := decodeVarint(b[i:])
				skip = int(l)
				i += m
			case 5:
				skip = 4
			default:
				return 0, 0, fmt.Errorf("unknown wire type")
			}
			i += skip
		}
	}
	return
}

func decodeVarint(b []byte) (uint64, int) {
	var x uint64
	var s uint
	for i, c := range b {
		if c < 0x80 {
			if i > 9 || i == 9 && c > 1 {
				return 0, 0
			}
			return x | uint64(c)<<s, i + 1
		}
		x |= uint64(c&0x7f) << s
		s += 7
	}
	return 0, 0
}

func serverKey() string {
	data := []byte{
		0x04, 0x33, 0x94, 0xf7, 0xdf, 0xa1, 0xeb, 0xb1,
		0xdc, 0x03, 0xa2, 0x5e, 0x15, 0x71, 0xdb, 0x48,
		0xd3, 0x2e, 0xed, 0xed, 0xb2, 0x34, 0xdb, 0xb7,
		0x47, 0x3a, 0x0c, 0x8f, 0xc4, 0xcc, 0xe1, 0x6f,
		0x3c, 0x8c, 0x84, 0xdf, 0xab, 0xb6, 0x66, 0x3e,
		0xf2, 0x0c, 0xd4, 0x8b, 0xfe, 0xe3, 0xf9, 0x76,
		0x2f, 0x14, 0x1c, 0x63, 0x08, 0x6a, 0x6f, 0x2d,
		0xb1, 0x1a, 0x95, 0xb0, 0xce, 0x37, 0xc0, 0x9c, 0x6e,
	}
	return base64.StdEncoding.EncodeToString(data)
}
