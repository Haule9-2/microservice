package main

import (
    "log"
    "net"
    "google.golang.org/grpc"
    "github.com/Haule9-2/microservice/adapter/userclient/generatedclient"
)

type server struct {
    generatedclient.UnimplementedUserServiceServer
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    s := grpc.NewServer()
    generatedclient.RegisterUserServiceServer(s, &server{})
    log.Println("gRPC server is running on port 50051...")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
