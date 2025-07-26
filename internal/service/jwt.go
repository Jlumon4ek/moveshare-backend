package service

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth interface {
	GenerateAccessToken(userID int64, username, email string) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	ValidateToken(tokenString string) (int64, error)
}

type jwtAuth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTAuth(privateKeyPath, publicKeyPath string) (JWTAuth, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	log.Printf("Keys loaded successfully: private key length %d, public key length %d", len(privateKeyData), len(publicKeyData))
	return &jwtAuth{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (j *jwtAuth) GenerateAccessToken(userID int64, username, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j *jwtAuth) GenerateRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j *jwtAuth) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		log.Printf("Token validation error: %v", err)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(float64)
		if !ok {
			log.Printf("Invalid user ID in token claims: %v", claims)
			return 0, fmt.Errorf("invalid user ID in token")
		}
		log.Printf("Token validated successfully for user ID: %d", int64(userID))
		return int64(userID), nil
	}

	log.Printf("Token invalid or claims unreadable: %v", token.Claims)
	return 0, fmt.Errorf("invalid token")
}
