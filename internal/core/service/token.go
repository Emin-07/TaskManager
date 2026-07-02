package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

type TokenServ struct {
	PrivKeyPath string
	PubKeyPath  string
}

func NewTokenService(privKeyPath, pubKeyPath string) TokenServ {
	return TokenServ{
		PrivKeyPath: privKeyPath,
		PubKeyPath:  pubKeyPath,
	}
}

//func (ts TokenServ) KeyFunc(token *jwt.Token) (any, error) {
//	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
//		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//	}
//
//	verifyBytes, err := os.ReadFile(ts.PubKeyPath)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
//}

type CustomClaims struct {
	jwt.RegisteredClaims
	Role string
}

func (ts TokenServ) CreateToken(id, role string) (string, error) {
	signBytes, err := os.ReadFile(ts.PrivKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256,
		&CustomClaims{
			jwt.RegisteredClaims{Subject: id, ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * 15))}, strings.ToLower(role),
		})
	return t.SignedString(signKey)
}

func (ts TokenServ) ParseFromRequest(r *http.Request) (map[string]string, error) {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		verifyBytes, err := os.ReadFile(ts.PubKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		return jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	}, request.WithClaims(&jwt.RegisteredClaims{}))
	if err != nil {
		return nil, err
	}
	id := token.Claims.(*CustomClaims).Subject
	role := token.Claims.(*CustomClaims).Role

	return map[string]string{"id": id, "role": role}, nil
}
