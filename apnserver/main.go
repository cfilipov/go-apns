// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
	"github.com/cfilipov/apns"
)

// AuthOptions contains options related to authenticating an APNs connection.
type AuthOptions struct {
	cerFile string
	keyFile string
	pemFile string
}

// ConnOptions contains options related to setting up and authenticating APNs 
// connections.
type ConnOptions struct {
	port int
}

// CMDOptions contains options which are used throughout this command.
type CMDOptions struct {
	verbose bool
}

// MockErrOptions contains options which determine how often a mocked error 
// should occur. Each option is configured as a percentage represented as an 
// integer from 0 to 100, 100 resulting in mock errors returned for every 
// notification.
type MockErrOptions struct {
	fail int
}

// Command line options grouped by type.
var (
	authOptions    *AuthOptions
	connOptions    *ConnOptions
	cmdOptions     *CMDOptions
	mockErrOptions *MockErrOptions
)

func init() {
	authOptions = &AuthOptions{}
	flag.StringVar(&authOptions.keyFile, "key", "", "X.509 private key in pem (Privacy Enhanced Mail) format")
	flag.StringVar(&authOptions.cerFile, "cer", "", "X.509 certificate in pem (Privacy Enhanced Mail) format")
	flag.StringVar(&authOptions.pemFile, "pem", "", "X.509 certificate/key pair stored in a pem file")

	cmdOptions = &CMDOptions{}
	flag.BoolVar(&cmdOptions.verbose, "v", false, "Verbose output")

	mockErrOptions = &MockErrOptions{}
	flag.IntVar(&mockErrOptions.fail, "fail", 0, "Determines how often the server should respond with an error. Accepted values are integers from 0 to 100, 100 causing all notifications to fail.")

	flag.Usage = func() {
		fmt.Println("apnserver - Push notification dummy server for Apple Push Notification system (APNs).\n")
		fmt.Fprintf(os.Stderr, "Usage: apnserver [OPTIONS] port\n")
		flag.PrintDefaults()
		fmt.Println("\nTo convert a pkcs#12 (.p12) certificate+key pair to pem, use opensll:")
		fmt.Println("\topenssl pkcs12 -in CertificateName.p12 -out CertificateName.pem -nodes")
	}

	flag.Parse()

	connOptions = &ConnOptions{}

	if flag.NArg() == 0 {
		connOptions.port = 2195
	} else {
		port, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		connOptions.port = port
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	cert, err := certificate(authOptions)
	if err != nil {
		fmt.Printf("Error loading certificate+key pair. %s\n", err)
		os.Exit(1)
	}

	if cert == nil {
		verbosePrintf("No certificate+key pair provided, using unauthenticated connection.\n")
	} else {
		verbosePrintf("Note: you may need to install the root certificate on the client machine.\n")
	}

	if mockErrOptions.fail == 0 {
		verbosePrintf("No mock errors will be used.\n")
	} else if mockErrOptions.fail > 100 {
		fmt.Printf("%d is an invalid value for --fail", mockErrOptions.fail)
		os.Exit(1)
	} else {
		verbosePrintf("Mock errors configured to %d%%.\n", mockErrOptions.fail)
	}

	conn, err := listen(cert, connOptions.port)
	if err != nil {
		fmt.Printf("Error starting TCP connection. %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listening on port %d\n", connOptions.port)

	for {
		client, err := conn.Accept()
		if err != nil {
			fmt.Printf("Unexpected error while accepting connection. %s\n", err)
			os.Exit(1)
		}
		verbosePrintf("[%v] Connected: %v\n", time.Now(), client.RemoteAddr())
		go handleClient(client, mockErrOptions)
	}
}

// verbosePrintf will print to the console only if the corresponding line option 
// is set.
func verbosePrintf(format string, a ...interface{}) (n int, err error) {
	if cmdOptions.verbose {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

// handleClient reads messages from a TCP connection.
func handleClient(conn net.Conn, mockErrOpts *MockErrOptions) {
	defer conn.Close()
	for {
		n, err := apns.ReadCommand(conn)
		if err == nil {
			verbosePrintf("Received: %s\n", n)
		}
		if err == nil {
			err = mockErr(mockErrOpts, n)
		}
		if err == nil {
			continue
		}
		// If the error is an ErrorResponse then write it to the stream.
		if resp, isResp := err.(*apns.ErrorResponse); isResp {
			verbosePrintf("Responding: %s\n", resp)
			err = resp.WriteTo(conn)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		verbosePrintf("%s\n", err)
		return
	}
}

// mockErr will randomly return an error to simulate notification failures.
func mockErr(mockErrOpts *MockErrOptions, n apns.Packet) error {
	i := rand.Intn(101-1) + 1
	if i < mockErrOpts.fail {
		if en, isEN := n.(*apns.EnhancedNotification); isEN {
			resp := &apns.ErrorResponse{
				Status:     apns.InvalidTokenStatus,
				Identifier: en.Identifier,
			}
			return resp
		}
		return io.EOF
	}
	return nil
}

// certificate creates an x.509 certificate based on the supplied options.
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

// Listen will create a TCP connection and listen for incoming
// clients. 
func listen(cer *tls.Certificate, port int) (conn net.Listener, err error) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)

	if cer != nil {
		config := &tls.Config{
			Certificates:       []tls.Certificate{*cer},
			InsecureSkipVerify: true,
		}
		conn, err = tls.Listen("tcp", addr, config)
	} else {
		taddr, _ := net.ResolveTCPAddr("tcp", addr)
		conn, err = net.ListenTCP("tcp", taddr)
	}

	return
}
