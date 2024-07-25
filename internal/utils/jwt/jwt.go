package jwt

import (
	"campaign/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_SECRET   = os.Getenv("JWT_SECRET")
	AUTH_CTX_KEY = &contextKey{"authcontext"}
)

type AuthContext struct {
	// Sub is the subject of the token
	// This is used to determine the user the token belongs to
	Sub string

	// Issuer is the issuer of the token
	// This is used to determine the source of the token
	Issuer string

	// Name is the name of the user
	Name string
}

type contextKey struct {
	name string
}

var (
	ErrUnauthorized = errors.New("token is unauthorized")
	ErrExpired      = errors.New("token is expired")
	ErrNBFInvalid   = errors.New("token nbf validation failed")
	ErrIATInvalid   = errors.New("token iat validation failed")
	ErrNoTokenFound = errors.New("no token found")
	ErrAlgoInvalid  = errors.New("algorithm mismatch")
)

// Authenticator is a middleware that handles jwt authentications
func Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			var token string

			findtokens := []func(*http.Request) string{GetTokenFromHeader}

			for _, fn := range findtokens {
				token = fn(r)

				if token != "" {
					break
				}

			}

			if token == "" {

				res := utils.ApiResponse{
					Message: "Unauthorized",
					Data:    nil,
				}

				w.WriteHeader(http.StatusUnauthorized)
				response, _ := json.Marshal(res)

				_, _ = w.Write(response)

				return

			}

			c, err := ParseToken(token)
			if err != nil {

				res := utils.ApiResponse{
					Message: "Unauthorized",
					Data:    nil,
				}

				w.WriteHeader(http.StatusUnauthorized)
				response, _ := json.Marshal(res)

				_, _ = w.Write(response)

				return
			}

			ctx, err := newContext(r.Context(), c)
			if err != nil {

				res := utils.ApiResponse{
					Message: "Unauthorized",
					Data:    nil,
				}

				w.WriteHeader(http.StatusUnauthorized)
				response, _ := json.Marshal(res)

				_, _ = w.Write(response)
				return

			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(hfn)
	}
}

func GetTokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")

	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		slog.Info("got token from request head")
		return bearer[7:]
	}

	return ""
}

/*
Generetes a signed token and return as byte or nil.
*/
func GenereteJWT(data AuthContext) ([]byte, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    data.Sub,
		"name":   data.Name,
		"issuer": "campaign",
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		slog.Error("Error signing jwt", "Error", err)

		return nil, fmt.Errorf("Something went wrong. Please try again. #2")
	}

	return []byte(tokenString), nil
}

// ParseToken is use to verify jwt tokens
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, alg := token.Method.(*jwt.SigningMethodHMAC); !alg {

			slog.Warn("invalid token alg", "ParseToken", token.Header["alg"])

			return nil, ErrAlgoInvalid
		}

		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		slog.Error("error parsing token [%s]", "Error", err)

		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, ErrNBFInvalid
	}

	return claims, nil
}

func newContext(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	sub := claims["sub"]

	if sub == "" {
		return nil, ErrUnauthorized
	}

	ctx = context.WithValue(ctx, AUTH_CTX_KEY, claims)
	return ctx, nil
}

// GetAuthContext returns decoded jwt data
func GetAuthContext(ctx context.Context) (AuthContext, error) {
	claims, ok := ctx.Value(AUTH_CTX_KEY).(jwt.MapClaims)

	if !ok {
		slog.Error("GetAuthContext - no claims found", "Error", ErrNoTokenFound)

		return AuthContext{}, ErrNoTokenFound
	}

	var sub string
	var issuer string
	var name string

	if claims["sub"] != nil {
		sub = claims["sub"].(string)
	}

	if claims["issuer"] != nil {
		issuer = claims["issuer"].(string)
	}

	if claims["name"] != nil {
		name = claims["name"].(string)
	}

	if sub == "" {
		slog.Error("GetAuthContext - missing claims", "Error", ErrNoTokenFound)

		return AuthContext{}, ErrNoTokenFound
	}

	return AuthContext{
		Sub:    sub,
		Issuer: issuer,
		Name:   name,
	}, nil
}
