apnsend
=======

A command line tool for sending push notifications using Apple's Push 
Notification System (APNs).

Example Usage
-------------

Send a push notification with an alert message using the production gateway

	$ apnsend -pem cert.pem -alert "Hello World" -device-token "beefca5e"

Send a push notification with an alert message using the sandbox gateway

	$ apnsend -sandbox -pem cert.pem -alert "Hello World" -device-token "beefca5e"

Send a push notification to trigger a background download of content

	$ apnsend -pem cert.pem -content-available 1 -device-token "beefca5e"

Send a push notification to update the app icon badge

	$ apnsend -pem cert.pem -badge 6 -device-token "beefca5e"

Send a push notification with a custom payload. The `-payload` argument will 
cause other payload-related arguments to be ignored (such as `-alert`, 
`-badge` etc...).

	$ apnsend -pem cert.pem -device-token "beedca5e" -payload '{"foo":"bar"}'

