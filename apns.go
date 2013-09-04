// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// UnknwonCommandErr error is used when APN data is encountered with a
// command that is unknown.
var UnknwonCommandErr = errors.New("Unknown command ID.")

const (
	simpleNotificationCMD   uint8 = 0
	enhancedNotificationCMD uint8 = 1
	errorResponseCMD        uint8 = 8
)

// Notification represents a specific set of APNs packets which are
// used for delivering push notifications.
type Notification interface {
	ReadFrom(r io.Reader) error
	Notification()
	String() string
	WriteTo(w io.Writer) error
}

// Packet represents the various data formats that may be encountered
// when communicating with APNs.
type Packet interface {
	ReadFrom(r io.Reader) error
	String() string
	WriteTo(w io.Writer) error
}

// ReadCommand will read an APNs data format from an input stream and
// return a Packet if successful.
func ReadCommand(r io.Reader) (p Packet, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Bad input.\n")
		}
	}()

	var command uint8

	err = binary.Read(r, binary.BigEndian, &command)
	if err != nil {
		return
	}

	switch command {
	case simpleNotificationCMD:
		p = new(SimpleNotification)

	case enhancedNotificationCMD:
		p = new(EnhancedNotification)

	case errorResponseCMD:
		p = new(ErrorResponse)

	default:
		err = UnknwonCommandErr
		return
	}

	err = p.ReadFrom(r)
	if err != nil {
		fmt.Println("Reading packet failed")
		return
	}

	return
}
