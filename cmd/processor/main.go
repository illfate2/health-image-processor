package main

import (
	"log"
	"net"
	"net/http"

	"github.com/illfate2/health-image-processor/internal/blinker"
	"github.com/illfate2/health-image-processor/internal/image"
	"github.com/illfate2/health-image-processor/internal/server/grpc"
	"github.com/illfate2/health-image-processor/internal/server/ws"
	"github.com/illfate2/health-image-processor/proto"
	grpclib "google.golang.org/grpc"
)

func main() {
	service := blinker.NewService(4)
	processor := image.NewProcessor()
	wsServer := ws.NewServer(service, processor)
	go func() {
		log.Fatal(http.ListenAndServe(":9999", wsServer))
	}()
	grpcServer := grpclib.NewServer()
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	proto.RegisterHealthServer(grpcServer, grpc.NewServer(service))
	log.Fatal(grpcServer.Serve(lis))
}
