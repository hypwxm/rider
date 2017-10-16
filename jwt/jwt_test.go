package jwt

import (
	"testing"
	"time"
)

func TestNewJWTer(t *testing.T) {
	jwt := NewJWTer("rider", time.Hour)
	if jwt.SecretKey != "rider" {
		t.Error("secretKey != rider")
	}
	if time.Now().Add(time.Hour).Sub(jwt.Expires) > time.Second {
		t.Error("Expires != Expires", time.Now().Add(time.Hour).Sub(jwt.Expires))
	}
}


func TestJwter_CreateJwt(t *testing.T) {
	jwt := NewJWTer("rider", time.Hour)
	token, err := jwt.CreateJwt(nil)
	if err != nil {
		t.Error(err)
	}

	//测试ValidateToken，token加密能否通过，并且能返回exp
	claims, err := ValidateToken(token, "rider")
	if err != nil {
		t.Error(err)
	}
	if exp, ok := claims["exp"]; ok {
		if exp == "" {
			t.Error("you must set exp")
		} else {
			_, err := time.Parse(time.RFC3339Nano, exp.(string))
			if err != nil {
				t.Error("exp must a time, ", err)
			}
		}
	} else {
		t.Error("you must set exp")
	}

	//测试给jwt加参数Set()，验证Set后jwt通过tokenstring获取的参数是否正确的改变
	jwt.Set("test", "test")
	claims, err = jwt.GetTokenClaims()
	if len(claims) != 2 {
		t.Error("Set jwt error")
	} else {
		if test, ok := claims["test"]; ok {
			if test != "test" {
				t.Error("Set jwt error")
			}
		} else {
			t.Error("Set jwt error", claims)
		}
	}

	//测试Delete
	jwt.Delete("test")
	claims, err = jwt.GetTokenClaims()

	if len(claims) != 1 {
		t.Error("delete claims error")
	} else {
		if _, ok := claims["test"]; ok {
			t.Error("delete claims error")
		}
	}
}