package middleware

import (
	"errors"
	"github.com/RaymondCode/simple-demo/constants"
	"github.com/RaymondCode/simple-demo/structs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

type Claims struct {
	ID         uint
	Identifier string
	jwt.StandardClaims
}

func GenerateToken(uid uint, identifier string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour * 24)
	issuer := "douyin"
	claims := Claims{
		ID:         uid,
		Identifier: identifier,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(constants.JwtSignKey)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return constants.JwtSignKey, nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware 基于JWT的认证中间件--验证用户是否登录
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		if token == "" {
			c.JSON(http.StatusOK, structs.Response{
				StatusCode: 2003,
				StatusMsg:  "缺少token",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, structs.Response{
				StatusCode: 2005,
				StatusMsg:  "登录已过期",
			})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next() // 后续的处理函数可以用过c.Get("user")来获取当前请求的用户信息
	}
}
