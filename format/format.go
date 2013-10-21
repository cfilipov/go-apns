// Copyright (c) 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

type JSON map[string]interface{}

const (
	SimpleNotificationCMD   int8 = 0
	EnhancedNotificationCMD int8 = 1
	NotificationCMD         int8 = 2
	NotificationErrorCMD    int8 = 8
)

type Command struct {
	Command int8 `json:"command"`
}