// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"crypto/tls"
	"io"
	"net"
	"time"
)

// Connection wraps a TCP connection to APNS servers. Connection implements 
// net.Conn so that anything that understands the APNS binary interface may use 
// this connection by writing to it.
type Connection struct {
	tcpconn *net.TCPConn
	tlsconn *tls.Conn
}

var pushHosts = [2]string{
	"gateway.push.apple.com:2195",
	"gateway.sandbox.push.apple.com:2195",
}

var feedbackHosts = [2]string{
	"feedback.push.apple.com:2196",
	"feedback.sandbox.push.apple.com:2196",
}

// Ensure that Connection implements io.Writer
var _ io.Writer = &Connection{}

// Ensure that Connection implements io.Reader
var _ io.Reader = &Connection{}

// Ensure that Connection implements net.Conn
var _ net.Conn = &Connection{}

type Environment int8

// Connections can be made to either sandbox or distribution (production) 
// environments. The Connection will select the correct address based on this 
// value.
const (
	DISTRIBUTION Environment = iota
	SANDBOX      Environment = iota
)

// PushConnection provides a type-safe mechanism for making sure push 
// notifications are only sent over push gateway Connections.
type PushConnection struct {
	Connection
}

// FeedbackConnection provides a type-safe mechanism for making sure push 
// notifications cannot be sent over feedback Connections.
type FeedbackConnection struct {
	Connection
}

// Connect opens a secured connection to the push APNS gateway using the  
// certificate and environment. The delay determines weather or not TCP packets
// should be batched using Nagle's algorithm.
func Connect(cer *tls.Certificate, env Environment, delay bool) (*PushConnection, error) {
	conn, err := connect(cer, pushHosts[env])
	if err != nil {
		return nil, err
	}
	// From the Local and Push Notification Programming Guide:
	// For optimum performance, you should batch multiple notifications in a 
	// single transmission over the interface, either explicitly or using a 
	// TCP/IP Nagle algorithm.
	conn.tcpconn.SetNoDelay(!delay)
	return &PushConnection{*conn}, nil
}

// Open a secured connection to the APNS feedback server.
func ConnectFeedback(cer *tls.Certificate, env Environment) (*Connection, error) {
	return connect(cer, feedbackHosts[env])
}

// Open a secured connection to an APNS server. 
func connect(cer *tls.Certificate, gateway string) (*Connection, error) {
	raddr, err := net.ResolveTCPAddr("tcp", gateway)
	if err != nil {
		return nil, err
	}

	// We want a net.TCPConn explicitly rather than just net.Conn so we can use 
	// SetNoDelay() to control TCP packet batching.
	tcpconn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
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

	conn := &Connection{
		tcpconn: tcpconn,
		tlsconn: tlsconn,
	}
	return conn, nil
}

// Implement net.Conn and io.Read
func (conn *Connection) Read(b []byte) (n int, err error) {
	return conn.tlsconn.Read(b)
}

// Implement net.Conn and io.Writer
func (conn *Connection) Write(b []byte) (n int, err error) {
	return conn.tlsconn.Write(b)
}

// Implement net.Conn
func (conn *Connection) Close() error {
	return conn.tlsconn.Close()
}

// Implement net.Conn
func (conn *Connection) LocalAddr() net.Addr {
	return conn.tlsconn.LocalAddr()
}

// Implement net.Conn
func (conn *Connection) RemoteAddr() net.Addr {
	return conn.tlsconn.RemoteAddr()
}

// Implement net.Conn
func (conn *Connection) SetDeadline(t time.Time) error {
	return conn.tlsconn.SetDeadline(t)
}

// Implement net.Conn
func (conn *Connection) SetReadDeadline(t time.Time) error {
	return conn.tlsconn.SetReadDeadline(t)
}

// Implement net.Conn
func (conn *Connection) SetWriteDeadline(t time.Time) error {
	return conn.tlsconn.SetWriteDeadline(t)
}
