package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"my/sberAuto/test_auto/pkg/service"
)

// Endpoints struct holds the list of endpoints definition...
type Endpoints struct {
	Fix      endpoint.Endpoint
	Validate endpoint.Endpoint
}

// FixReq struct holds the endpoint request definition...
type FixReq struct {
	StrIn string
}

// FixResp struct holds the endpoint response definition...
type FixResp struct {
	StrOut string
}

// ValidateReq struct holds the endpoint request definition...
type ValidateReq struct {
	StrIn string
}

// ValidateResp struct holds the endpoint response definition...
type ValidateResp struct {
	Result string
}

// MakeEndpoints func initializes the Endpoint instances...
func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Fix:      makeFixEndpoint(s),
		Validate: makeValidateEndpoint(s),
	}
}

func makeValidateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ValidateReq)
		res := s.Validate(ctx, req.StrIn)

		return ValidateResp{Result: res}, nil
	}
}

func makeFixEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(FixReq)
		result := s.Fix(ctx, req.StrIn)

		return FixResp{StrOut: result}, nil
	}
}
