// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package apns implements the Apple Push Notification System (APNS) binary 
interface. The library is organized into two high level parts: the connection 
and data. 

Connection wraps a TCP connection to an APNS server and handles the details of 
authenticating and selecting environments. Connection implements the net.Conn 
interface, after creating the connection one can use the Connection like an 
ordinary net.Conn. 

Connection comes in two flavors, PushConnection and FeedbackConnection. 
PushConnection can only be used to send push notifications and the 
FeedbackConnection can only be used to read feedback messages, this restriction 
is enforced by the type system when using the Send() method.

The data part of the library is made up of Notification, Feedback and 
ErrorResponse. These work via standard interfaces io.Writer and io.Reader.
*/
package apns
