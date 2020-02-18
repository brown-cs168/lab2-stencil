package example

import "testing"

func TestSimple(t *testing.T) {
	addr := "localhost:30000"
	server, err := NewServer(addr)
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(addr)
	if err != nil {
		t.Fatal(err)
	}
	reply, err := client.Greet("This is a test message")
	if err != nil {
		t.Fatal(err)
	}
	if reply != "OK" {
		t.Fatalf("Expected: %v, Got: %v", "OK", reply)
	}
	server.GracefulStop()
}
