package example

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// Server is used to implement grpc.GreeterServer.
type Server struct {
	server *grpc.Server
}

// NewServer starts gRPC server and returns server for graceful stop in future
func NewServer(addr string) (*Server, error) {
	// Creates a new gRPC server with interceptor
	s := grpc.NewServer(grpc.UnaryInterceptor(serverUnaryInterceptor))

	// gRPC uses TCP as its underlying network protocol so we need to open a tcp listener
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Create our server struct. This struct is responsible for implementing
	// the RPC interface provided by our proto file
	server := &Server{
		server: s,
	}

	// Registers gRPC server for rpcs defined in proto file
	RegisterGreeterServer(s, server)

	// Begins serving requests on separate goroutine
	go s.Serve(lis)

	return server, nil
}

// serverUnaryInterceptor is a server unary interceptor that gets called before
// any request is processed. In fact, the unary interceptor is responsible for
// calling the request handler itself.
func serverUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	// Print the method being called
	fmt.Println(info.FullMethod)

	return handler(ctx, req)
}

// GreetRPC prints msg in GreetRequest and responds with Ok
func (s *Server) GreetRPC(ctx context.Context, msg *GreetRequest) (*GreetReply, error) {
	fmt.Printf("Msg received: %v\n", msg.Msg)
	return &GreetReply{Reply: "OK"}, nil
}

// GracefulStop stops the gRPC server. Useful for testing
func (s *Server) GracefulStop() {
	s.server.GracefulStop()
}
