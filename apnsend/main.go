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
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/cfilipov/apns"
	"net"
	"os"
	"time"
)

var tokens []string
var customGateway = flag.String("apn-gateway", "", "A custom APNs gateway (for testing or proxy)")
var tcpDelay = flag.Bool("tcp-delay", false, "Determines weather to delay TCP packet until it's full")
var verbose = flag.Bool("v", false, "Verbose output")
var expiry = flag.Int("expiry", 0, "UNIX date in seconds (UTC) that identifies when the notification can be discarded")
var repeat = flag.Int("repeat", 1, "Number of times this notification should be sent.")
var keyFile = flag.String("key", "apns-key.pem", "X.509 private key in pem (Privacy Enhanced Mail) format")
var cerFile	= flag.String("cer", "apns-cer.pem", "X.509 certificate in pem (Privacy Enhanced Mail) format")
var pemFile = flag.String("pem", "apns.pem", "X.509 certificate/key pair stored in a pem file")
var sandbox = flag.Bool("sandbox", false, "Indicates the push notification should use the sandbox environment")
var badge = flag.String("badge", "", "Badge value to use in payload")
var sound =	flag.String("sound", "", "Notification sound key")
var text = flag.String("text", "", "Text to send as an APN alert")
var ttl = flag.Int("ttl", 0, "Time-to-live, in seconds. Signifies how long to wait before the notification can be discarded by APNs. Differs from --expiry in that --expiry requires an actual UNIX time stamp. If both flags are provided, expiry takes precedence.")

func init() {
	flag.Parse()

	tokens = make([]string, 0)
	for _, a := range flag.Args() {
		tokens = append(tokens, a)
	}

	flag.Usage = func() {
		fmt.Println("apnsend - Push notification sending utility for Apple's Push Notification system (APNs)\n")
		fmt.Fprintf(os.Stderr, "Usage: apnsend [OPTIONS] token... \n")
		flag.PrintDefaults()
		fmt.Println("\nTo convert a pkcs#12 (.p12) certificate+key pair to pem, use opensll:")
		fmt.Println("\topenssl pkcs12 -in CertificateName.p12 -out CertificateName.pem -nodes")
	}
}

func main() {
	// At least one token is required.
	if len(tokens) == 0 {
		IfErrExitWithUsagePrompt(errors.New("No notifications to send. You must specify at least one push token."))
	}

	// Load certificate.
	var err error
	var cert tls.Certificate
	switch {
	case *pemFile != "":
		cert, err = apns.LoadPemFile(*pemFile)
	case *cerFile != "" && *keyFile != "":
		cert, err = tls.LoadX509KeyPair(*cerFile, *keyFile)
	}

	IfErrExitWithUsagePrompt(err)

	// Setup a secure connection to an APNs server.
	var conn net.Conn
	switch {
	case *sandbox:
		verbosePrintf("Using sandbox environment.\n")
		conn, err = apns.DialAPN(&cert, apns.SANDBOX, false)
	case *customGateway != "":
		verbosePrintf("Using custom gateway: %s\n", *customGateway)
		conn, err = apns.Dial(&cert, *customGateway, false)
	default:
		verbosePrintf("Using production environment.\n")
		conn, err = apns.DialAPN(&cert, apns.DISTRIBUTION, false)
	}

	IfErrExit(err)
	defer conn.Close()

	// Listen for errors.
	go func() {
		for {
			p, err := apns.ReadCommand(conn)
			IfErrExit(err)
			IfErrExit(fmt.Errorf("%s", p))
		}
	}()

	// JSON payload.
	payload := make(map[string]interface{})
	aps := map[string]string{}
	aps["alert"] = *text
	aps["badge"] = *badge
	aps["sound"] = *sound
	payload["aps"] = aps
	jsonPayload, err := json.Marshal(payload)
	IfErrExit(err)

	// Expiry = Specific DateTime, TTL = Length of Time
	var expiryTime uint32
	switch {
	case *expiry != 0:
		expiryTime = uint32(*expiry)
	case *ttl != 0:
		unixTime := uint32(time.Now().Unix())
		expiryTime = unixTime + uint32(*ttl)
	}

	verbosePrintf("%d tokens to send.\n", len(tokens))

	var notification apns.Notification
	for i := *repeat; i > 0; i-- {
		for _, t := range tokens {
			token, err := hex.DecodeString(t)
			IfErrExit(err)

			// Expiry is only available in enhanced notifications.
			if expiryTime != 0 {
				notification = &apns.EnhancedNotification{
					Identifier:    1,
					Expiry:        expiryTime,
					TokenLength:   uint16(len(token)),
					DeviceToken:   token,
					PayloadLength: uint16(len(jsonPayload)),
					Payload:       jsonPayload,
				}
			} else {
				notification = &apns.SimpleNotification{
					TokenLength:   uint16(len(token)),
					DeviceToken:   token,
					PayloadLength: uint16(len(jsonPayload)),
					Payload:       jsonPayload,
				}
			}

			// Send the notification.
			verbosePrintf("Sending %d of %d: %s\n", (*repeat - i + 1), *repeat, notification)
			notification.WriteTo(conn)
		}
	}

	// Wait for a short time before quitting to give APNs a chance to
	// return error responses, if any.
	time.Sleep(5000 * time.Millisecond)
	return
}

func IfErrExit(err error) {
	if err != nil {
		fmt.Printf("\nERROR: %s\n", err)
		os.Exit(1)
	}
}

func IfErrExitWithUsagePrompt(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n\nUsage: apnsend [OPTIONS] token... \n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func verbosePrintf(format string, a ...interface{}) (n int, err error) {
	if *verbose {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
