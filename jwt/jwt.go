package jwt

//信息加密
//信息解密
//信息认证
//刷新token
//http接受认证信息
//http发送加密token
/*
type StandardClaims struct {
    Audience  string `json:"aud,omitempty"`
    ExpiresAt int64  `json:"exp,omitempty"`
    Id        string `json:"jti,omitempty"`
    IssuedAt  int64  `json:"iat,omitempty"`
    Issuer    string `json:"iss,omitempty"`
    NotBefore int64  `json:"nbf,omitempty"`
    Subject   string `json:"sub,omitempty"`
}

iss(Issuser)：代表这个JWT的签发主体；

sub(Subject)：代表这个JWT的主体，即它的所有人；

aud(Audience)：代表这个JWT的接收对象；

exp(Expiration time)：是一个时间戳，代表这个JWT的过期时间；

nbf(Not Before)：是一个时间戳，代表这个JWT生效的开始时间，意味着在这个时间之前验证JWT是会失败的；

iat(Issued at)：是一个时间戳，代表这个JWT的签发时间；

jti(JWT ID)：是JWT的唯一标识。
*/

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"errors"
	"time"
)

type Jwter struct {
	SecretKey string
	TokenString string
	Expires time.Time
	Claims jwt.MapClaims
}

func NewJWTer(secretKey string, expires time.Duration) *Jwter {
	j := &Jwter{
		SecretKey: secretKey,
		Expires: time.Now().Add(expires),
		Claims: make(map[string]interface{}),
	}
	return j
}

//data中的expires字段会被忽略，使用初始化时设置的Expires
func (j *Jwter) CreateJwt(data map[string]interface{}) (string, error) {
	if j.Claims == nil {
		j.Claims = make(jwt.MapClaims)
	}
	for k, v := range data {
		j.Claims[k] = v
	}
	//设置过期时间
	j.Claims["exp"] = j.Expires
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.Claims)
	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		fmt.Println("Error while signing the token")
		return "", err
	}
	j.TokenString = tokenString
	return tokenString, nil
}

//设置token属性和值
func (j *Jwter) Set(key string, value interface{}) (string, error) {
	j.Claims[key] = value
	return j.CreateJwt(j.Claims)
}

//删除token属性和值
func (j *Jwter) Delete(key string) (string, error) {
	if _, ok := j.Claims[key]; ok {
		delete(j.Claims, key)
	}
	return j.CreateJwt(j.Claims)
}

//验证客服端token可用性
func ValidateToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if !token.Valid {
		return nil, errors.New("invalid token")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("That's not even a token. ")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return nil, errors.New("Token is either expired or not active yet. ")
		} else {
			return nil, errors.New("Couldn't handle this token:" + err.Error())
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

//获取claims里面的参数
func (j *Jwter) GetTokenClaims() (jwt.MapClaims, error) {
	claims, err := ValidateToken(j.TokenString, j.SecretKey)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
