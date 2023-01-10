package bookmark

import (
	"context"
	"fmt"
	"log"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"github.com/golang-jwt/jwt/v4"
)

type Data struct {
	Name    string
	Picture string
}

//
//encore:authhandler
func AuthHandler(ctx context.Context, tokenString string) (auth.UID, *Data, error) {
	key := secrets.JwtSecret

	// Parse and decrypt the Json web token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return "", nil, err
	}

	// Validate the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println("claims=", claims)
		return auth.UID(claims["sub"].(string)), &Data{
			Name:    claims["name"].(string),
			Picture: claims["picture"].(string),
		}, nil
	}

	return "", nil, &errs.Error{
		Code:    errs.Unauthenticated,
		Message: "invalid token",
	}
}
