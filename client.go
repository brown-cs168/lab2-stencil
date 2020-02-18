package lab2

import (
	"context"

	"google.golang.org/grpc"
)

// Client is a wrapper for a gRPC client
type Client struct {
	conn  StoreClient
	group Group
}

// NewClient returns a gRPC client which can send requests to a gRPC server
func NewClient(addr string, group Group) (*Client, error) {
	// Establishes connection with gRPC server
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientUnaryInterceptor))
	if err != nil {
		return nil, err
	}

	// Registers client to call rpcs defined in proto file
	c := NewStoreClient(conn)

	return &Client{
		conn:  c,
		group: group,
	}, nil
}

// clientUnaryInterceptor is a client unary interceptor that injects a default timeout
func clientUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// TODO: Augment function input <ctx context.Context> with a 2 second timeout.
	// Don't forget to defer cancel or else timeout won't work!

	return invoker(ctx, method, req, reply, cc, opts...)
}

// Get sends a GetRequest to the server
func (c *Client) Get(key string) ([]byte, error) {
	// TODO: Create a background context and send the gRPC request using the client's conn.
	// If no err, return the response's item's value ([]byte)
	return nil, nil
}

// Set sends a SetRequest to the server
func (c *Client) Set(key string, value []byte) error {
	// TODO: Create a background context and send the gRPC request using the client's conn
	// You can ignore the response from the RPC call, but don't ignore the error!
	return nil
}
