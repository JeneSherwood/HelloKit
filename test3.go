package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func (receiver concatRequest) name() {

}

// concatRequest is a simple struct to represent the incoming request
type concatRequest struct {
	Text1 string `json:"text1"`
	Text2 string `json:"text2"`
	A     string
}

// concatResponse is a simple struct to represent the outgoing response
type concatResponse struct {
	ConcatenatedText string `json:"concatenatedText"`
}

// ConcatService is an interface that defines the concatenation logic
type ConcatService interface {
	Concat(text1, text2 string) string
}

// concatServiceImpl implements ConcatService interface
type concatServiceImpl struct{}

// Concat implements ConcatService interface
func (s concatServiceImpl) Concat(text1, text2 string) string {
	return text1 + text2
}

// MakeConcatEndpoint returns an endpoint that concatenates two strings
func MakeConcatEndpoint(svc ConcatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(concatRequest)
		concatenatedText := svc.Concat(req.Text1, req.Text2)
		return concatResponse{ConcatenatedText: concatenatedText}, nil
	}
}

func main() {
	svc := concatServiceImpl{}
	endpoints := Endpoints{
		ConcatEndpoint: MakeConcatEndpoint(svc),
	}
	// Now endpoints.ConcatEndpoint is a valid endpoint.
	// You can use it in your server or client.
	_ = endpoints
	return
}
