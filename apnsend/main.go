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

func main() {
	options := ParseOptions()

	// At least one token is required.
	if len(options.Tokens) == 0 {
		IfErrExitWithUsagePrompt(errors.New("No notifications to send. You must specify at least one push token."))
	}

	// Load certificate.
	var err error
	var cert tls.Certificate
	switch {
	case options.PemFile != "":
		cert, err = apns.LoadPemFile(options.PemFile)
	case options.CerFile != "" && options.KeyFile != "":
		cert, err = tls.LoadX509KeyPair(options.CerFile, options.KeyFile)
	}

	IfErrExitWithUsagePrompt(err)

	// Setup a secure connection to an APNs server.
	var conn net.Conn
	switch {
	case options.Sandbox:
		verbosePrintf("Using sandbox environment.\n")
		conn, err = apns.DialAPN(&cert, apns.SANDBOX, false)
	case options.CustomGateway != "":
		verbosePrintf("Using custom gateway: %s\n", options.CustomGateway)
		conn, err = apns.Dial(&cert, options.CustomGateway, false)
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
	aps["alert"] = options.Text
	aps["badge"] = options.Badge
	aps["sound"] = options.Sound
	payload["aps"] = aps
	jsonPayload, err := json.Marshal(payload)
	IfErrExit(err)

	// Expiry = Specific DateTime, TTL = Length of Time
	var expiry uint32
	switch {
	case options.Expiry != 0:
		expiry = uint32(options.Expiry)
	case options.Ttl != 0:
		unixTime := uint32(time.Now().Unix())
		expiry = unixTime + uint32(options.Ttl)
	}

	verbosePrintf("%d tokens to send.\n", len(options.Tokens))

	var notification apns.Notification
	for i := options.Repeat; i >= 0; i-- {
		for _, t := range options.Tokens {
			token, err := hex.DecodeString(t)
			IfErrExit(err)

			// Expiry is only available in enhanced notifications.
			if expiry != 0 {
				notification = &apns.EnhancedNotification{
					Identifier:    1,
					Expiry:        expiry,
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
			verbosePrintf("Sending: %s\n", notification)
			notification.WriteTo(conn)
		}

		// Wait for a short time before quitting to give APNs a chance to
		// return error responses, if any.
		time.Sleep(5000 * time.Millisecond)
		return
	}
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
	if verbose {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
