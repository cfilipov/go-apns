// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ErrorResponse implements the APNS error response format.
//
// From the Local and Push Notification Programming Guide:
//
// If you send a notification and APNs finds the notification
// malformed or otherwise unintelligible, it returns an error-response
// packet prior to disconnecting. (If there is no error, APNs doesn't
// return anything.)
type ErrorResponse struct {
	Command    uint8 // = 8
	Status     uint8
	Identifier uint32
}

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

var errorResponseCodes = map[uint8]string{
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

// ReadFrom will read an error response from an io.Reader. Note this
// assumes a command ID has already been read and taken off the
// stream.
func (nerr *ErrorResponse) ReadFrom(r io.Reader) error {
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
func (nerr *ErrorResponse) WriteTo(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, nerr.Command)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, nerr.Status)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, nerr.Identifier)
	if err != nil {
		return err
	}

	return nil
}

// Error implements the error interface.
func (nerr *ErrorResponse) Error() string {
	return nerr.String()
}

func (nerr *ErrorResponse) String() string {
	format := "[Error Response][\n\tcommand=%v\n" +
		"\tstatus=%v (%s)\n" +
		"\tidentifier=%v\n" +
		"]"

	return fmt.Sprintf(format, nerr.Command, nerr.Status,
		errorResponseCodes[nerr.Status], nerr.Identifier)
}
