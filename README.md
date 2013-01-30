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

There are already several APNS projects written in Go. In contrast to most 
implementations, this package aims to expose a minimal functionality necessary 
for APNS communication. For this reason the basic abstraction adds very little 
on top of the APNS interface and is implemented in a way that allows one to use
selective parts of the package. For example, because a Notification writes to
an io.Writer one may elect to implement their own networking layer but continue 
to use the Notification implementations provided. 

The following is a list of notable Go implementations of APNS:

- [nicolaspaton/goapn](https://github.com/nicolaspaton/goapn)
- [mugenken/apnsender](https://github.com/mugenken/apnsender)
- [virushuo/Go-Apns](https://github.com/virushuo/Go-Apns)
- [uniqush/uniqush-push](https://github.com/uniqush/uniqush-push)

Some non-Go APNS projects:

- [notnoop/java-apns](https://github.com/notnoop/java-apns)
- [jpoz/APNS](https://github.com/jpoz/APNS)
- [simonwhitaker/PyAPNs](https://github.com/simonwhitaker/PyAPNs)

License
-------

Distribution and use of this project is governed by the [3-clause] Revised BSD 
license that can be found in the LICENSE file.

Further Reading
---------------

- [Local and Push Notification Programming Guide. Mac Developer Library, Apple.](http://developer.apple.com/library/mac/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingWIthAPS/CommunicatingWIthAPS.html)
- [Apple Push Notification Service. Wikipedia](http://en.wikipedia.org/wiki/Apple_Push_Notification_Service)
- [Optimizing Connections to the Apple Push Notification Service. Apple.](https://developer.apple.com/news/index.php?id=03212012a)
- [Nagle's algorithm](http://en.wikipedia.org/wiki/Nagle's_algorithm)
- [RFC 896 - Congestion Control in IP/TCP Internetworks. IETF.](http://tools.ietf.org/html/rfc896)
- [Socket-level Programming. Network programming with Go.](http://jan.newmarch.name/go/socket/chapter-socket.html)