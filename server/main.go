package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-ping/gen/proto"

	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type server struct {
	pb.UnimplementedPingServiceServer
}

func (s *server) Ping(ctx context.Context, in *pb.Empty) (*pb.PongResponse, error) {
	// Get client IP address
	clientIP := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
		// Extract just the IP part if it's in IP:port format
		if host, _, err := net.SplitHostPort(clientIP); err == nil {
			clientIP = host
		}
	}

	zerolog.Info().Str("client_ip", clientIP).Msg("📡 Ping API called - responding with Pong")
	return &pb.PongResponse{Message: "Pong"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPingServiceServer(s, &server{})

	fmt.Println("🚀 gRPC Ping server started on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
