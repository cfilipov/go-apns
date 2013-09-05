Go APNs Package
===============

This Go package implements the Apple Push Notification System (APNs) binary 
interface. 

Usage
-----

This package provides simple interfaces for establishing authenticated 
connections to APNs gateways and sending notifications. The function 
`apns.DialAPN(...)` returns a `net.Conn` which is authenticated and ready to 
receive data. Notifications implement the `io.Writer` and `io.Reader` 
interfaces so that sending a notification is done by having it write to a 
connection.

	func main() {
		// Load the pem file from the current dir.
		cert, _ := apns.LoadPemFile("notifyme_cert.pem")
		conn, _ := apns.DialAPN(&cert, apns.SANDBOX, false)
		// Use a real APNs token.
		token, _ := hex.DecodeString("beefca5e")
		payload := []byte(`{"aps":{"alert":"Hello World!"}}`)

		notification := apns.SimpleNotification{
			TokenLength:   uint16(len(token)),
			DeviceToken:   token,
			PayloadLength: uint16(len(payload)),
			Payload:       payload,
		}

		notification.WriteTo(conn)
	}

Other Go implementations of APNs:

- [nicolaspaton/goapn](https://github.com/nicolaspaton/goapn)
- [mugenken/apnsender](https://github.com/mugenken/apnsender)
- [virushuo/Go-Apns](https://github.com/virushuo/Go-Apns)
- [uniqush/uniqush-push](https://github.com/uniqush/uniqush-push)

Other non-Go APNs projects:

- [notnoop/java-apns](https://github.com/notnoop/java-apns)
- [jpoz/APNS](https://github.com/jpoz/APNS)
- [simonwhitaker/PyAPNs](https://github.com/simonwhitaker/PyAPNs)

apnsend
-------

The [apnsend](https://github.com/cfilipov/go-apns/tree/master/apnsend) utility is a command line tool which uses the apns package for 
sending push notifications.

apnserver
---------

The apnserver utility will respond to the APNs protocol with mock data. The 
server can be configured to a specific mock failure rate to simulate errors 
and dropped connections.

License
-------

Distribution and use of this project is governed by the [3-clause] Revised BSD 
license that can be found in the LICENSE file.

Related Info
------------

- [Local and Push Notification Programming Guide. Mac Developer Library, Apple.](http://developer.apple.com/library/mac/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingWIthAPS/CommunicatingWIthAPS.html)
- [Optimizing Connections to the Apple Push Notification Service. Apple.](https://developer.apple.com/news/index.php?id=03212012a)
- [Apple Push Service Protocol. ios-rev tumblr.](http://ios-rev.tumblr.com/post/13032664009/apple-push-service-protocol-ios5-os-x-10-7)