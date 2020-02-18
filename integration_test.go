package lab2

import (
	fmt "fmt"
	"sync"
	"testing"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

func TestGRPC(t *testing.T) {
	addr := "localhost:30000"
	server, err := NewServer(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer server.GracefulStop()

	aClient, err := NewClient(addr, Group_ADMIN)
	if err != nil {
		t.Fatal(err)
	}
	uClient, err := NewClient(addr, Group_USER)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("set", func(t *testing.T) {
		err := aClient.Set("key1", []byte("value1"))
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get", func(t *testing.T) {
		err = aClient.Set("key1", []byte("value1"))
		if err != nil {
			t.Fatal(err)
		}
		value, err := aClient.Get("key1")
		if err != nil {
			t.Fatal(err)
		}
		if string(value) != "value1" {
			t.Fatalf("Expected: %v, Got: %v\n", "value1", string(value))
		}
		value, err = aClient.Get("key2")
		s := status.Convert(err)
		if s.Message() != ErrKeyNotFound.Error() {
			t.Fatalf("Expected: %v, Got: %v\n", ErrKeyNotFound, err)
		}
	})

	t.Run("concurrent", func(t *testing.T) {
		errChan := make(chan error, 2000)
		var wg sync.WaitGroup

		for j := 0; j < 100; j++ {

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(i int) {
					err := aClient.Set(fmt.Sprintf("key%d", i), []byte(fmt.Sprintf("value%d", i)))
					if err != nil {
						errChan <- err
					}
					wg.Done()
				}(i)
			}
			wg.Wait()

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(i int) {
					value, err := aClient.Get(fmt.Sprintf("key%d", i))
					if err != nil {
						errChan <- err
					} else if string(value) != fmt.Sprintf("value%d", i) {
						errChan <- fmt.Errorf("expected: %v, got: %v", fmt.Sprintf("value%d", i), string(value))
					}
					wg.Done()
				}(i)
			}
			wg.Wait()

		}

		close(errChan)

		for err := range errChan {
			t.Fatal(err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		server.slowConn.Store(true)
		defer server.slowConn.Store(false)

		err := aClient.Set("key0", []byte("value0"))
		st := status.Convert(err)
		if st.Code() != codes.DeadlineExceeded {
			t.Fatalf("Expected: %v, Got: %v\n", codes.DeadlineExceeded, st.Code())
		}

		server.blockConn.Store(true)
		defer server.blockConn.Store(false)

		err = aClient.Set("key0", []byte("value0"))
		st = status.Convert(err)
		if st.Message() != ErrConnBlocked.Error() {
			t.Fatalf("Expected: %v, Got: %v\n", ErrConnBlocked.Error(), st.Message())
		}
	})

	t.Run("permissions", func(t *testing.T) {
		err := aClient.Set("key1", []byte("admin"))
		if err != nil {
			t.Fatal(err)
		}
		err = uClient.Set("key1", []byte("user"))
		st := status.Convert(err)
		if st.Message() != ErrUnauthorized.Error() {
			t.Fatalf("Expected: %v, Got: %v\n", ErrUnauthorized.Error(), st.Message())
		}
		_, err = uClient.Get("key1")
		st = status.Convert(err)
		if st.Message() != ErrUnauthorized.Error() {
			t.Fatalf("Expected: %v, Got: %v\n", ErrUnauthorized.Error(), st.Message())
		}
		value, err := aClient.Get("key1")
		if err != nil {
			t.Fatal(err)
		}
		if string(value) != "admin" {
			t.Fatalf("Expected: %v, Got: %v\n", "admin", string(value))
		}
	})

	t.Run("net.Listen error", func(t *testing.T) {
		_, err := NewServer(addr)
		if err == nil {
			t.Fatal("Expected error, Got none")
		}
	})
}
