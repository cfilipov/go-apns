package main

import (
	"encoding/hex"
	"github.com/cfilipov/apns"
)

func main() {
	cert, _ := apns.LoadPemFile("notifyme_cert.pem") // Load the pem file from the current dir.
	conn, _ := apns.DialAPN(&cert, apns.SANDBOX, false)
	token, _ := hex.DecodeString("beefca5e") // Use a real APNs token.
	payload := []byte(`{"aps":{"alert":"Hello World!"}}`)

	notification := apns.SimpleNotification{
		TokenLength:   uint16(len(token)),
		DeviceToken:   token,
		PayloadLength: uint16(len(payload)),
		Payload:       payload,
	}

	notification.WriteTo(conn)
}
