package service

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(name string, admin bool) string
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// jwtCustomClaims are custom claims extending default ones
type jwtCustomClaims struct {
	Name string `json:"name"`
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer string
}

func NewJWTService() JWTService {
	return &jwtService {
		secretKey: getSecretKey(),
		issuer: "github.com",
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (jwtSrv *jwtService) GenerateToken(username string, admin bool) string {

	// set custom and standard claims
	claims := &jwtCustomClaims {
		username,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer: jwtSrv.issuer,
			IssuedAt: time.Now().Unix(),
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate encoded token using the secret signing key
	t, err := token.SignedString([]byte(jwtSrv.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (jwtSrv *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return the secret signing key
		return []byte(jwtSrv.secretKey), nil
	})
}
