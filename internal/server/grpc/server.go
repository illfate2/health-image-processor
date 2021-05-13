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

func (s *Server) UserBlinked(req proto.Health_UserBlinkedServer) error {
	for {
		_, err := req.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		select {
		case <-req.Context().Done():
			return nil
		default:
			s.service.Blinked()
		}
	}
}

func (s *Server) ShoulderChangeAngle(req proto.Health_ShoulderChangeAngleServer) error {
	for {
		_, err := req.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		select {
		case <-req.Context().Done():
			return nil
		default:
			// TODO
		}
	}
}

func (s *Server) NoseChangeAngle(req proto.Health_NoseChangeAngleServer) error {
	for {
		_, err := req.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		select {
		case <-req.Context().Done():
			return nil
		default:
			// TODO
		}
	}
}
