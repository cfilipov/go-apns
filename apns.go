// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"fmt"
	"io"
	"net"
)

// Notification represents push notifications that may be sent to APNs. 
type Notification interface {
	GetCommand() int
	GetPayload() []byte
	GetToken() []byte
	SetPayload(p []byte) error
	SetToken(t []byte) error
	Validate() error
	WriteTo(w io.Writer) error
}

// APNData represents the various data formats that may be encountered when 
// communicating with APNs.
type APNData interface {
	ReadFrom(r io.Reader) error
	String() string
	Validate() error
	WriteTo(w io.Writer) error
}

// Send will deliver a notification by writing its data to a tcp connection. A 
// secure connection to APNS may be established by using the DialXXX(...) 
// functions.
func Send(conn net.Conn, apn Notification) error {
	// if err := apn.Validate(); err != nil {
	// 	return err
	// }
	return apn.WriteTo(conn)
}

// ResolveCommand will return the corresponding data format for the command ID.
// If no corresponding data is found an error is returned instead.
func ResolveCommand(cmd uint8) (APNData, error) {
	switch cmd {
	case SimpleNotificationCMD:
		return NewSimpleNotification(), nil
	case EnhancedNotificationCMD:
		return NewEnhancedNotification(), nil
	case ErrorResponseCMD:
		return NewErrorResponse(), nil
	}
	return nil, fmt.Errorf("Unknown command: %d", cmd)
}
