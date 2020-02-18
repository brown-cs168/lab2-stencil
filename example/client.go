package example

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// Client is a wrapper for a gRPC client
type Client struct {
	conn GreeterClient
}

// NewClient returns a gRPC client which can send requests to a gRPC server
func NewClient(addr string) (*Client, error) {
	// Establishes connection with gRPC server
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientUnaryInterceptor))
	if err != nil {
		return nil, err
	}

	// Registers client to call rpcs defined in proto file
	c := NewGreeterClient(conn)
	return &Client{
		conn: c,
	}, nil
}

// clientUnaryInterceptor is a client unary interceptor that gets called before
// any request is sent. In fact, the unary interceptor is responsible for
// invoking the request sender itself.
func clientUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// Add default 5 second timeout to client interceptor
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return invoker(ctx, method, req, reply, cc, opts...)
}

// Greet sends a GreetRequest to the server
func (c *Client) Greet(msg string) (string, error) {
	// First create a context for the RPC call. A context can specify
	// many things such as a timeout, metadata, and tracing info.
	ctx := context.Background()

	// Call the RPC function defined from our proto file. The response
	// will return the message type we defined from our service.
	resp, err := c.conn.GreetRPC(ctx, &GreetRequest{Msg: msg})
	if err != nil {
		return "", err
	}
	return resp.Reply, nil
}
