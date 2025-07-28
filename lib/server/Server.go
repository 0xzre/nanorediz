package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/raft"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/app"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/config"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type GRPCServer struct {
	address    *Address
	grpcServer *grpc.Server
	listener   net.Listener
	config     *config.Config
}

func NewServer(_address *Address, isContact bool, contactAddress *Address) *GRPCServer {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	server := &GRPCServer{
		address: _address,
		config:  cfg,
	}

	netListen, err := net.Listen("tcp", server.address.ToString())
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", server.address.ToString(), err)
	}
	server.listener = netListen

	server.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(server.unaryInterceptor),
		grpc.MaxConcurrentStreams(uint32(cfg.Server.MaxConcurrentRPC)),
	)

	app := app.NewKVStore()
	raft := raft.NewRaftNode(app, server.address, isContact, contactAddress)

	kvservice := service.NewKVService(raft)
	raftservice := service.NewRaftService(raft)
	pb.RegisterKeyValueServiceServer(server.grpcServer, kvservice)
	pb.RegisterRaftServiceServer(server.grpcServer, raftservice)

	return server
}

func (server *GRPCServer) Serve() {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the gRPC server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Printf("Starting nanorediz server on %s", server.address.ToString())
		if err := server.grpcServer.Serve(server.listener); err != nil {
			errChan <- err
		}
	}()

	// Wait for either an error or a signal
	select {
	case err := <-errChan:
		log.Fatalf("Server failed to start: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		server.Shutdown()
	}
}

func (server *GRPCServer) Shutdown() {
	log.Println("Shutting down server...")
	
	// Create a timeout context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), server.config.Server.ShutdownTimeout)
	defer cancel()

	// Channel to track graceful shutdown completion
	done := make(chan struct{})

	go func() {
		server.grpcServer.GracefulStop()
		close(done)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-done:
		log.Println("Server shut down gracefully")
	case <-ctx.Done():
		log.Println("Shutdown timeout exceeded, forcing stop")
		server.grpcServer.Stop()
	}
}

func (server *GRPCServer) unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Add timeout context
	ctx, cancel := context.WithTimeout(ctx, server.config.Server.GrpcTimeout)
	defer cancel()

	// Log method name and request
	p, _ := peer.FromContext(ctx)
	log.Printf("Method: %s, Request: %+v, From: %v", info.FullMethod, req, p.Addr)

	// Handling request
	resp, err := handler(ctx, req)

	// Log response and elapsed time
	log.Printf("Method: %s, Response: %+v, ElapsedTime: %s, Error: %v\n\n", info.FullMethod, resp, time.Since(start), err)

	return resp, err
}
