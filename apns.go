// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"io"
)

// From the Local and Push Notification Programming Guide:
// The payload must not exceed 256 bytes and must not be null-terminated.
const MaxPayloadSize uint16 = 256

// Notification provides a uniform way of handling a notification after it has 
// already been created. 
type Notification interface {
	Write(w io.Writer) error
	Validate() error
}

// Send will deliver a notification by writing its data to a tcp connection. A 
// secure connection to APNS may be established by creating a PushConnection.
// This function is typed such that a Notification cannot accidentally be sent 
// over a FeedbackConnection.
func Send(conn *PushConnection, apn Notification) error {
	if err := apn.Validate(); err != nil {
		return err
	}
	return apn.Write(conn)
}
