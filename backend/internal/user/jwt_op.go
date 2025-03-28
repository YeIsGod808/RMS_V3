package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/RMS_V3/pkg/errcode"

	"github.com/RMS_V3/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type userClaims struct {
	jwt.StandardClaims
	*User
}

func jwtGenerateToken(u *User, d time.Duration) (string, error) {
	u.Password = ""
	expireTime := time.Now().Add(d)
	stdClams := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        u.Id,
		Issuer:    config.GetGlobalConfig().JwtConfig.Issuer,
	}
	uClaims := userClaims{
		StandardClaims: stdClams,
		User:           u,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uClaims)
	return token.SignedString([]byte(config.GetGlobalConfig().JwtConfig.JwtSalt))
}

func JwtParseToken(token string) (*User, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	uClaims := userClaims{}
	_, err := jwt.ParseWithClaims(token, &uClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetGlobalConfig().JwtConfig.JwtSalt), nil
	})
	return uClaims.User, err
}

// parse token and confirm user_type,
//
// return false and response StatusUnauthorized if token invalid or user_type lower than expected.
//
// Admin > Teacher > Student
func CheckUserPermission(token string, userType UserTypes, c *gin.Context) (user *User, ok bool) {
	ok = false
	user, err := JwtParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"ret": errcode.JWT_ERR,
			"msg": "登录信息过期，请重新登录",
		})
		return
	}
	user.Password = ""
	
	ok = PermissionCmp(user.UserType, userType)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"ret": errcode.AUTH_ERR,
			"msg": "权限不足",
		})
	}
	return
}

// return true if a>=b
func PermissionCmp(a UserTypes, b UserTypes) bool {
	return typeVal[a] >= typeVal[b]
}
