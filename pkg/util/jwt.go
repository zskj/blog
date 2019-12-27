package util

import (
	"blog/models"
	"blog/pkg/setting"
	jwt "github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username string, password string, user models.User) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			Audience:  user.Username,         // 受众
			ExpiresAt: expireTime.Unix(),     // 失效时间
			Id:        strconv.Itoa(user.ID), // 编号
			IssuedAt:  time.Now().Unix(),     // 签发时间
			Issuer:    "blog",                // 签发人
			NotBefore: time.Now().Unix(),     // 生效时间
			Subject:   "login",               // 场景
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
