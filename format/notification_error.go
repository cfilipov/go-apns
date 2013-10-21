// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

const (
	NoErrStatus              uint8 = 0
	ProcessingErrorsStatus   uint8 = 1
	MissingTokenStatus       uint8 = 2
	MissingTopicStatus       uint8 = 3
	MissingPayloadStatus     uint8 = 4
	InvalidTokenSizeStatus   uint8 = 5
	InvalidTopicSizeStatus   uint8 = 6
	InvalidPayloadSizeStatus uint8 = 7
	InvalidTokenStatus       uint8 = 8
	UnknownStatus            uint8 = 255
)

var ErrorStatusCodes = map[uint8]string{
	0:   "No errors encountered",
	1:   "Processing Errors",
	2:   "Missing Device Token",
	3:   "Missing Topic",
	4:   "Missing Payload",
	5:   "Invalid Token Size",
	6:   "Invalid Topic Size",
	7:   "Invalid Payload Size",
	8:   "Invalid Token",
	255: "None (Unknown)",
}

// ErrorResponse implements the APNS error response format.
//
// From the Local and Push Notification Programming Guide:
//
// 		If you send a notification and APNs finds the notification
// 		malformed or otherwise unintelligible, it returns an error-response
// 		packet prior to disconnecting. (If there is no error, APNs doesn't
// 		return anything.)
type NotificationError struct {
	// The packet has a command value of 8.
	// This field is automatically set.
	Command int8 // = 8

	// A one-byte status code which identifies the type of error.
	Status uint8

	// The notification identifier in the error response indicates the last
	// notification that was successfully sent. Any notifications you sent
	// after it have been discarded and must be resent. When you receive this
	// status code, stop using this connection and open a new connection.
	Identifier int32
}

// ReadFrom will read an error response from an io.Reader. Note this
// assumes a command ID has already been read and taken off the
// stream.
func (nerr NotificationError) ReadFrom(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &nerr.Status)
	if err != nil {
		return err
	}
	err = binary.Read(r, binary.BigEndian, &nerr.Identifier)
	if err != nil {
		return err
	}
	return nil
}

// WriteTo will write the entire error response to an io.Writer.
func (nerr NotificationError) WriteTo(w io.Writer) error {
	// Write Command
	err := binary.Write(w, binary.BigEndian, nerr.Command)
	if err != nil {
		return err
	}
	// Write Status
	err = binary.Write(w, binary.BigEndian, nerr.Status)
	if err != nil {
		return err
	}
	// Write Identifier
	err = binary.Write(w, binary.BigEndian, nerr.Identifier)
	if err != nil {
		return err
	}
	return nil
}

// Implement the error interface.
func (nerr NotificationError) Error() string {
	return nerr.String()
}

func (nerr NotificationError) String() string {
	nerr.Command = NotificationErrorCMD
	n, _ := json.Marshal(nerr)
	return string(n)
}
