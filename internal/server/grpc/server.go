package grpc

import (
	"io"

	"github.com/illfate2/health-image-processor/internal/body"
	"github.com/illfate2/health-image-processor/proto"
)

type Server struct {
	proto.UnimplementedHealthServer
	service *body.Processor
}

func NewServer(service *body.Processor) *Server {
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

func (s *Server) ShouldersPositionChange(req proto.Health_ShouldersPositionChangeServer) error {
	for {
		msg, err := req.Recv()
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
			s.service.BackCrooked(msg.IsCrooked)
		}
	}
}

func (s *Server) NosePositionChange(req proto.Health_NosePositionChangeServer) error {
	for {
		msg, err := req.Recv()
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
			s.service.NoseCrooked(msg.IsCrooked)
		}
	}
}
