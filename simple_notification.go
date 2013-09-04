// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"fmt"
	"io"
)

// APNS simple notification protocol data.
//
// From the Local and Push Notification Programming Guide: 
//
// The first
// byte in the simple format is a command value of 0 (zero). The
// lengths of the device token and the payload must be in network
// order (that is, big endian). In addition, you should encode the
// device token in binary format. The payload must not exceed 256
// bytes and must not be null-terminated.
type SimpleNotification struct {
	Command       uint8 // = 0
	TokenLength   uint16
	DeviceToken   []byte
	PayloadLength uint16
	Payload       []byte
}

// WriteTo implements the io.Write interface to write the notification
// packets.
func (sn *SimpleNotification) WriteTo(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, simpleNotificationCMD) // = 0
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, sn.TokenLength)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, sn.DeviceToken)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, sn.PayloadLength)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, sn.Payload)
	if err != nil {
		return
	}

	return
}

// ReadFrom will fill a SimpleNotification with data from an input
// stream.
func (sn *SimpleNotification) ReadFrom(r io.Reader) (err error) {
	err = binary.Read(r, binary.BigEndian, &(sn.TokenLength))
	if err != nil {
		return
	}

	sn.DeviceToken = make([]byte, sn.TokenLength)
	_, err = r.Read(sn.DeviceToken)
	if err != nil {
		return
	}

	err = binary.Read(r, binary.BigEndian, &(sn.PayloadLength))
	if err != nil {
		return
	}

	sn.Payload = make([]byte, sn.PayloadLength)
	_, err = r.Read(sn.Payload)
	if err != nil {
		return
	}

	return
}

func (sn *SimpleNotification) Notification() {}

func (sn *SimpleNotification) String() string {
	format := "[Simple Notification][\n" +
		"\tcommand=%v\n" +
		"\ttoken_length=%v\n" +
		"\ttoken=%x\n" +
		"\tpayload_length=%v\n" +
		"\tpayload=%s\n" +
		"]"

	return fmt.Sprintf(format, sn.Command, sn.TokenLength, sn.DeviceToken,
		sn.PayloadLength, sn.Payload)
}
