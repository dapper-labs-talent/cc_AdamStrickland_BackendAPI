package security

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/golang-jwt/jwt"
)

type TokenPayload struct {
	Token string `json:"token"`
}

func parsedToken(cfg *config.Config, token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		return []byte(cfg.GetString("secret")), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, errors.New("Invalid token (token is invalid)")
	}

	return t, nil
}

func tokenClaims(cfg *config.Config, token string) (jwt.MapClaims, error) {
	t, err := parsedToken(cfg, token)

	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("Claims could not be extracted")
	}

	return claims, nil
}

func TokenSubject(cfg *config.Config, token string) (*string, error) {
	claims, err := tokenClaims(cfg, token)

	if err != nil {
		return nil, err
	}

	subj := fmt.Sprintf("%v", claims["sub"])

	return &subj, nil
}

func IsValidToken(cfg *config.Config, token string) (bool, error) {
	claims, err := tokenClaims(cfg, token)

	if err != nil {
		return false, err
	}

	err = claims.Valid()

	if err != nil {
		return false, err
	}

	return true, nil
}

func NewTokenPayload(cfg *config.Config, subj string) ([]byte, error) {
	ts, err := NewTokenForSubject(cfg, subj)

	if err != nil {
		log.Printf("Unable to generate token: %e", err)
		return nil, err
	}

	sp := &TokenPayload{
		Token: ts,
	}

	var data bytes.Buffer
	err = json.NewEncoder(&data).Encode(sp)

	if err != nil {
		log.Printf("Unable to generate JWT payload: %e", err)
		return nil, err
	}

	return data.Bytes(), nil
}

func newTokenWithClaims(cfg *config.Config, claims *jwt.StandardClaims) (string, error) {

	log.Printf("Issuing authorization: '%+v'", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := cfg.GetString("secret")

	ss, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Printf("Unable to generate JWT: %e", err)
	}

	return ss, err
}

func NewTokenForSubject(cfg *config.Config, subj string) (string, error) {
	ts := time.Now()

	claims := &jwt.StandardClaims{
		Audience:  "dapper-client",
		ExpiresAt: ts.Add(time.Hour * 24).Unix(),
		Issuer:    "dapper-api",
		IssuedAt:  ts.Unix(),
		Subject:   subj,
	}

	return newTokenWithClaims(cfg, claims)
}
