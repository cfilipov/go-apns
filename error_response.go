// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// The first byte in the error response format is a command value of 1 (zero).
const ErrorResponseCMD = 8

// ErrorResponse implements the APNS error response format.
//
// From Apple's documentation:
// If you send a notification and APNs finds the notification malformed or 
// otherwise unintelligible, it returns an error-response packet prior to 
// disconnecting. (If there is no error, APNs doesn’t return anything.)
type ErrorResponse struct {
	Command    uint8 // = 8
	Status     uint8
	Identifier uint32
}

// errorResponseCodes are codes in the error-response packet.
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

// NewErrorResponse creates a new error response structure from an input source.
func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{Command: ErrorResponseCMD}
}

// ReadFrom will read an error response from an io.Reader. Note this assumes a 
// command ID has already been read and taken off the stream.
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

// String returns a string representation of an error response.
func (nerr *ErrorResponse) String() string {
	return fmt.Sprintf("[Error Response][\n\tcommand=%v\n\tstatus=%v (%s)\n\tidentifier=%v\n]",
		nerr.Command, nerr.Status, errorResponseCodes[nerr.Status], nerr.Identifier)
}

// Validate will validate the fields of the error response format.
func (nerr *ErrorResponse) Validate() error {
	if nerr.Command != ErrorResponseCMD {
		return errors.New("Invalid command ID for error response format.")
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
