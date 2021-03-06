// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"crypto/tls"
	"net"
)

var pushHosts = [2]string{
	"gateway.push.apple.com:2195",
	"gateway.sandbox.push.apple.com:2195",
}

var feedbackHosts = [2]string{
	"feedback.push.apple.com:2196",
	"feedback.sandbox.push.apple.com:2196",
}

// Environment represents an APNs production or sandbox environment
// configuration for connections.
//
// From the Local and Push Notification Programming Guide:
//
// 		The binary interface of the production environment is available
// 		through gateway.push.apple.com, port 2195; the binary interface of
// 		the sandbox (development) environment is available through
// 		gateway.sandbox.push.apple.com, port 2195. You may establish
// 		multiple, parallel connections to the same gateway or to multiple
// 		gateway instances.
type Environment int8

const (
	DISTRIBUTION Environment = iota
	SANDBOX      Environment = iota
)

// DialAPN will create a TCP connection to Apple's APNs server using
// the certificate provided. The delay parameter tells the network
// stack to use Nagle's algorithm to batch data in TCP packets.
func DialAPN(cer *tls.Certificate, env Environment, delay bool) (net.Conn, error) {
	return Dial(cer, pushHosts[env], delay)
}

// DialFeedback will create a TCP connection to Apple's feedback service.
func DialFeedback(cer *tls.Certificate, env Environment) (net.Conn, error) {
	return Dial(cer, feedbackHosts[env], false)
}

// Dial will connect to an APNs server provided in the host parameter.
// Unless you plan on using a non-standard APNs server (like a mock
// server) then it's preferable to use DialAPN or DialFeedback.
func Dial(cer *tls.Certificate, host string, delay bool) (net.Conn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}

	// We want a net.TCPConn explicitly rather than just net.Conn so we can use 
	// SetNoDelay() to control TCP packet batching.
	tcpconn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}

	// From the Local and Push Notification Programming Guide:
	// For optimum performance, you should batch multiple notifications in a 
	// single transmission over the interface, either explicitly or using a 
	// TCP/IP Nagle's algorithm.
	tcpconn.SetNoDelay(!delay)

	// We should provide the option to connect without certificates for testing 
	// (this is convenient when one wants to setup a dummy APNs server.)
	if cer == nil {
		return tcpconn, nil
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{*cer},
	}
	tlsconn := tls.Client(tcpconn, conf)

	// From the Local and Push Notification Programming Guide:
	// To establish a trusted provider identity, you should present this 
	// certificate to APNs at connection time using peer-to-peer authentication
	err = tlsconn.Handshake()
	if err != nil {
		return nil, err
	}

	return tlsconn, nil
}
