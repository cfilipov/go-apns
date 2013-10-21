// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"encoding/json"
	"github.com/cfilipov/apns/format"
	"io"
)

// UnknwonCommandErr error is used when APN data is encountered with a
// command that is unknown.
var UnknwonCommandErr = errors.New("Unknown command ID.")

// Notification represents a specific set of APNs packets which are
// used for delivering push notifications.
type PushNotification interface {
	PushNotification()
	ReadFrom(r io.Reader) error
	WriteTo(w io.Writer) error
	String() string
}

// Packet represents the various data formats that may be encountered
// when communicating with APNs.
type Packet interface {
	ReadFrom(r io.Reader) error
	String() string
	WriteTo(w io.Writer) error
}

func MakeNotification(data []byte) (pn PushNotification) {
	var notif format.Command
	json.Unmarshal(data, &notif)

	switch notif.Command {
	case 0:
		var n format.SimpleNotification
		json.Unmarshal([]byte(data), &n)
		pn = n
		return
	case 1:
		var n format.EnhancedNotification
		json.Unmarshal([]byte(data), &n)
		pn = n
		return
	case 2:
		var n format.Notification
		json.Unmarshal([]byte(data), &n)
		pn = n
		return
	}
	return
}

// ReadCommand will read an APNs data format from an input stream and
// return a Packet if successful.
func ReadCommand(r io.Reader) (p Packet, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Bad input.\n")
		}
	}()

	var command int8
	err = binary.Read(r, binary.BigEndian, &command)
	if err != nil {
		return
	}

	switch command {
	case format.SimpleNotificationCMD:
		p = new(format.SimpleNotification)
	case format.EnhancedNotificationCMD:
		p = new(format.EnhancedNotification)
	case format.NotificationErrorCMD:
		p = new(format.NotificationError)
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
