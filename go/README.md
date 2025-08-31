# pushreceiver (Go)

This directory contains a Go port of the `push-receiver-v2` library.
It implements the registration flow against the public Google and
Firebase endpoints and exposes structures mirroring the JavaScript
version.

## Usage

```go
cfg := pushreceiver.FirebaseConfig{
    APIKey:    "<firebase api key>",
    AppID:     "<firebase app id>",
    ProjectID: "<firebase project id>",
    VapidKey:  "<vapid public key>",
}

creds, err := pushreceiver.Register(cfg)
if err != nil {
    panic(err)
}
fmt.Println("GCM Token:", creds.GCM.Token)
```

The client connection used to receive messages is still under
development.

