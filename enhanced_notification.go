// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"fmt"
	"io"
)

const EnhancedNotificationCMD = 1

// APNS enhanced notification format datagram. This format is the same as the 
// simple notification format except for two additional fields: Identifier and
// Expiry.
//
// From the Local and Push Notification Programming Guide:
//
// Identifier — An arbitrary value that identifies this notification. This same 
// identifier is returned in a error-response packet if APNs cannot interpret a 
// notification.
//
// Expiry — A fixed UNIX epoch date expressed in seconds (UTC) that identifies 
// when the notification is no longer valid and can be discarded. The expiry 
// value should be in network order (big endian). If the expiry value is 
// positive, APNs tries to deliver the notification at least once. You can 
// specify zero or a value less than zero to request that APNs not store the 
// notification at all.
type EnhancedNotification struct {
	Command       uint8  // = 1
	Identifier    uint32 // An arbitrary value that identifies this notification
	Expiry        uint32 // UNIX epoch date expressed in seconds (UTC)
	TokenLength   uint16 // must be big endian
	DeviceToken   []byte
	PayloadLength uint16 // must be big endian
	Payload       []byte // max 256 bytes
}

// NewEnhancedNotification creates a new push notification using the enhanced 
// notification format. 
func NewEnhancedNotification(token []byte, payload []byte, id uint32, expiry uint32) *EnhancedNotification {
	return &EnhancedNotification{
		Command:       EnhancedNotificationCMD,
		Identifier:    id,
		Expiry:        expiry,
		TokenLength:   uint16(len(token)),
		DeviceToken:   token,
		PayloadLength: uint16(len(payload)),
		Payload:       payload,
	}
}

// Write implements the io.Write interface to write the notification data. 
func (sn *EnhancedNotification) Write(w io.Writer) error {
	// From the Local and Push Notification Programming Guide:
	// The lengths of the device token and the payload must be in network order 
	// (that is, big endian).
	binary.Write(w, binary.BigEndian, sn.Command) // = 0
	binary.Write(w, binary.BigEndian, sn.Identifier)
	binary.Write(w, binary.BigEndian, sn.Expiry)
	binary.Write(w, binary.BigEndian, sn.TokenLength)
	binary.Write(w, binary.BigEndian, sn.DeviceToken)
	binary.Write(w, binary.BigEndian, sn.PayloadLength)
	binary.Write(w, binary.BigEndian, sn.Payload)
	return nil
}

// Validate will validate the notification data. This will check that the 
// command ID is correct (1 for enhanced notification), the lengths of the token 
// and payload data match the length fields and that the payload size does not 
// exceed the limit (256 bytes).
func (apn *EnhancedNotification) Validate() error {
	switch {
	case apn.Command != EnhancedNotificationCMD:
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
