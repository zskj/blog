package util

import (
	"blog/models"
	"blog/pkg/setting"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	jwt.StandardClaims
}

//生成令牌
func GenerateToken(user models.User) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
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
	return tokenClaims.SignedString(jwtSecret)
}

//解析令牌
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

//刷新令牌
func RefreshToken(tokenString string) (string, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			id, err := strconv.Atoi(claims.Id)
			if err != nil {
				return "", err
			}
			user, userError := models.FindUserById(id)
			if userError != nil {
				return "", err
			}
			nowTime := time.Now()
			expireTime := nowTime.Add(1 * time.Hour)
			claims := Claims{
				jwt.StandardClaims{
					Audience:  user.Username,         // 受众
					ExpiresAt: expireTime.Unix(),     // 失效时间
					Id:        strconv.Itoa(user.ID), // 编号
					IssuedAt:  time.Now().Unix(),     // 签发时间
					Issuer:    "blog",                // 签发人
					NotBefore: time.Now().Unix(),     // 生效时间
					Subject:   "refresh",             // 场景
				},
			}
			tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			return tokenClaims.SignedString(jwtSecret)
		}
	}
	return "", err
}
