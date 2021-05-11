package grpc

import (
	"io"

	"github.com/illfate2/health-image-processor/internal/blinker"
	"github.com/illfate2/health-image-processor/proto"
)

type Server struct {
	proto.UnimplementedHealthServer
	service *blinker.Service
}

func NewServer(service *blinker.Service) *Server {
	return &Server{service: service}
}

func (s *Server) UserBlinked(server proto.Health_UserBlinkedServer) error {
	for {
		_, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		select {
		case <-server.Context().Done():
			return nil
		default:
			s.service.Blinked()
		}
	}
}
