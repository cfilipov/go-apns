package main

import (
	"encoding/json"
	"fmt"
	"github.com/cfilipov/apns"
	"github.com/cfilipov/apns/format"
)

var data = `{
	"command": 2,
	"device-token": "beefca5e",
	"identifier": 42,
	"expiry": 0,
	"priority": 50,
	"payload": {
		"aps" : {
   	    	"alert" : "Bob wants to play poker",
   			"badge" : 0
		}
	}
}`

func main() {
	var notif format.Command
	json.Unmarshal([]byte(data), &notif)

	fmt.Printf("Command: %d\n", notif.Command)

	switch notif.Command {
	case 0:
		
	case 1:
		
	case 2:
		var n format.Notification
		json.Unmarshal([]byte(data), &n)
		fmt.Printf("Notification: %s\n", n.String())
	}

	var n format.Notification
	json.Unmarshal([]byte(data), &n)
	fmt.Printf("Notification: %s\n", n.String())

	sn := &format.SimpleNotification{
		Token:   "beefca5e",
		Payload: n.Payload,
	}
	fmt.Printf("Simple Notification: %s\n", sn.String())

	en := &format.EnhancedNotification{
		Identifier: 1,
		Expiry:     0,
		Token:      "beefca5e",
		Payload:    n.Payload,
	}
	fmt.Printf("Enhanced Notification: %s\n", en.String())

	nn := &format.Notification{
		Identifier: 1,
		Expiry:     0,
		Token:      "beefca5e",
		Priority:   5,
		Payload:    n.Payload,
	}
	fmt.Printf("Notification: %s\n", nn.String())

	xn := apns.MakeNotification([]byte(data))
	fmt.Printf("Notification: %s\n", xn.String())

	return
}
