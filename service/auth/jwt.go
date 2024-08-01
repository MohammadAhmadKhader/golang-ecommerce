package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohammadahmadkhader/golang-ecommerce/config"
	"github.com/mohammadahmadkhader/golang-ecommerce/types"
	"github.com/mohammadahmadkhader/golang-ecommerce/utils"
)

type tokenKey string
const tokenPayloadKey = tokenKey("tokenPayload")

type tokenPayload struct {
	Email string `json:"email"`
	UserId int `json:"userId"`
}

func CreateJWT(secret []byte, user types.User) (string, error) {
	durationInt, err := strconv.Atoi(config.Envs.JWTExpirationInSeconds)
	if err != nil {
		return "", err
	}

	expiration := time.Second * time.Duration(durationInt)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(user.ID),
		"email":     user.Email,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func deCryptToken(r *http.Request) (*jwt.MapClaims, error) {
	jwtToken := r.Header.Get("Authorization")
	jwtSecret := config.Envs.JWTSecret
	if jwtToken == "" {
		return nil, fmt.Errorf("missing token")
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unauthenticated")
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	isValidToken := token.Valid
	if !isValidToken {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		return nil, fmt.Errorf("invalid token")
	} else {

		return &claims, nil
	}
}

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := deCryptToken(r)
		if err != nil {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("forbidden"))
			return
		}

		
		ctx := context.WithValue(r.Context(), tokenPayloadKey, claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func GetTokenPayload(ctx context.Context) (tokenPayload, error) {
	payload, ok := ctx.Value(tokenPayloadKey).(tokenPayload)
	if !ok {
		return tokenPayload{}, fmt.Errorf("invalid token")
	}

	return payload, nil
}