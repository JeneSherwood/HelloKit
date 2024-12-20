package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakeConcatEndpointTest(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(concatRequest)
		v := svc.Concat(req.A, req.B)
		return concatResponse{V: v}, nil
	}
}
