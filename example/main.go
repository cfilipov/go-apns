package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/cfilipov/apns"
)

func main() {
	cert, _ := apns.LoadPemFile("notifyme_cert.pem") // Load the pem file from the current dir.
	conn, _ := apns.DialAPN(&cert, apns.SANDBOX, false)
	token, _ := hex.DecodeString("beefca5e") // Use a real APNs token.

	jsonData := make(map[string]interface{})
	aps := map[string]string{}
	aps["alert"] = "Hello World!"
	jsonData["aps"] = aps
	payload, _ := json.Marshal(jsonData)

	notification := apns.SimpleNotification{
		TokenLength:   uint16(len(token)),
		DeviceToken:   token,
		PayloadLength: uint16(len(payload)),
		Payload:       payload,
	}

	notification.WriteTo(conn)
}
