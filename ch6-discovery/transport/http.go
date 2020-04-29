package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	endpts "github.com/longjoy/micro-go-book/ch6-discovery/endpoint"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpts.DiscoveryEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/say-hello").Handler(kithttp.NewServer(
		endpoints.SayHelloEndpoint,
		decodeSayHelloRequest,
		encodeJsonResponse,
		options...,
	))

	r.Methods("GET").Path("/discovery").Handler(kithttp.NewServer(
		endpoints.DiscoveryEndpoint,
		decodeDiscoveryRequest,
		encodeJsonResponse,
		options...,
	))


	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		options...,
	))

	return r
}

// decodeSayHelloRequest decode request params to struct
func decodeSayHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpts.SayHelloRequest{}, nil
}


// decodeDiscoveryRequest decode request params to struct
func decodeDiscoveryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	serviceName := r.URL.Query().Get("serviceName")
	if serviceName == ""{
		return nil, ErrorBadRequest
	}
	return endpts.DiscoveryRequest{
		ServiceName:serviceName,
	}, nil
}


// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

// encodeJsonResponse encode response to return
func encodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}


// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

