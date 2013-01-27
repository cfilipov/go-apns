// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"fmt"
	"io"
)

const SimpleNotificationCMD = 0

// APNS simple notification protocol data.
type SimpleNotification struct {
	Command       uint8  // = 0
	TokenLength   uint16 // must be big endian
	DeviceToken   []byte
	PayloadLength uint16 // must be big endian
	Payload       []byte // max 256 bytes
}

// NewSimpleNotification creates a new push notification using the simple 
// notification format. 
func NewSimpleNotification(token, payload []byte) *SimpleNotification {
	return &SimpleNotification{
		Command:       SimpleNotificationCMD,
		TokenLength:   uint16(len(token)),
		DeviceToken:   token,
		PayloadLength: uint16(len(payload)),
		Payload:       payload,
	}
}

// Write implements the io.Write interface to write the notification data. 
func (sn *SimpleNotification) Write(w io.Writer) error {
	// From the Local and Push Notification Programming Guide:
	// The lengths of the device token and the payload must be in network order 
	// (that is, big endian).
	binary.Write(w, binary.BigEndian, sn.Command) // = 0
	binary.Write(w, binary.BigEndian, sn.TokenLength)
	binary.Write(w, binary.BigEndian, sn.DeviceToken)
	binary.Write(w, binary.BigEndian, sn.PayloadLength)
	binary.Write(w, binary.BigEndian, sn.Payload)
	return nil
}

// Validate will validate the notification data. This will check that the 
// command ID is correct (0 for simple notification), the lengths of the token 
// and payload data match the length fields and that the payload size does not 
// exceed the limit (256 bytes).
func (apn *SimpleNotification) Validate() error {
	switch {
	case apn.Command != SimpleNotificationCMD:
		return fmt.Errorf("Incorrect command ID in notification: %d.\n", apn.Command)
	case apn.TokenLength != uint16(len(apn.DeviceToken)):
		return fmt.Errorf("Token length %d does not match actal device token length %d.\n", apn.TokenLength, len(apn.DeviceToken))
	case apn.PayloadLength != uint16(len(apn.Payload)):
		return fmt.Errorf("Payload length %d does not match actual payload length %d.\n", apn.PayloadLength, len(apn.Payload))
	case apn.PayloadLength > MaxPayloadSize:
		return fmt.Errorf("Payload too large: %d", apn.PayloadLength)
	}
	return nil
}
