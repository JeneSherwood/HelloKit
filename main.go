package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Service defines the service business logic
type Service interface {
	Concat(string, string) string
}

// stringService implements Service interface
type stringService struct{}

func (stringService) Concat(a, b string) string {
	return a + b
}

// Endpoints defines the service endpoints
type Endpoints struct {
	ConcatEndpoint endpoint.Endpoint
}

// MakeConcatEndpoint creates the Concat endpoint
func MakeConcatEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(concatRequest)
		v := svc.Concat(req.A, req.B)
		return concatResponse{V: v}, nil
	}
}

// concatRequest and concatResponse define the request and response structures
type concatRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}

type concatResponse struct {
	V string `json:"v"`
}

// decodeConcatRequest decodes the HTTP request
func decodeConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request concatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// encodeResponse encodes the HTTP response
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	var svc Service
	svc = stringService{}
	endpoints := Endpoints{
		ConcatEndpoint: MakeConcatEndpoint(svc),
	}

	handler := httptransport.NewServer(
		endpoints.ConcatEndpoint,
		decodeConcatRequest,
		encodeResponse,
	)

	http.Handle("/concat", handler)
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log(http.ListenAndServe(":8080", nil))
}
