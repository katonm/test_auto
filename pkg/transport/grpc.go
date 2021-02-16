package transport

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"

	"my/sberAuto/test_auto/pb"
	"my/sberAuto/test_auto/pkg/endpoint"
)

type gRPCServer struct {
	validate grpc.Handler
	fix      grpc.Handler
	pb.ParenthesesServiceServer
}

// NewGRPCServer initializes a new gRPC server...
func NewGRPCServer(endpoints endpoint.Endpoints, logger log.Logger) pb.ParenthesesServiceServer {
	return &gRPCServer{
		validate: grpc.NewServer(
			endpoints.Validate,
			decodeValidateRequest,
			encodeValidateResponse,
		),
		fix: grpc.NewServer(
			endpoints.Fix,
			decodeFixRequest,
			encodeFixResponse,
		),
	}
}

func (s *gRPCServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	_, resp, err := s.validate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ValidateResponse), nil
}

func decodeValidateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ValidateRequest)

	return endpoint.ValidateReq{StrIn: req.StrIn}, nil
}

func encodeValidateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.ValidateResp)

	return &pb.ValidateResponse{Result: resp.Result}, nil
}

func (s *gRPCServer) Fix(ctx context.Context, req *pb.FixRequest) (*pb.FixResponse, error) {
	_, resp, err := s.fix.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FixResponse), nil
}

func decodeFixRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.FixRequest)

	return endpoint.FixReq{StrIn: req.StrIn}, nil
}

func encodeFixResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.FixResp)

	return &pb.FixResponse{StrOut: resp.StrOut}, nil
}
