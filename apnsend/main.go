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
	"io"
	"net"
	"os"
	"time"
	"github.com/cfilipov/apns"
)

// NoTokensErr is an error that occurs when the required token
// parameter(s) are missing.
var NoTokensErr = errors.New("No notifications to send. You must specify at least one push token.")

// ConnOptions contains options related to setting up an APNs
// connection.
type ConnOptions struct {
	customGateway string
	sandbox       bool
	tcpdelay      bool
	noverify      bool
}

// AuthOptions contains options related to authenticating an APNs
// connection.
type AuthOptions struct {
	auth    bool
	cerFile string
	keyFile string
	pemFile string
}

// APNOptions contains options related to creating a push
// notification.
type APNOptions struct {
	expiry int
	repeat int
	ttl    int
	tokens []string
}

// PayloadOptions contains fields and options used for creating APNs
// payloads.
type PayloadOptions struct {
	text  string
	badge string
	sound string
}

// CMDOptions contains options which are used throughout this command.
type CMDOptions struct {
	verbose bool
}

// Command line options grouped by type.
var (
	authOptions    *AuthOptions
	connOptions    *ConnOptions
	apnOptions     *APNOptions
	cmdOptions     *CMDOptions
	payloadOptions *PayloadOptions
)

func init() {
	connOptions = &ConnOptions{}
	flag.StringVar(&connOptions.customGateway, "apn-gateway", "", "A custom APNs gateway (for testing or proxy)")
	flag.BoolVar(&connOptions.sandbox, "sandbox", false, "Indicates the push notification should use the sandbox environment")
	flag.BoolVar(&connOptions.tcpdelay, "tcp-delay", false, "Determines weather to delay TCP packet until it's full")

	authOptions = &AuthOptions{}
	flag.StringVar(&authOptions.keyFile, "key", "", "X.509 private key in pem (Privacy Enhanced Mail) format")
	flag.StringVar(&authOptions.cerFile, "cer", "", "X.509 certificate in pem (Privacy Enhanced Mail) format")
	flag.StringVar(&authOptions.pemFile, "pem", "", "X.509 certificate/key pair stored in a pem file")

	cmdOptions = &CMDOptions{}
	flag.BoolVar(&cmdOptions.verbose, "v", false, "Verbose output")

	apnOptions = &APNOptions{}
	flag.IntVar(&apnOptions.expiry, "expiry", 0, "UNIX date in seconds (UTC) that identifies when the notification can be discarded")
	flag.IntVar(&apnOptions.ttl, "ttl", 0, "Tim-to-live, in seconds. Signifies how long to wait before the notification can be discarded by APNs. Differs from --expiry in that --expiry requires an actual UNIX time stamp. If both flags are provided, expiry takes precedence.")
	flag.IntVar(&apnOptions.repeat, "repeat", 0, "Number of times this notification should be sent.")

	payloadOptions = &PayloadOptions{}
	flag.StringVar(&payloadOptions.text, "text", "", "Text to send as an APN alert")
	flag.StringVar(&payloadOptions.badge, "badge", "", "Badge value to use in payload")
	flag.StringVar(&payloadOptions.sound, "sound", "", "Notification sound key")

	flag.Usage = func() {
		fmt.Println("apnsend - Push notification sending utility for Apple's Push Notification system (APNs)\n")
		fmt.Fprintf(os.Stderr, "Usage: apnsend [OPTIONS] token... \n")
		flag.PrintDefaults()
		fmt.Println("\nTo convert a pkcs#12 (.p12) certificate+key pair to pem, use opensll:")
		fmt.Println("\topenssl pkcs12 -in CertificateName.p12 -out CertificateName.pem -nodes")
	}

	flag.Parse()

	apnOptions.tokens = make([]string, 0)
	for _, a := range flag.Args() {
		apnOptions.tokens = append(apnOptions.tokens, a)
	}
}

func main() {
	conn, err := connection(connOptions, authOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(apnOptions.tokens) == 0 {
		fmt.Printf("ERROR: %s\n\nUsage: apnsend [OPTIONS] token... \n", err)
		flag.PrintDefaults()
	}

	defer conn.Close()

	apnc := make(chan apns.Notification)
	errc := make(chan error)
	quit := make(chan int)

	// APNs will respond with either EOF or error response data
	// depending on the type of notification or error.
	go readResp(conn, errc)
	go apnMake(apnOptions, payloadOptions, apnc, quit, errc)

	for {
		select {
		case n := <-apnc:
			verbosePrintf("Sending: %s\n", n)
			err := n.WriteTo(conn)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				os.Exit(1)
			}

		case err := <-errc:
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)

		case <-quit:
			return
		}
	}
}

// verbosePrintf will print to the console only if the corresponding
// line option is set.
func verbosePrintf(format string, a ...interface{}) (n int, err error) {
	if cmdOptions.verbose {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

// certificate creates an x.509 certificate based on the supplied
// options.
func certificate(authOpts *AuthOptions) (cert *tls.Certificate, err error) {
	var c tls.Certificate
	switch {
	case authOpts.pemFile != "":
		c, err = apns.LoadPemFile(authOpts.pemFile)
		if err != nil {
			return
		}
		cert = &c

	case authOpts.cerFile != "" && authOpts.keyFile != "":
		c, err = tls.LoadX509KeyPair(authOpts.cerFile, authOpts.keyFile)
		if err != nil {
			return
		}
		cert = &c

	default:
		cert, err = nil, nil
	}
	return
}

// connection creates a TCP connection to APNs based on the supplied
// options.
func connection(connOpts *ConnOptions, authOpts *AuthOptions) (conn net.Conn, err error) {
	verbosePrintf("Connecting to APNS. ")

	cert, err := certificate(authOpts)
	if err != nil {
		return
	}

	switch {
	case connOpts.sandbox:
		verbosePrintf("Using sandbox environment.\n")
		conn, err = apns.DialAPN(cert, apns.SANDBOX, false)

	case connOpts.customGateway != "":
		verbosePrintf("Using custom gateway: %s\n", connOptions.customGateway)
		conn, err = apns.Dial(cert, connOptions.customGateway, false)

	default:
		verbosePrintf("Using production environment.\n")
		conn, err = apns.DialAPN(cert, apns.DISTRIBUTION, false)
	}
	return
}

// readResp will listen for error responses.
func readResp(r io.Reader, errc chan error) {
	for {
		p, err := apns.ReadCommand(r)
		if err != nil {
			errc <- err
			break
		}

		errc <- fmt.Errorf("%s", p)
		break
	}
	return
}

// apnMake creates push notifications based on provided options and
// sends them along a channel.
func apnMake(apnOpts *APNOptions, payloadOpts *PayloadOptions, apnc chan apns.Notification, quit chan int, errc chan error) {
	verbosePrintf("%d tokens to send.\n", len(apnOptions.tokens))

	p, err := payload(payloadOpts)
	if err != nil {
		errc <- err
	}

	for i := apnOpts.repeat; i >= 0; i-- {
		for _, t := range apnOpts.tokens {
			token, err := hex.DecodeString(t)
			if err != nil {
				errc <- err
			}

			switch {
			case apnOpts.expiry != 0:
				apnc <- &apns.EnhancedNotification{
					Identifier:    1,
					Expiry:        uint32(apnOpts.expiry),
					TokenLength:   uint16(len(token)),
					DeviceToken:   token,
					PayloadLength: uint16(len(p)),
					Payload:       p,
				}
			case apnOpts.ttl != 0:
				unixTime := uint32(time.Now().Unix())
				apnc <- &apns.EnhancedNotification{
					Identifier:    1,
					Expiry:        unixTime + uint32(apnOpts.ttl),
					TokenLength:   uint16(len(token)),
					DeviceToken:   token,
					PayloadLength: uint16(len(p)),
					Payload:       p,
				}
			default:
				apnc <- &apns.SimpleNotification{
					TokenLength:   uint16(len(token)),
					DeviceToken:   token,
					PayloadLength: uint16(len(p)),
					Payload:       p,
				}
			}
		}
	}

	// Wait for a short time before quitting to give APNs a chance to
	// return error responses, if any.
	time.Sleep(5000 * time.Millisecond)
	quit <- 0
}

// payload creates an APNs payload based on the supplied options.
func payload(payloadOpts *PayloadOptions) (p []byte, err error) {
	var hasAps bool = false
	jsonPayload := make(map[string]interface{})
	aps := map[string]string{}

	if len(payloadOpts.text) != 0 {
		hasAps = true
		aps["alert"] = payloadOpts.text
	}

	if len(payloadOpts.badge) != 0 {
		hasAps = true
		aps["badge"] = payloadOpts.badge
	}

	if len(payloadOpts.sound) != 0 {
		hasAps = true
		aps["sound"] = payloadOpts.sound
	}

	if hasAps {
		jsonPayload["aps"] = aps
	}

	p, err = json.Marshal(jsonPayload)
	return
}
