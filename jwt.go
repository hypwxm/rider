package rider

import (
	"github.com/hypwxm/rider/jwt"
	"time"
	"net/http"
	jwtgo "github.com/dgrijalva/jwt-go"

)

type riderJwter struct {
	jwt *jwt.Jwter
	context Context
	expires time.Duration
}

func RiderJwt(secret string, expires time.Duration) HandlerFunc {
	return func(c Context) {
		rj := &riderJwter{context: c, expires: expires}
		c.setJwt(rj)
		if token, err := c.CookieValue("token"); err == nil {
			//如果cookie里面存在token，验证token
			claims, err := jwt.ValidateToken(token, secret)
			if err == nil {
				//token通过验证
				//这里即使初始化了expires，但是不set，delete，对token重新赋值，expires不会起作用
				rj.jwt = jwt.NewJWTer(secret, expires)
				rj.jwt.TokenString = token
				rj.jwt.Claims = claims
				c.Next()
				return
			}
		}
		//请求中的cookie不存在token，或者token验证不通过
		rj.jwt = jwt.NewJWTer(secret, expires)
		_, err := rj.jwt.CreateJwt(nil)
		if err != nil {
			c.Next(err)
			return
		}
		c.SetCookie(http.Cookie{
			Name:     "token",
			Value:    rj.jwt.TokenString,
			MaxAge:   int(expires),
			HttpOnly: true,
		})
		c.Next()
	}
}

//发送tokenCookie
func(rj *riderJwter) SetTokenCookie(claims jwtgo.MapClaims) (string, error) {
	tokenString, err := rj.jwt.CreateJwt(rj.jwt.Claims)
	if err != nil {
		return "", err
	}
	rj.context.SetCookie(http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   int(rj.expires),
		HttpOnly: true,
	})
	return tokenString, nil
}

func (rj *riderJwter) Set(key string, value interface{}) (string, error) {
	rj.jwt.Claims[key] = value
	return rj.SetTokenCookie(rj.jwt.Claims)
}

//删除token属性和值
func (rj *riderJwter) Delete(key string) (string, error) {
	if _, ok := rj.jwt.Claims[key]; ok {
		delete(rj.jwt.Claims, key)
	}
	return rj.SetTokenCookie(rj.jwt.Claims)
}

//获取token中的信息，payload
func (rj *riderJwter) Claims() (jwtgo.MapClaims, error) {
	return rj.jwt.GetTokenClaims()
}

//获取claims中指定字段的值
func (rj *riderJwter) ClaimsValue(key string) interface{} {
	if rj.jwt.Claims == nil {
		return nil
	}
	if v, ok := rj.jwt.Claims[key]; ok {
		return v
	}
	return nil
}

//删除claims的所有参数
func (rj *riderJwter) DeleteAll() (string, error) {
	rj.jwt.Claims = make(jwtgo.MapClaims)
	return rj.SetTokenCookie(rj.jwt.Claims)
}