package main

import (
	"github.com/cfilipov/apns"
	"fmt"
	"os"
	"time"
)

var notif = `
{
	"command": 2,
	"device-token": "44722300f16f335fe3495cb1f760c80a9913b49b03d6ba7cdcb8726ada3a9655",
	"identifier": 22,
	"expiry": 0,
	"priority": 10,
	"payload": {
		"aps" : {
			"content-available": 1,
   	    	"alert" : "Hello World 2",
   			"badge" : 24
		},
		"msg": "Boo!"
	}
}`

func main() {
	cert, _ := apns.LoadPemFile("notifyme_cert.pem") // Load the pem file from the current dir.
	conn, _ := apns.DialAPN(&cert, apns.SANDBOX, false)

	defer conn.Close()

	// Listen for errors.
	go func() {
		for {
			p, err := apns.ReadCommand(conn)
			if err != nil {
				fmt.Printf("\nERROR: %s\n", err)
				os.Exit(1)
			}
			if p != nil {
				fmt.Printf("\nResponse: %s\n", p)
				os.Exit(1)
			}
		}
	}()

	n := apns.MakeNotification([]byte(notif))
	fmt.Printf("Sending %s\n", n.String())
	err := n.WriteTo(conn)
	if err != nil {
		fmt.Printf("\nERROR: %s\n", err)
	}

	// Wait for a short time before quitting to give APNs a chance to
	// return error responses, if any.
	time.Sleep(5000 * time.Millisecond)

	return
}
