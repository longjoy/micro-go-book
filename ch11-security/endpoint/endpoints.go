package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-book/ch11-security/model"
	"github.com/longjoy/micro-go-book/ch11-security/service"
	"net/http"
)

// CalculateEndpoint define endpoint
type OAuth2Endpoints struct {
	TokenEndpoint		endpoint.Endpoint
	CheckTokenEndpoint	endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
	SimpleEndpoint 		endpoint.Endpoint
	AdminEndpoint		endpoint.Endpoint
}



func MakeClientAuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok{
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails); !ok{
				return  nil, ErrInvalidClientRequest
			}
			return next(ctx, request)
		}
	}
}

func MakeOAuth2AuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok{
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); !ok{
				return  nil, ErrInvalidUserRequest
			}
			return next(ctx, request)
		}
	}
}
func MakeAuthorityAuthorizationMiddleware(authority string, logger log.Logger) endpoint.Middleware  {
	return func(next endpoint.Endpoint) endpoint.Endpoint {

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok{
				return nil, err
			}
			if details, ok := ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details); !ok{
				return  nil, ErrInvalidClientRequest
			}else {
				for _, value := range details.User.Authorities{
					if value == authority{
						return next(ctx, request)
					}
				}
				return nil, ErrNotPermit
			}
		}
	}
}

const (

	OAuth2DetailsKey       = "OAuth2Details"
	OAuth2ClientDetailsKey = "OAuth2ClientDetails"
	OAuth2ErrorKey         = "OAuth2Error"

)


var (
	ErrInvalidClientRequest = errors.New("invalid client message")
	ErrInvalidUserRequest = errors.New("invalid user message")
	ErrNotPermit = errors.New("not permit")
)



type TokenRequest struct {
	GrantType string
	Reader *http.Request
}


type TokenResponse struct {
	AccessToken *model.OAuth2Token `json:"access_token"`
	Error string `json:"error"`
}

//  make endpoint
func MakeTokenEndpoint(svc service.TokenGranter, clientService service.ClientDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*TokenRequest)
		token, err := svc.Grant(ctx, req.GrantType, ctx.Value(OAuth2ClientDetailsKey).(*model.ClientDetails), req.Reader)
		var errString = ""
		if err != nil{
			errString = err.Error()
		}

		return TokenResponse{
			AccessToken:token,
			Error:errString,
		}, nil
	}
}


type CheckTokenRequest struct {
	Token string
	ClientDetails model.ClientDetails
}

type CheckTokenResponse struct {
	OAuthDetails *model.OAuth2Details `json:"o_auth_details"`
	Error string `json:"error"`

}

func MakeCheckTokenEndpoint(svc service.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CheckTokenRequest)
		tokenDetails, err := svc.GetOAuth2DetailsByAccessToken(req.Token)

		var errString = ""
		if err != nil{
			errString = err.Error()
		}

		return CheckTokenResponse{
			OAuthDetails:tokenDetails,
			Error:errString,
		}, nil
	}
}

type SimpleRequest struct {
}

type SimpleResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

type AdminRequest struct {
}

type AdminResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

func MakeSimpleEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result := svc.SimpleData(ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details).User.Username)
		return &SimpleResponse{
			Result:result,
		}, nil
	}

}

func MakeAdminEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result := svc.AdminData(ctx.Value(OAuth2DetailsKey).(*model.OAuth2Details).User.Username)
		return &AdminResponse{
			Result:result,
		}, nil
	}
}




// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{
			Status:status,
		}, nil
	}
}
