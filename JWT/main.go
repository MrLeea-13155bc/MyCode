package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

const (
	SECRETKEY = "54188wozhenweida"
	MAXAGE    = 3 // 1天
)

type CustomClaims struct {
	UserId int64
	jwt.StandardClaims
}

func main() {
	token := CreateToken(11)
	time.Sleep(time.Second * 5)
	c, err := ParseToken(token)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(c)
}

func CreateToken(Userid int64) string {
	c := &CustomClaims{
		UserId: Userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(MAXAGE) * time.Second).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(tokenString)
	return tokenString
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})
	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return claims, nil
	}
	t := time.Unix(claims.StandardClaims.ExpiresAt, 0)
	timeExceed := int(time.Now().Sub(t).Seconds())
	if timeExceed < MAXAGE {
		CreateToken(claims.UserId)
	}
	return nil, err
}
