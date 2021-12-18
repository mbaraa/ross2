package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mbaraa/ross2/config"
)

type GoogleJWTTokenValidator struct{}

func NewGoogleJWTTokenValidator() *GoogleJWTTokenValidator {
	return new(GoogleJWTTokenValidator)
}

func (g *GoogleJWTTokenValidator) Validate(token string) error {
	_, err := g.validateGoogleJWT(token)
	return err
}

func (g *GoogleJWTTokenValidator) validateGoogleJWT(tokenString string) (googleClaims, error) {
	gclaims := googleClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &gclaims, func(token *jwt.Token) (interface{}, error) {
		pem, err := g.getGooglePublicKey(token.Header["kid"].(string))
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return googleClaims{}, err
	}

	claims, ok := token.Claims.(*googleClaims)
	if !ok {
		return googleClaims{}, err
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return googleClaims{}, err
	}

	if claims.Audience != config.GetInstance().GoogleClientID {
		return googleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return googleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func (g *GoogleJWTTokenValidator) getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

type googleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}
