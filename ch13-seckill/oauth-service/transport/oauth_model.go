package transport

import (
	"context"
	"github.com/longjoy/micro-go-book/ch13-seckill/oauth-service/endpoint"
	"github.com/longjoy/micro-go-book/ch13-seckill/oauth-service/model"
	"github.com/longjoy/micro-go-book/ch13-seckill/pb"
)

func EncodeGRPCCheckTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.CheckTokenRequest)
	return &pb.CheckTokenRequest{
		Token: req.Token,
	}, nil
}

func DecodeGRPCCheckTokenRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CheckTokenRequest)
	return &endpoint.CheckTokenRequest{
		Token: req.Token,
	}, nil
}

func EncodeGRPCCheckTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.CheckTokenResponse)

	if resp.Error != "" {
		return &pb.CheckTokenResponse{
			IsValidToken: false,
			Err:          resp.Error,
		}, nil
	} else {
		return &pb.CheckTokenResponse{
			UserDetails: &pb.UserDetails{
				UserId:      resp.OAuthDetails.User.UserId,
				Username:    resp.OAuthDetails.User.Username,
				Authorities: resp.OAuthDetails.User.Authorities,
			},
			ClientDetails: &pb.ClientDetails{
				ClientId:                    resp.OAuthDetails.Client.ClientId,
				AccessTokenValiditySeconds:  int32(resp.OAuthDetails.Client.AccessTokenValiditySeconds),
				RefreshTokenValiditySeconds: int32(resp.OAuthDetails.Client.RefreshTokenValiditySeconds),
				AuthorizedGrantTypes:        resp.OAuthDetails.Client.AuthorizedGrantTypes,
			},
			IsValidToken: true,
			Err:          "",
		}, nil
	}
}

func DecodeGRPCCheckTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.CheckTokenResponse)
	if resp.Err != "" {
		return endpoint.CheckTokenResponse{
			Error: resp.Err,
		}, nil
	} else {
		return endpoint.CheckTokenResponse{
			OAuthDetails: &model.OAuth2Details{
				User: &model.UserDetails{
					UserId:      resp.UserDetails.UserId,
					Username:    resp.UserDetails.Username,
					Authorities: resp.UserDetails.Authorities,
				},
				Client: &model.ClientDetails{
					ClientId:                    resp.ClientDetails.ClientId,
					AccessTokenValiditySeconds:  int(resp.ClientDetails.AccessTokenValiditySeconds),
					RefreshTokenValiditySeconds: int(resp.ClientDetails.RefreshTokenValiditySeconds),
					AuthorizedGrantTypes:        resp.ClientDetails.AuthorizedGrantTypes,
				},
			},
		}, nil

	}

}
