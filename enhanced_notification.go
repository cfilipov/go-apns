// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// The first byte in the enhanced format is a command value of 1 (zero).
	EnhancedNotificationCMD = 1

	// From the Local and Push Notification Programming Guide:
	// The payload must not exceed 256 bytes and must not be null-terminated.
	MaxEnhancedPayloadSize uint16 = 256
)

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
func NewEnhancedNotification() *EnhancedNotification {
	return &EnhancedNotification{Command: EnhancedNotificationCMD}
}

// Write implements the io.Write interface to write the notification data. 
func (sn *EnhancedNotification) WriteTo(w io.Writer) error {
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

func (sn *EnhancedNotification) ReadFrom(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &(sn.Identifier))
	err = binary.Read(r, binary.BigEndian, &(sn.Expiry))

	// Read the toke length
	err = binary.Read(r, binary.BigEndian, &(sn.TokenLength))
	if err != nil {
		return err
	}

	// Read the device token
	sn.DeviceToken = make([]byte, sn.TokenLength)
	_, err = r.Read(sn.DeviceToken)
	if err != nil {
		return err
	}

	// Read the payload length
	err = binary.Read(r, binary.BigEndian, &(sn.PayloadLength))
	if err != nil {
		return err
	}

	// Read the device token
	sn.Payload = make([]byte, sn.PayloadLength)
	_, err = r.Read(sn.Payload)
	if err != nil {
		return err
	}
	return nil
}

func (sn *EnhancedNotification) String() string {
	return fmt.Sprintf("[Enhanced Notification][\n\tcommand=%v\n\tidentifier=%v\n\texpiry=%v\n\ttoken_length=%v\n\ttoken=%x\n\tpayload_length=%v\n\tpayload=%s\n]",
		sn.Command, sn.Identifier, sn.Expiry, sn.TokenLength, sn.DeviceToken, sn.PayloadLength, sn.Payload)
}

func (sn *EnhancedNotification) GetToken() []byte {
	return sn.DeviceToken
}

func (sn *EnhancedNotification) GetPayload() []byte {
	return sn.Payload
}

func (sn *EnhancedNotification) GetCommand() int {
	return EnhancedNotificationCMD
}

func (sn *EnhancedNotification) SetPayload(p []byte) error {
	sn.PayloadLength = uint16(len(p))
	sn.Payload = p
	return nil
}

func (sn *EnhancedNotification) SetToken(t []byte) error {
	sn.TokenLength = uint16(len(t))
	sn.DeviceToken = t
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
	case apn.PayloadLength > MaxEnhancedPayloadSize:
		return fmt.Errorf("Payload too large: %d", apn.PayloadLength)
	}
	return nil
}
