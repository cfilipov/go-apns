Go APNS Package
===============

This Go package implements the Apple Push Notification System (APNS) binary 
interface. 

Introduction
------------

This package aims to implement a light abstraction of the APNS binary interface 
for the Go programming language. The following are the design goals of this 
package:

1. Provide a minimal implementation necessary to communicate with APNS.
2. Use common interfaces in Go to allow for flexibility via composition.

Usage
-----

The various components responsible for the network and data format implement 
common interfaces such as `io.Writer`, `io.Reader` and `net.Conn`. Sending a 
push notification involves creating an `apns.Connection` to Apple's server and 
then writing to it an `apns.Notification`. 

    conn, err := apns.Connect(&cert, apns.DISTRIBUTION, false)
    notification := apns.NewSimpleNotification(token, payload)
    apns.Send(conn, notification)
    conn.Close()

Comparison with other implementations
-------------------------------------

There are already several APNS projects written in Go. The goal of this package 
is to implement the APNS interface in a style that is idiomatic to the Go 
programming language.

License
-------

Distribution and use of this project is governed by the [3-clause] Revised BSD 
license that can be found in the LICENSE file.

Further Reading
---------------

[Local and Push Notification Programming Guide. Mac Developer Library, Apple.](http://developer.apple.com/library/mac/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingWIthAPS/CommunicatingWIthAPS.html)
[Apple Push Notification Service. Wikipedia](http://en.wikipedia.org/wiki/Apple_Push_Notification_Service)
[Optimizing Connections to the Apple Push Notification Service. Apple.](https://developer.apple.com/news/index.php?id=03212012a)
