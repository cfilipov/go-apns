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
	// The first byte in the simple format is a command value of 0 (zero).
	SimpleNotificationCMD = 0

	// From the Local and Push Notification Programming Guide:
	// The payload must not exceed 256 bytes and must not be null-terminated.
	MaxSimplePayloadSize uint16 = 256
)

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
func NewSimpleNotification() *SimpleNotification {
	return &SimpleNotification{Command: SimpleNotificationCMD}
}

// WriteTo implements the io.Write interface to write the notification data. 
func (sn *SimpleNotification) WriteTo(w io.Writer) error {
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

func (sn *SimpleNotification) ReadFrom(r io.Reader) error {
	// Read the toke length
	err := binary.Read(r, binary.BigEndian, &(sn.TokenLength))
	if err != nil {
		return fmt.Errorf("Error reading token length from input.\n%s", err)
	}

	// Read the device token
	sn.DeviceToken = make([]byte, sn.TokenLength)
	_, err = r.Read(sn.DeviceToken)
	if err != nil {
		return fmt.Errorf("Error reading device token from input.\n%s", err)
	}

	// Read the payload length
	err = binary.Read(r, binary.BigEndian, &(sn.PayloadLength))
	if err != nil {
		return fmt.Errorf("Error reading payload length from input.\n%s", err)
	}

	// Read the device token
	sn.Payload = make([]byte, sn.PayloadLength)
	_, err = r.Read(sn.Payload)
	if err != nil {
		return fmt.Errorf("Error reading payload data from input.\n%s", err)
	}

	return nil
}

func (sn *SimpleNotification) String() string {
	return fmt.Sprintf("[Simple Notification][\n\tcommand=%v\n\ttoken_length=%v\n\ttoken=%x\n\tpayload_length=%v\n\tpayload=%s\n]",
		sn.Command, sn.TokenLength, sn.DeviceToken, sn.PayloadLength, sn.Payload)
}

func (sn *SimpleNotification) GetToken() []byte {
	return sn.DeviceToken
}

func (sn *SimpleNotification) GetPayload() []byte {
	return sn.Payload
}

func (sn *SimpleNotification) GetCommand() int {
	return SimpleNotificationCMD
}

func (sn *SimpleNotification) SetPayload(p []byte) error {
	sn.PayloadLength = uint16(len(p))
	sn.Payload = p
	return nil
}

func (sn *SimpleNotification) SetToken(t []byte) error {
	sn.TokenLength = uint16(len(t))
	sn.DeviceToken = t
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
	case apn.PayloadLength > MaxSimplePayloadSize:
		return fmt.Errorf("Payload too large: %d", apn.PayloadLength)
	}
	return nil
}
