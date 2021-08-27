package myjwt

import (
	"crypto/rand"
	"errors"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kascas/httpserver/logs"
)

// UserClaims Token的载荷，此处只有User一个字段
type UserClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

// KeyStruct 签名结构
type SignKeys struct {
	AccessTokenSignKey  []byte
	RefreshTokenSignKey []byte
}

var keys *SignKeys
var (
	ErrTokenExpired     = errors.New("TokenExpired")
	ErrTokenNotValidYet = errors.New("TokenNotValid")
	ErrTokenMalformed   = errors.New("TokenMalformed")
	ErrTokenInvalid     = errors.New("TokenInvalid")
)

func init() {
	keys = &SignKeys{
		AccessTokenSignKey:  RandomBytes(64),
		RefreshTokenSignKey: RandomBytes(64),
	}
}

func RandomBytes(size int) []byte {
	min, max := new(big.Int), new(big.Int)
	min.Lsh(big.NewInt(1), uint(size*8-1))
	max.Lsh(big.NewInt(1), uint(size*8))
	for {
		tmp, err := rand.Int(rand.Reader, max)
		if err != nil {
			logs.ErrorPanic(err, "JWTSecretInitFailed")
		}
		if tmp.Cmp(min) >= 0 {
			return tmp.Bytes()
		}
	}
}

// JWTAuth 中间件，检查Token是否合法
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中获取token内容
		authHeader := c.Request.Header.Get(`Authorization`)
		// 如果Token为空，则返回-1状态码
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "AuthHeaderNotFound",
			})
			c.Abort()
			return
		}
		// 将authHeader分割，如果不符合Bearer Auth则丢弃该请求
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "AuthHeaderMalformed",
			})
			c.Abort()
			return
		}
		// 解析token信息
		claims, err := parseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}

// parseToken用于解析Token，如果错误则返回（nil，err）
func parseToken(tokenString string) (*jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return keys.AccessTokenSignKey, nil
	})
	// 若干Error
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			switch v.Errors {
			case jwt.ValidationErrorMalformed:
				return nil, ErrTokenMalformed
			case jwt.ValidationErrorExpired:
				return nil, ErrTokenExpired
			case jwt.ValidationErrorNotValidYet:
				return nil, ErrTokenNotValidYet
			default:
				return nil, ErrTokenInvalid
			}
		}
	}
	// 如果token合法，则返回claims
	if token.Valid {
		return &token.Claims, nil
	} else {
		return nil, ErrTokenInvalid
	}
}

// GetAccessToken 生成一个token
func GetAccessToken(user string) (string, error) {
	claims := UserClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Add(-1 * time.Minute).Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "kascas",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(keys.AccessTokenSignKey)
}
