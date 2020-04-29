package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/longjoy/micro-go-book/ch11-security/endpoint"
	"github.com/longjoy/micro-go-book/ch11-security/service"
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
func MakeHttpHandler(ctx context.Context, endpoints endpoint.OAuth2Endpoints, tokenService service.TokenService, clientService service.ClientDetailsService, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}
	r.Path("/metrics").Handler(promhttp.Handler())

	clientAuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeClientAuthorizationContext(clientService, logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
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

	oauth2AuthorizationOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(makeOAuth2AuthorizationContext(tokenService, logger)),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

    r.Methods("Get").Path("/simple").Handler(kithttp.NewServer(
    	endpoints.SimpleEndpoint,
    	decodeSimpleRequest,
    	encodeJsonResponse,
    	oauth2AuthorizationOptions...,
    	))

	r.Methods("Get").Path("/admin").Handler(kithttp.NewServer(
		endpoints.AdminEndpoint,
		decodeAdminRequest,
		encodeJsonResponse,
		oauth2AuthorizationOptions...,
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

func makeOAuth2AuthorizationContext(tokenService service.TokenService, logger log.Logger) kithttp.RequestFunc {

	return func(ctx context.Context, r *http.Request) context.Context {

		// 获取访问令牌
		accessTokenValue := r.Header.Get("Authorization")
		var err error
		if accessTokenValue != ""{
			// 获取令牌对应的用户信息和客户端信息
			oauth2Details, err := tokenService.GetOAuth2DetailsByAccessToken(accessTokenValue)
			if err == nil {
				return context.WithValue(ctx, endpoint.OAuth2DetailsKey, oauth2Details)
			}
		}else {
			err = ErrorTokenRequest
		}

		return context.WithValue(ctx, endpoint.OAuth2ErrorKey, err)
	}
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

func decodeSimpleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &endpoint.SimpleRequest{}, nil
}

func decodeAdminRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &endpoint.AdminRequest{}, nil
}

func encodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}


// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthRequest{}, nil
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

