package model

import "time"

type OAuth2Token struct {

	// 刷新令牌
	RefreshToken *OAuth2Token
	// 令牌类型
	TokenType string
	// 令牌
	TokenValue string
	// 过期时间
	ExpiresTime *time.Time

}


func (oauth2Token *OAuth2Token) IsExpired() bool  {
	return oauth2Token.ExpiresTime != nil &&
		oauth2Token.ExpiresTime.Before(time.Now())
}

type OAuth2Details struct {
	Client *ClientDetails
	User *UserDetails
}