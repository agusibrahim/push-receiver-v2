# pushreceiver (Go)

This directory contains an experimental Go port of the
`push-receiver-v2` library. It exposes a minimal API compatible with the
JavaScript version while leaving the underlying network protocol as
future work.

## Usage

```go
creds := pushreceiver.Credentials{
    GCM: pushreceiver.GCMCredentials{
        AndroidID:    "...",
        SecurityToken: "...",
    },
    Keys: pushreceiver.Keys{
        PrivateKey: "...",
        AuthSecret: "...",
    },
}

client, err := pushreceiver.Listen(creds, func(n pushreceiver.Notification) {
    // handle notification
})
if err != nil {
    // handle error
}
// use client ...
```

The current implementation focuses on the public API surface; message
parsing and decryption are not yet implemented.

