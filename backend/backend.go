package backend

import (
	"log"
	"net"

	"github.com/microapis/messages-api"
	"github.com/microapis/messages-api/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ListenAndServe ...
func ListenAndServe(addr string, backend messages.Backend) error {
	lis, err := net.Listen("tcp", string(addr))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	proto.RegisterMessageBackendServiceServer(s, &service{backend})

	return s.Serve(lis)
}

type service struct {
	backend messages.Backend
}

func (s *service) Approve(ctx context.Context, r *proto.MessageBackendApproveRequest) (*proto.MessageBackendApproveResponse, error) {
	var resp proto.MessageBackendApproveResponse
	var err error

	resp.Valid, err = s.backend.Approve(r.Content)
	if err != nil {
		resp.Error = &proto.MessagesError{
			Code:    500,
			Message: err.Error(),
		}
	}
	return &resp, nil
}

func (s *service) Deliver(ctx context.Context, r *proto.MessageBackendDeliverRequest) (*proto.MessageBackendDeliverResponse, error) {
	var resp proto.MessageBackendDeliverResponse
	if err := s.backend.Deliver(r.Content); err != nil {
		resp.Error = &proto.MessagesError{
			Code:    500,
			Message: err.Error(),
		}
	}
	return &resp, nil
}
