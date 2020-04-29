package model

import (
	"encoding/json"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/mysql"
	"log"
)

type ClientDetails struct {
	// client 的标识
	ClientId string
	// client 的密钥
	ClientSecret string
	// 访问令牌有效时间，秒
	AccessTokenValiditySeconds int
	// 刷新令牌有效时间，秒
	RefreshTokenValiditySeconds int
	// 重定向地址，授权码类型中使用
	RegisteredRedirectUri string
	// 可以使用的授权类型
	AuthorizedGrantTypes []string
}

func (clientDetails *ClientDetails) IsMatch(clientId string, clientSecret string) bool {
	return clientId == clientDetails.ClientId && clientSecret == clientDetails.ClientSecret
}

type ClientDetailsModel struct {

}


func NewClientDetailsModel() *ClientDetailsModel {
	return &ClientDetailsModel{}
}

func (p *ClientDetailsModel) getTableName() string {
	return "client_details"
}


func (p *ClientDetailsModel) GetClientDetailsByClientId(clientId string) (*ClientDetails,  error)  {

	conn := mysql.DB()
	if result, err := conn.Table(p.getTableName()).Where(map[string]interface{}{"client_id": clientId}).First(); err == nil{

		var authorizedGrantTypes []string
		_ = json.Unmarshal([]byte(result["authorized_grant_types"].(string)), &authorizedGrantTypes)

		return &ClientDetails{
			ClientId:                   result["client_id"].(string),
			ClientSecret:               result["client_secret"].(string),
			AccessTokenValiditySeconds: int(result["access_token_validity_seconds"].(int64)),
			RefreshTokenValiditySeconds:int(result["refresh_token_validity_seconds"].(int64)),
			RegisteredRedirectUri:result["registered_redirect_uri"].(string),
			AuthorizedGrantTypes:authorizedGrantTypes,
		}, nil

	}else {
		return nil, err
	}

}


func (p *ClientDetailsModel) CreateClientDetails(clientDetails *ClientDetails) error {
	conn := mysql.DB()

	grantTypeString, _ := json.Marshal(clientDetails.AuthorizedGrantTypes)
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"client_id":     clientDetails.ClientId,
		"client_secret":   clientDetails.ClientSecret,
		"access_token_validity_seconds":    clientDetails.AccessTokenValiditySeconds,
		"refresh_token_validity_seconds":         clientDetails.RegisteredRedirectUri,
		"registered_redirect_uri": clientDetails.RegisteredRedirectUri,
		"authorized_grant_types":grantTypeString,
	}).Insert()
	if err != nil {
		log.Printf("Error : %v", err)
		return err
	}
	return nil
}


