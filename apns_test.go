package apns

import (
	// "crypto/tls"
	"encoding/hex"
	"encoding/json"
	"net"
	"testing"
)

var (
	conn *PushConnection
	apn  Notification
)

func BenchmarkSendBuffered(b *testing.B) {
	b.StopTimer()
	token, _ := hex.DecodeString("47ee04b9e673f7ddc86cd126d2504b3661336a60c17e06cec382881b1bd839f8")
	jsonPayload := make(map[string]interface{})
	jsonPayload["aps"] = map[string]string{"alert": "Hello World 1"}
	payload, _ := json.Marshal(jsonPayload)
	apn = NewEnhancedNotification()
	apn.SetToken(token)
	apn.SetPayload(payload)

	// cert, err := tls.LoadX509KeyPair("test_data/cert.pem", "test_data/key.pem")
	// if err != nil {
	// 	b.Error("Failed to load certificate/key pair.")
	// }

	c, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		b.Error("Connection error.")
	}
	conn := &PushConnection{c}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SendBuffered(conn, apn)
	}
	b.StopTimer()
	c.Close()
}

func BenchmarkSendBufferedWithDelay(b *testing.B) {
	b.StopTimer()
	token, _ := hex.DecodeString("47ee04b9e673f7ddc86cd126d2504b3661336a60c17e06cec382881b1bd839f8")
	jsonPayload := make(map[string]interface{})
	jsonPayload["aps"] = map[string]string{"alert": "Hello World 1"}
	payload, _ := json.Marshal(jsonPayload)
	apn = NewEnhancedNotification()
	apn.SetToken(token)
	apn.SetPayload(payload)

	// cert, err := tls.LoadX509KeyPair("test_data/cert.pem", "test_data/key.pem")
	// if err != nil {
	//  b.Error("Failed to load certificate/key pair.")
	// }

	raddr, _ := net.ResolveTCPAddr("tcp", "localhost:8080")
	c, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		b.Error("Connection error.")
	}
	c.SetNoDelay(false)
	conn := &PushConnection{c}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SendBuffered(conn, apn)
	}
	b.StopTimer()
	c.Close()
}

func BenchmarkSend(b *testing.B) {
	b.StopTimer()
	token, _ := hex.DecodeString("47ee04b9e673f7ddc86cd126d2504b3661336a60c17e06cec382881b1bd839f8")
	jsonPayload := make(map[string]interface{})
	jsonPayload["aps"] = map[string]string{"alert": "Hello World 2"}
	payload, _ := json.Marshal(jsonPayload)
	apn = NewEnhancedNotification()
	apn.SetToken(token)
	apn.SetPayload(payload)

	// cert, err := tls.LoadX509KeyPair("test_data/cert.pem", "test_data/key.pem")
	// if err != nil {
	//  b.Error("Failed to load certificate/key pair.")
	// }

	c, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		b.Error("Connection error.")
	}
	conn := &PushConnection{c}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Send(conn, apn)
	}
	b.StopTimer()
	c.Close()
}

func BenchmarkSendWithDelay(b *testing.B) {
	b.StopTimer()
	token, _ := hex.DecodeString("47ee04b9e673f7ddc86cd126d2504b3661336a60c17e06cec382881b1bd839f8")
	jsonPayload := make(map[string]interface{})
	jsonPayload["aps"] = map[string]string{"alert": "Hello World 2"}
	payload, _ := json.Marshal(jsonPayload)
	apn = NewEnhancedNotification()
	apn.SetToken(token)
	apn.SetPayload(payload)

	// cert, err := tls.LoadX509KeyPair("test_data/cert.pem", "test_data/key.pem")
	// if err != nil {
	// 	b.Error("Failed to load certificate/key pair.")
	// }

	raddr, _ := net.ResolveTCPAddr("tcp", "localhost:8080")
	c, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		b.Error("Connection error.")
	}
	c.SetNoDelay(false)
	conn := &PushConnection{c}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Send(conn, apn)
	}
	b.StopTimer()
	c.Close()
}
