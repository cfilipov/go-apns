// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
)

const (
	TokenItemNumber      int8 = 1
	PayloadItemNumber    int8 = 2
	IdentifierItemNumber int8 = 3
	ExpiryItemNumber     int8 = 4
	PriorityItemNumber   int8 = 5
)

// New Notification Format (command 2)
//
// This format is a superset of the data in the enhanced notification format,
// adding a 'priority' field. However, this binary format differs in more than
// just an additional field.
//
// 		{
//			"command": 2,
// 			"device-token": "beefca5e",
// 			"identifier": 42,
// 			"expiry": 0,
// 			"priority": 50,
// 			"payload": {
// 				"aps" : {
// 		   	    	"alert" : "Hello World",
// 	       			"badge" : 0
//     			}
// 			}
//		}
type Notification struct {
	// The new notification data format is specified by command 2.
	// This field is automatically set.
	Command int8 `json:"command"` // = 2

	// The device token in binary form, as was registered by the device.
	Token string `json:"device-token"`

	// An arbitrary, opaque value that identifies this notification. This
	// identifier is used for reporting errors to your server.
	Identifier int32 `json:"identifier"`

	// A UNIX epoch date expressed in seconds (UTC) that identifies when the
	// notification is no longer valid and can be discarded.
	//
	// If this value is non-zero, APNs stores the notification tries to
	// deliver the notification at least once. Specify zero to indicate that
	// the notification expires immediately and that APNs should not store
	// the notification at all.
	Expiry int32 `json:"expiry"`

	// The notification’s priority. Provide one of the following values:
	//
	// 		10	The push message is sent immediately.
	//
	// 			The push notification must trigger an alert, sound, or badge
	//			on the device. It is an error use this priority for a push
	//			that contains only the content-available key.
	//
	// 		5 	The push message is sent at a time that conserves power on
	// 			the device receiving it.
	Priority int8 `json:"priority"`

	// The JSON-formatted payload. The payload must not be null-terminated.
	Payload JSON `json:"payload"`
}

// Implement the PushNotification interface.
func (en Notification) PushNotification() {}

func (en Notification) ReadFrom(r io.Reader) (err error) {
	return // TODO: Stub.
}

func (n Notification) WriteTo(w io.Writer) (err error) {
	token, err := hex.DecodeString(n.Token)
	if err != nil {
		return
	}
	payload, err := json.Marshal(n.Payload)
	if err != nil {
		return
	}

	tokenLen := len(token)
	payloadLen := len(payload)
	identifierLen := 4 // 4 bytes
	expiryLen := 4 // 4 bytes
	priorityLen := 1 // 1 byte

	// Calculate the size of the frame data.
	// The size of the frame data is the sum of the sizes of all items. The 
	// sum of an item is the sum of the sizes of its fields.
	//
	//                         | Number | Data len | Data         | 
	// ------------------------+--------+----------+--------------+
	// Device token            | 1 byte | 2 bytes  | 32 bytes     | 
	// Payload                 | 1 byte | 2 bytes  | <= 256 bytes | 
	// Notification identifier | 1 byte | 2 bytes  | 4 bytes      | 
	// Expiration date         | 1 byte | 2 bytes  | 4 bytes      | 
	// Priority                | 1 byte | 2 bytes  | 1 bytes      | 

	frameLen := 0 +
		1 + 2 + tokenLen + 
		1 + 2 + payloadLen + 
		1 + 2 + identifierLen + 
		1 + 2 + expiryLen + 
		1 + 2 + priorityLen

	// It is not documented, but it is possible to leave off all but the 
	// token and payload items from the frame data.

	// Write Command
	err = binary.Write(w, binary.BigEndian, NotificationCMD) // = 2
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int32(frameLen))
	if err != nil {
		return
	}
	// Write Device Token Item
	err = binary.Write(w, binary.BigEndian, TokenItemNumber)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int16(tokenLen))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, token)
	if err != nil {
		return
	}
	// Write Payload Item
	err = binary.Write(w, binary.BigEndian, PayloadItemNumber)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int16(payloadLen))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, payload)
	if err != nil {
		return
	}
	// Write Identifier Item
	err = binary.Write(w, binary.BigEndian, IdentifierItemNumber)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int16(identifierLen))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, n.Identifier)
	if err != nil {
		return
	}
	// Expiry Item
	err = binary.Write(w, binary.BigEndian, ExpiryItemNumber)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int16(expiryLen))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, n.Expiry)
	if err != nil {
		return
	}
	// Priority Item
	err = binary.Write(w, binary.BigEndian, PriorityItemNumber)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int16(priorityLen))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, n.Priority)
	if err != nil {
		return
	}
	return
}

func (nn Notification) String() string {
	nn.Command = NotificationCMD
	n, _ := json.Marshal(nn)
	return string(n)
}
