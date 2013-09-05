package main

import (
	"flag"
	"fmt"
	"os"
)

var verbose bool = false

type Options struct {
	CustomGateway string
	Sandbox       bool
	Tcpdelay      bool
	Noverify      bool
	Auth          bool
	CerFile       string
	KeyFile       string
	PemFile       string
	Expiry        int
	Repeat        int
	Ttl           int
	Tokens        []string
	Text          string
	Badge         string
	Sound         string
	Verbose       bool
}

func ParseOptions() (o Options) {
	opts := Options{}
	flag.StringVar(&opts.CustomGateway, "apn-gateway", "", "A custom APNs gateway (for testing or proxy)")
	flag.BoolVar(&opts.Sandbox, "sandbox", false, "Indicates the push notification should use the sandbox environment")
	flag.BoolVar(&opts.Tcpdelay, "tcp-delay", false, "Determines weather to delay TCP packet until it's full")
	flag.StringVar(&opts.KeyFile, "key", "", "X.509 private key in pem (Privacy Enhanced Mail) format")
	flag.StringVar(&opts.CerFile, "cer", "", "X.509 certificate in pem (Privacy Enhanced Mail) format")
	flag.StringVar(&opts.PemFile, "pem", "", "X.509 certificate/key pair stored in a pem file")
	flag.BoolVar(&opts.Verbose, "v", false, "Verbose output")
	flag.IntVar(&opts.Expiry, "expiry", 0, "UNIX date in seconds (UTC) that identifies when the notification can be discarded")
	flag.IntVar(&opts.Ttl, "ttl", 0, "Time-to-live, in seconds. Signifies how long to wait before the notification can be discarded by APNs. Differs from --expiry in that --expiry requires an actual UNIX time stamp. If both flags are provided, expiry takes precedence.")
	flag.IntVar(&opts.Repeat, "repeat", 0, "Number of times this notification should be sent.")
	flag.StringVar(&opts.Text, "text", "", "Text to send as an APN alert")
	flag.StringVar(&opts.Badge, "badge", "", "Badge value to use in payload")
	flag.StringVar(&opts.Sound, "sound", "", "Notification sound key")

	flag.Usage = func() {
		fmt.Println("apnsend - Push notification sending utility for Apple's Push Notification system (APNs)\n")
		fmt.Fprintf(os.Stderr, "Usage: apnsend [OPTIONS] token... \n")
		flag.PrintDefaults()
		fmt.Println("\nTo convert a pkcs#12 (.p12) certificate+key pair to pem, use opensll:")
		fmt.Println("\topenssl pkcs12 -in CertificateName.p12 -out CertificateName.pem -nodes")
	}

	flag.Parse()

	opts.Tokens = make([]string, 0)
	for _, a := range flag.Args() {
		opts.Tokens = append(opts.Tokens, a)
	}

	verbose = opts.Verbose

	return opts
}
