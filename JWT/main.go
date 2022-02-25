package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

const (
	SECRETKEY = "243223ffslsfsldfl412fdsfsdf"
	MAXAGE    = 60 * 60 * 24 // 1天
)

type CustomClaims struct {
	UserId int64
	jwt.StandardClaims
}

func main() {
	c := CreateToken(11)
	token := c.Encrypt()
	fmt.Println(token)
	c, err := ParseToken(token)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(c)
}

func CreateToken(Userid int64) *CustomClaims {
	return &CustomClaims{
		UserId: Userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(MAXAGE) * time.Second).Unix(),
		},
	}
}

func (c *CustomClaims) Encrypt() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		log.Panicln(err)
	}
	return tokenString
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})
	clamis, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return clamis, nil
	}
	return nil, err
}
