apnsend
=======

A command line tool for sending push notifications using Apple's Push 
Notification System (APNs).

Example Usage
-------------

	$ apnsend -v -pem="notifyme_cert.pem" -sandbox -text="Foo" beefca5e
	Connecting to APNS. Using sandbox environment.
	1 tokens to send.
	Sending: [Simple Notification][
		command=0
		token_length=32
		token=beefca5e
		payload_length=23
		payload={"aps":{"alert":"Foo"}}
	]
