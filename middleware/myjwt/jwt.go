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

// SignKeys 不同token的签名密钥
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

	StatusAccessTokenExpired  = 1
	StatusRefreshTokenExpired = 2
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    "AuthHeaderNotFound",
			})
			c.Abort()
			return
		}
		// 将authHeader分割，如果不符合Bearer Auth则丢弃该请求
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    "AuthHeaderMalformed",
			})
			c.Abort()
			return
		}
		// 解析token信息
		claims, types, err := parseToken(parts[1])
		if err == nil {
			switch types {
			case 1:
				// 继续交由下一个路由处理,并将解析出的信息传递下去
				c.Set("TokenClaims", claims)
				c.Set("TokenType", types)
				return
			case 2:
				token, innerErr := GetAccessToken(claims.User)
				if innerErr != nil {
					c.JSON(http.StatusOK, gin.H{
						"status": -1,
						"msg":    innerErr.Error(),
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"status": 0,
						"msg":    "RefreshSucess",
						"data":   token,
					})
				}
				c.Abort()
				return
			}
		} else {
			// 如果AccessToken过期，则返回状态码498提醒更新
			if types == 1 && err == ErrTokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": StatusAccessTokenExpired,
					"msg":    err.Error(),
				})
				c.Abort()
				return
			} else if types == 2 && err == ErrTokenExpired {
				// 如果RefreshTokenToken过期，则返回状态码499提醒更新
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": StatusRefreshTokenExpired,
					"msg":    err.Error(),
				})
				c.Abort()
				return
			} else {
				// 其他错误返回err字符串
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": -1,
					"msg":    err.Error(),
				})
				c.Abort()
				return
			}
		}

	}
}

// parseToken用于解析Token
func parseToken(tokenString string) (*UserClaims, int, error) {
	var types int
	claims := new(UserClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 判断Token的类型，根据类型返回不同的signkey
		t, ok := token.Header["type"]
		if !ok {
			return nil, errors.New("TokenTypeEmpty")
		}
		typeFloat, ok := t.(float64)
		if ok {
			switch int(typeFloat) {
			case 1:
				types = 1
				return keys.AccessTokenSignKey, nil
			case 2:
				types = 2
				return keys.RefreshTokenSignKey, nil
			default:
				types = -1
				return nil, errors.New("TokenTypeError")
			}
		} else {
			types = -1
			return nil, errors.New("TokenTypeError")
		}
	})
	// 返回err时进行细化
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			switch v.Errors {
			case jwt.ValidationErrorMalformed:
				return nil, types, ErrTokenMalformed
			case jwt.ValidationErrorExpired:
				return nil, types, ErrTokenExpired
			case jwt.ValidationErrorNotValidYet:
				return nil, types, ErrTokenNotValidYet
			default:
				return nil, types, errors.New(v.Inner.Error())
			}
		}
	}
	// 如果token合法，则返回claims，token种类（err为nil）
	if token.Valid {
		return claims, types, nil
	} else {
		return nil, types, ErrTokenInvalid
	}
}

// GetAccessToken 生成一个token
func GetAccessToken(user string) (string, error) {
	claims := UserClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["type"] = 1
	return token.SignedString(keys.AccessTokenSignKey)
}

// GetAccessToken 生成一个token
func GetRefreshToken(user string) (string, error) {
	claims := UserClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["type"] = 2
	return token.SignedString(keys.RefreshTokenSignKey)
}
