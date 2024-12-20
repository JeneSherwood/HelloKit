package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"testing"
)

func MakeConcatEndpointTest(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(concatRequest)
		v := svc.Concat(req.A, req.B)
		return concatResponse{V: v}, nil
	}
}

type concatService struct {
}

func (c concatService) Concat(s string, s2 string) string {
	//TODO implement me
	panic("implement me")
}

func TestConcatEndpoint(t *testing.T) {
	svc := &concatService{}
	ep := MakeConcatEndpointTest(svc)

	request := concatRequest{A: string(10), B: string(20)}
	_, err := ep(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("concat endpoint test passed")
	t.Log(svc.Concat("hello", "world"))
}
