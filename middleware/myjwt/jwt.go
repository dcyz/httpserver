package myjwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 关于Token验证的若干Error
var (
	TokenExpired     = errors.New("Token已过期")
	TokenNotValidYet = errors.New("Token认证错误")
	TokenMalformed   = errors.New("Token格式错误")
	TokenInvalid     = errors.New("Token不合法")
	SignKey          = "httpserver"
)

// CustomClaims Token的载荷，此处只有User一个字段
type CustomClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// GetSignKey 获取SignKey
func GetSignKey() string {
	return SignKey
}

// SetSignKey 设置SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
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
				"msg":    "无Token信息",
			})
			c.Abort()
			return
		}
		// 将authHeader分割，如果不符合Bearer Auth则丢弃该请求
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "AuthHeader格式错误",
			})
			c.Abort()
			return
		}
		// 新建JWT实例
		j := &JWT{
			[]byte(GetSignKey()),
		}
		// 解析token信息
		claims, err := j.ParseToken(parts[1])
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

// ParseToken用于解析Token，如果错误则返回（nil，err）
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	// 若干Error
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors == jwt.ValidationErrorMalformed {
				return nil, TokenMalformed
			} else if v.Errors == jwt.ValidationErrorExpired {
				return nil, TokenExpired
			} else if v.Errors == jwt.ValidationErrorNotValidYet {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	// 如果token合法，则返回claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return ``, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return ``, TokenInvalid
}
