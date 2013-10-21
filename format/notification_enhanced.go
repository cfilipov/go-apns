// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

import (
	"encoding/json"
	"encoding/hex"
	"encoding/binary"
	"io"
)

// APNS New Notification Format (command 1)
//
// This format is the same as the simple notification format except for two
// additional fields: Identifier and Expiry.
//
// 		{
//			"command": 1,
// 			"device-token": "beefca5e",
// 			"identifier": 42,
// 			"expiry": 0,
// 			"payload": {
// 				"aps" : {
// 		   	    	"alert" : "Hello World",
// 	       			"badge" : 0
//     			}
// 			}
//		}
type EnhancedNotification struct {
	// The first byte in the enhanced format is a command value of 1 (one).
	// This field is automatically set.
	Command int8 `json:"command"` // = 1

	// The device token in binary form, as was registered by the device.
	Token string `json:"device-token"`

	// A fixed UNIX epoch date expressed in seconds (UTC) that identifies when
	// the notification is no longer valid and can be discarded. The expiry
	// value should be in network order (big endian). If the expiry value is
	// positive, APNs tries to deliver the notification at least once. You can
	// specify zero or a value less than zero to request that APNs not store
	// the notification at all.
	Identifier int32 `json:"identifier"`

	// An arbitrary value that identifies this notification. This same
	// identifier is returned in a error-response packet if APNs cannot
	// interpret a notification.
	Expiry int32 `json:"expiry"`

	// The JSON-formatted payload. The payload must not be null-terminated.
	Payload JSON `json:"payload"`
}

// Implement the PushNotification interface.
func (en EnhancedNotification) PushNotification() {}

func (en EnhancedNotification) ReadFrom(r io.Reader) (err error) {
	err = binary.Read(r, binary.BigEndian, &(en.Identifier))
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &(en.Expiry))
	if err != nil {
		return
	}
	var tokenLen uint16
	err = binary.Read(r, binary.BigEndian, &(tokenLen))
	if err != nil {
		return
	}
	token := make([]byte, tokenLen)
	_, err = r.Read(token)
	if err != nil {
		return
	}
	en.Token = string(token);
	var payloadLen uint16
	err = binary.Read(r, binary.BigEndian, &(payloadLen))
	if err != nil {
		return
	}
	payloadData := make([]byte, payloadLen)
	_, err = r.Read(payloadData)
	if err != nil {
		return
	}
	payload := make(map[string]interface{})
	json.Unmarshal(payloadData, &payload)
	en.Payload = payload
	return
}

func (en EnhancedNotification) WriteTo(w io.Writer) (err error) {
	// Write Command
	err = binary.Write(w, binary.BigEndian, EnhancedNotificationCMD) // = 1
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
	token, err := hex.DecodeString(en.Token)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, uint16(len(token)))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, token)
	if err != nil {
		return
	}
	payload, err := json.Marshal(en.Payload)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, uint16(len(payload)))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, payload)
	if err != nil {
		return
	}
	return
}

func (en EnhancedNotification) String() string {
	en.Command = EnhancedNotificationCMD
	n, _ := json.Marshal(en)
	return string(n)
}
