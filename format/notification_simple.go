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

// Simple Notification Format (command 0)
//
// 		{
//			"command": 0,
// 			"device-token": "beefca5e",
// 			"payload": {
// 				"aps" : {
// 		   	    	"alert" : "Hello World",
// 	       			"badge" : 0
//     			}
// 			}
//		}
type SimpleNotification struct {
	// The first byte in the simple format is a command value of 0 (zero). 
	// This field is automatically set.
	Command int8 `json:"command"` // = 0

	// The device token in binary form, as was registered by the device.
	Token string `json:"device-token"`

	// The JSON-formatted payload. The payload must not be null-terminated.
	Payload JSON `json:"payload"`
}

// Implement the PushNotification interface.
func (sn SimpleNotification) PushNotification() {}

func (sn SimpleNotification) ReadFrom(r io.Reader) (err error) {
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
	sn.Token = string(token);
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
	sn.Payload = payload
	return
}

func (sn SimpleNotification) WriteTo(w io.Writer) (err error) {
	// Write Command
	err = binary.Write(w, binary.BigEndian, SimpleNotificationCMD) // = 0
	if err != nil {
		return
	}
	// Write Token
	token, err := hex.DecodeString(sn.Token)
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
	// Write Payload
	payload, err := json.Marshal(sn.Payload)
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

func (sn SimpleNotification) String() string {
	sn.Command = SimpleNotificationCMD
	n, _ := json.Marshal(sn)
	return string(n)
}
