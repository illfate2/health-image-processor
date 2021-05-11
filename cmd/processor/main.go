package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/illfate2/health-image-processor/internal/blinker"
	"github.com/illfate2/health-image-processor/internal/image"
	"github.com/illfate2/health-image-processor/internal/server/grpc"
	"github.com/illfate2/health-image-processor/internal/server/ws"
	"github.com/illfate2/health-image-processor/proto"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcServerPortEnv = "GRPC_SERVER_PORT"
const httpServerPortEnv = "HTTP_SERVER_PORT"

func main() {
	service := blinker.NewService(4)
	processor := image.NewProcessor()
	wsServer := ws.NewServer(service, processor)
	go func() {
		log.Fatal(http.ListenAndServe(":"+os.Getenv(httpServerPortEnv), wsServer))
	}()
	grpcServer := grpclib.NewServer()
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", ":"+os.Getenv(grpcServerPortEnv))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	proto.RegisterHealthServer(grpcServer, grpc.NewServer(service))
	log.Fatal(grpcServer.Serve(lis))
}
