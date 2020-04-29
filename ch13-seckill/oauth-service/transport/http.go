package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/longjoy/micro-go-book/ch13-seckill/oauth-service/endpoint"
	"github.com/longjoy/micro-go-book/ch13-seckill/oauth-service/service"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
	ErrorGrantTypeRequest = errors.New("invalid request grant type")
	ErrorTokenRequest = errors.New("invalid request token")
	ErrInvalidClientRequest = errors.New("invalid client message")

)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpoint.OAuth2Endpoints, tokenService service.TokenService, clientService service.ClientDetailsService,zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}
	r.Path("/metrics").Handler(promhttp.Handler())


	clientAuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeClientAuthorizationContext(clientService, logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}


	r.Methods("POST").Path("/oauth/token").Handler(kithttp.NewServer(
		endpoints.TokenEndpoint,
		decodeTokenRequest,
		encodeJsonResponse,
		clientAuthorizationOptions...,
	))

	r.Methods("POST").Path("/oauth/check_token").Handler(kithttp.NewServer(
		endpoints.CheckTokenEndpoint,
		decodeCheckTokenRequest,
		encodeJsonResponse,
		clientAuthorizationOptions...,
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


func makeClientAuthorizationContext(clientDetailsService service.ClientDetailsService, logger log.Logger) kithttp.RequestFunc {

	return func(ctx context.Context, r *http.Request) context.Context {

		if clientId, clientSecret, ok := r.BasicAuth(); ok {
			clientDetails, err := clientDetailsService.GetClientDetailByClientId(ctx, clientId, clientSecret)
			if err == nil {
				return context.WithValue(ctx, endpoint.OAuth2ClientDetailsKey, clientDetails)
			}
		}
		return context.WithValue(ctx, endpoint.OAuth2ErrorKey, ErrInvalidClientRequest)
	}
}


func decodeTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	grantType := r.URL.Query().Get("grant_type")
	if grantType == ""{
		return nil, ErrorGrantTypeRequest
	}
	return &endpoint.TokenRequest{
		GrantType:grantType,
		Reader:r,
	}, nil

}

func decodeCheckTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	tokenValue := r.URL.Query().Get("token")
	if tokenValue == ""{
		return nil, ErrorTokenRequest
	}

	return &endpoint.CheckTokenRequest{
		Token:tokenValue,
	}, nil

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


func encodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}


// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthRequest{}, nil
}
