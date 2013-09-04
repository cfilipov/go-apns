// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"fmt"
	"io"
)

// APNS enhanced notification format datagram. This format is the same
// as the simple notification format except for two additional fields:
// Identifier and Expiry.
//
// From the Local and Push Notification Programming Guide:
//
// Identifier — An arbitrary value that identifies this notification.
// This same identifier is returned in a error- response packet if
// APNs cannot interpret a notification.
//
// Expiry — A fixed UNIX epoch date expressed in seconds (UTC) that
// identifies when the notification is no longer valid and can be
// discarded. The expiry value should be in network order (big
// endian). If the expiry value is positive, APNs tries to deliver the
// notification at least once. You can specify zero or a value less
// than zero to request that APNs not store the notification at all.
type EnhancedNotification struct {
	Command       uint8 // = 1
	Identifier    uint32
	Expiry        uint32
	TokenLength   uint16
	DeviceToken   []byte
	PayloadLength uint16
	Payload       []byte
}

// WriteTo implements the io.Write interface to write the notification
// data.
func (en *EnhancedNotification) WriteTo(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, enhancedNotificationCMD) // = 0
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, en.Identifier)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, en.Expiry)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, en.TokenLength)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, en.DeviceToken)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, en.PayloadLength)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, en.Payload)
	if err != nil {
		return
	}

	return
}

// ReadFrom will fill an EnhancedNotification with data from an input
// stream.
func (en *EnhancedNotification) ReadFrom(r io.Reader) (err error) {
	err = binary.Read(r, binary.BigEndian, &(en.Identifier))
	if err != nil {
		return
	}

	err = binary.Read(r, binary.BigEndian, &(en.Expiry))
	if err != nil {
		return
	}

	err = binary.Read(r, binary.BigEndian, &(en.TokenLength))
	if err != nil {
		return
	}

	en.DeviceToken = make([]byte, en.TokenLength)
	_, err = r.Read(en.DeviceToken)
	if err != nil {
		return
	}

	err = binary.Read(r, binary.BigEndian, &(en.PayloadLength))
	if err != nil {
		return
	}

	en.Payload = make([]byte, en.PayloadLength)
	_, err = r.Read(en.Payload)
	if err != nil {
		return
	}

	return
}

func (en *EnhancedNotification) Notification() {}

func (en *EnhancedNotification) String() string {
	format := "[Enhanced Notification][\n" +
		"\tcommand=%v\n" +
		"\tidentifier=%v\n" +
		"\texpiry=%v\n" +
		"\ttoken_length=%v\n" +
		"\ttoken=%x\n" +
		"\tpayload_length=%v\n" +
		"\tpayload=%s\n" +
		"]"

	return fmt.Sprintf(format, en.Command, en.Identifier, en.Expiry,
		en.TokenLength, en.DeviceToken, en.PayloadLength, en.Payload)
}
