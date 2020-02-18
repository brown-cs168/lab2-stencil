package lab2

import (
	"context"
	"errors"
	"net"
	"sync"

	"go.uber.org/atomic"
	"google.golang.org/grpc"
)

// Errors produced by server
var (
	ErrKeyNotFound  = errors.New("Key not found")
	ErrConnBlocked  = errors.New("gRPC connection blocked")
	ErrUnauthorized = errors.New("Operation unauthorized")
)

// Server is used to implement grpc.GreeterServer.
type Server struct {
	server *grpc.Server
	store  map[string]Item

	blockConn *atomic.Bool
	slowConn  *atomic.Bool

	// We can embed an RWMutex in the struct itself. However, according to
	// https://github.com/uber-go/guide/blob/master/style.md#avoid-embedding-types-in-public-structs
	// one should avoid embedding types in public structs. This is done just to
	// show an example of embedding types in structs. For reminder, public structs
	// are capitalized and private structs are not.
	sync.RWMutex
}

// NewServer starts gRPC server and returns server for graceful stop in future
func NewServer(addr string) (*Server, error) {
	// Creates a new gRPC server with unary interceptor
	s := grpc.NewServer(grpc.UnaryInterceptor(serverUnaryInterceptor))

	// gRPC uses TCP as its underlying network protocol so we need to open a tcp listener
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Create struct that implements RPC interface
	server := &Server{
		server:    s,
		store:     make(map[string]Item),
		blockConn: atomic.NewBool(false),
		slowConn:  atomic.NewBool(false),
	}

	// Registers gRPC server for rpcs defined in proto file
	RegisterStoreServer(s, server)

	// Begins serving requests on separate goroutine
	go s.Serve(lis)

	return server, nil
}

// serverUnaryInterceptor is a server unary interceptor that can simulate
// slow nodes, network partition, or killed nodes
func serverUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// TODO: Get access to Server struct defined above using info.Server

	// TODO: If server blockConn is true, return ErrConnBlocked.
	// This simulates network partitions

	// TODO: If server slowConn is true, sleep for more than 2 seconds. This
	// This simulates a slow node.

	return handler(ctx, req)
}

// SetRPC inserts a key-value pair into the in-memory store
func (s *Server) SetRPC(ctx context.Context, msg *SetReq) (*SetReply, error) {
	// TODO: Fill out method. Make sure to check for Group in request and compare
	// with item belonging to key if item exists. Since enums are just ints, you can
	// use >, ==, < comparators with enums! Also don't forget to lock mutex!

	return nil, nil
}

// GetRPC returns a value for the given key or ErrKeyNotFound if key does not exist in the in-memory store.
// Please use the provided ErrKeyNotFound var defined in errors.go
func (s *Server) GetRPC(ctx context.Context, msg *GetReq) (*GetReply, error) {
	// TODO: Fill out method. If key is not found in store, return ErrKeyNotFound
	// If Group is not authorized to access item, return ErrUnauthorized.
	// Don't forget to lock mutex!

	return nil, nil
}

// GracefulStop stops the gRPC server. Useful for testing
func (s *Server) GracefulStop() {
	s.server.GracefulStop()
}
