// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Utility for sending push notifications using the Apple's Push
Notification System (APNs) Go library.
*/
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cfilipov/apns"
	"github.com/cfilipov/apns/format"
	"net"
	"os"
	"time"
)

var token = flag.String("device-token", "", "A custom APNs gateway (for testing or proxy)")
var notifJSON = flag.String("notification-json", "", "A custom APNs gateway (for testing or proxy)")
var notifCMD = flag.Int("command", 2, "An identifier specifying the apns binary data format to use. 0: Simple, 1: Enhanced, 2:Default")
var customGateway = flag.String("apn-gateway", "", "A custom APNs gateway (for testing or proxy)")
var tcpDelay = flag.Bool("tcp-delay", false, "Determines weather to delay TCP packet until it's full")
var verbose = flag.Bool("v", false, "Verbose output")
var priority = flag.Int("priority", 10, "The notification’s priority. Default is 10. Possible values: 10 (The push message is sent immediately), 5 (The push message is sent at a time that conserves power on the device receiving it).")
var expiry = flag.Int("expiry", 0, "UNIX date in seconds (UTC) that identifies when the notification can be discarded")
var keyFile = flag.String("key", "apns-key.pem", "X.509 private key in pem (Privacy Enhanced Mail) format")
var cerFile = flag.String("cer", "apns-cer.pem", "X.509 certificate in pem (Privacy Enhanced Mail) format")
var pemFile = flag.String("pem", "apns.pem", "X.509 certificate/key pair stored in a pem file. If this argument is specified then other certificate/key arguments are ignored.")
var sandbox = flag.Bool("sandbox", false, "Indicates the push notification should use the sandbox environment")
var badge = flag.String("badge", "", "Badge value to use in payload")
var sound = flag.String("sound", "", "Notification sound key")
var contentAvailable = flag.String("content-available", "", "Provide this key with a value of 1 to indicate that new content is available. This is used to support Newsstand apps and background content downloads.")
var alert = flag.String("alert", "", "Alert text to send as an APN alert")
var payload = flag.String("payload", "", "Raw (JSON) payload to send. This overrides all other aps payload arguments such as -text -badge and -sound options.")
var ttl = flag.Int("ttl", 0, "Time-to-live, in seconds. Signifies how long to wait before the notification can be discarded by APNs. Differs from --expiry in that --expiry requires an actual UNIX time stamp. If both flags are provided, expiry takes precedence.")

func init() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Println("apnsend - Push notification sending utility for Apple's Push Notification system (APNs)\n")
		fmt.Fprintf(os.Stderr, "Usage: apnsend -pem <certificate> -alert <text> -device-token <token> \n")
		flag.PrintDefaults()
		fmt.Println("\nTo convert a pkcs#12 (.p12) certificate+key pair to pem, use opensll:")
		fmt.Println("\topenssl pkcs12 -in CertificateName.p12 -out CertificateName.pem -nodes")
	}
}

func main() {
	var err error

	// Sanity check the arguments.

	if *token == "" {
		fmt.Println("Missing argument: -device-token")
		flag.Usage()
		os.Exit(1)
	}
	if *pemFile == "" && *cerFile == "" && *keyFile == "" {
		fmt.Println("Missing argument: -pem, -cer, or -key required")
		flag.Usage()
		os.Exit(1)
	}
	if *payload == "" && *alert == "" && *badge == "" && *sound == "" && *contentAvailable == "" {
		fmt.Println("Missing argument: -payload, -alert, -badge, -sound, or -content-available required")
		flag.Usage()
		os.Exit(1)
	}

	// Calculate the expiry, if applicable.

	var expiryTime int32 // Expiry = Specific DateTime, TTL = Length of Time

	if *expiry != 0 {
		expiryTime = int32(*expiry)
	}
	if *ttl != 0 {
		unixTime := int32(time.Now().Unix())
		expiryTime = unixTime + int32(*ttl)
	}

	// Load the certificate.

	var cert tls.Certificate

	if *pemFile != "" {
		cert, err = apns.LoadPemFile(*pemFile)
	} else {
		cert, err = tls.LoadX509KeyPair(*cerFile, *keyFile)
	}

	if err != nil {
		fmt.Printf("\nERROR: %s\n", err)
		os.Exit(1)
	}

	// Setup a secure connection to an APNs server.

	var conn net.Conn

	if *sandbox {
		if *verbose {
			fmt.Printf("Using sandbox environment.\n")
		}
		conn, err = apns.DialAPN(&cert, apns.SANDBOX, *tcpDelay)
	} else if *customGateway != "" {
		if *verbose {
			fmt.Printf("Using custom gateway: %s\n", *customGateway)
		}
		conn, err = apns.Dial(&cert, *customGateway, *tcpDelay)
	} else {
		if *verbose {
			fmt.Printf("Using production environment.\n")
		}
		conn, err = apns.DialAPN(&cert, apns.DISTRIBUTION, *tcpDelay)
	}

	if err != nil {
		fmt.Printf("\nERROR: %s\n", err)
		os.Exit(1)
	}

	defer conn.Close()

	// Listen for error responses.

	go func() {
		for {
			p, err := apns.ReadCommand(conn)
			if err != nil {
				fmt.Printf("\nERROR: %s\n", err)
				os.Exit(1)
			}
			if err != nil {
				fmt.Printf("\nAPNs Response: %s\n", p)
				os.Exit(1)
			}
		}
	}()

	// Create a notification instance.

	var notif apns.PushNotification

	if *notifJSON != "" {
		notif = apns.MakeNotification([]byte(*notifJSON))
	} else {
		var p format.JSON

		// Create a payload unless one is provided by the -payload argument.

		if len(*payload) == 0 {
			p = make(map[string]interface{})
			aps := map[string]string{}
			if *alert != "" {
				aps["alert"] = *alert
			}
			if *badge != "" {
				aps["badge"] = *badge
			}
			if *sound != "" {
				aps["sound"] = *sound
			}
			if *contentAvailable != "" {
				aps["content-available"] = *contentAvailable
			}
			p["aps"] = aps
			if err != nil {
				fmt.Printf("\nERROR: %s\n", err)
				os.Exit(1)
			}
		} else {
			json.Unmarshal([]byte(*payload), &p)
		}

		// Create a notification instance.

		if *notifCMD == 0 {
			notif = format.SimpleNotification{
				Token:   *token,
				Payload: p,
			}
		} else if *notifCMD == 1 {
			notif = format.EnhancedNotification{
				Identifier: 1, 
				Expiry:     expiryTime,
				Token:      *token,
				Payload:    p,
			}
		} else { // *notifCMD == 2
			notif = format.Notification{
				Identifier: 1,
				Expiry:     expiryTime,
				Token:      *token,
				Priority:   int8(*priority),
				Payload:    p,
			}
		}
	}

	// Write the notification to output.

	if *verbose {
		fmt.Printf("Sending: %s\n", notif)
	}
	notif.WriteTo(conn)

	// Wait for a short time before quitting to give APNs a chance to
	// return error responses, if any.

	time.Sleep(5000 * time.Millisecond)

	return
}
