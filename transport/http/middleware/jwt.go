package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/jwtmodel"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)



type JWTAuthentication struct {
}


const (
	HeaderJWTAuthorization = "Authorization"
)

func ProvideJWTAuthentication(db *infras.MySQLConn) *JWTAuthentication {
	return &JWTAuthentication{}
}


func validateJWT(tokenString string) (*jwtmodel.Claims, error)  {
	secret := viper.GetString("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &jwtmodel.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("JWT validation failed: %v", err)
	}

	if claims, ok := token.Claims.(*jwtmodel.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("JWT is not valid")
}


func (a *JWTAuthentication) JWTMiddlewareValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			log.Println("no header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		claims, err := validateJWT(tokenString)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}


		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}