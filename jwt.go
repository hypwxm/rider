package rider

import (
	"net/http"
	"time"

	"github.com/hypwxm/rider/jwt"

	jwtgo "github.com/hypwxm/jwt-go"
)

type riderJwter struct {
	jwt         *jwt.Jwter
	context     Context
	expires     time.Duration
	tokenKey    string
	riderCookie *RiderCookie
}

type RiderCookie struct {
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

func RiderJwt(tokenKey string, secret string, expires time.Duration, riderCookie *RiderCookie) HandlerFunc {
	return func(c Context) {
		rj := &riderJwter{
			context:     c,
			expires:     expires,
			tokenKey:    tokenKey,
			riderCookie: riderCookie,
		}
		c.setJwt(rj)

		var claims jwtgo.MapClaims
		var err error
		var token string
		// 如果是app进行的请求，token回放在请求头里面，headers的token优先级大于cookie，所以先验证headers
		if token = c.HeaderValue(tokenKey); token != "" {
			claims, err = jwt.ValidateToken(token, secret)
		} else if token, err = c.CookieValue(tokenKey); err == nil {
			//如果cookie里面存在token，验证token
			claims, err = jwt.ValidateToken(token, secret)
		}

		if err == nil {
			//token通过验证
			//这里即使初始化了expires，但是不set，delete，对token重新赋值，expires不会起作用
			rj.jwt = jwt.NewJWTer(secret, expires)
			rj.jwt.TokenString = token
			rj.jwt.Claims = claims
			c.Next()
			return
		}

		//请求中的cookie不存在token，或者token验证不通过
		rj.jwt = jwt.NewJWTer(secret, expires)
		_, err = rj.jwt.CreateJwt(nil)
		if err != nil {
			c.Next(NError{500, err.Error()})
			return
		}
		c.SetCookie(http.Cookie{
			Name:     tokenKey,
			Value:    rj.jwt.TokenString,
			MaxAge:   int(expires),
			HttpOnly: riderCookie.HttpOnly,
			Secure:   riderCookie.Secure,
			SameSite: riderCookie.SameSite,
			Raw:      riderCookie.Raw,
			Unparsed: riderCookie.Unparsed,
		})
		c.Next()
	}
}

//发送tokenCookie
func (rj *riderJwter) SetTokenCookie(claims jwtgo.MapClaims) (string, error) {
	tokenString, err := rj.jwt.CreateJwt(rj.jwt.Claims)
	if err != nil {
		return "", err
	}
	rj.context.SetCookie(http.Cookie{
		Name:     rj.tokenKey,
		Value:    tokenString,
		MaxAge:   int(rj.expires),
		HttpOnly: rj.riderCookie.HttpOnly,
		Secure:   rj.riderCookie.Secure,
		SameSite: rj.riderCookie.SameSite,
		Raw:      rj.riderCookie.Raw,
		Unparsed: rj.riderCookie.Unparsed,
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
func (rj *riderJwter) Values() (jwtgo.MapClaims, error) {
	return rj.jwt.GetTokenClaims()
}

// 获取claims中指定字段的值
// token验证已在请求进来时RiderJwt中进行验证了，所以这里可以直接取值
func (rj *riderJwter) Get(key string) interface{} {
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

// 获取riderjwt上的token
func (rj *riderJwter) GetToken() string {
	return rj.jwt.TokenString
}
