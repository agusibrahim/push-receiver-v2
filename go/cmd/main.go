package main

import (
	"fmt"
	"log"

	"pushreceiver"
)

func main() {
	cfg := pushreceiver.FirebaseConfig{
		APIKey:    "AIzaSyC23PJFvsGcPV-mxk-OOc0d3o9uCuiVZX4",
		AppID:     "1:640680687175:android:0144cefbfc47ad6c",
		ProjectID: "authexample-ffdf9",
	}

	creds, err := pushreceiver.Register(cfg)
	if err != nil {
		log.Fatalf("registration failed: %v", err)
	}

	fmt.Println("GCM Token:", creds.GCM.Token)
	fmt.Println("FCM Token:", creds.FCM.Token)
	fmt.Println("Public Key:", creds.Keys.PublicKey)
}
