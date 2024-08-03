package Utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	// "strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// var jwtSecretKey = []byte("your_secret_key")

func ExtractGuestIDFromToken(tokenString string) (string, error) {
	var jwtSecretKey = []byte("your_secret_key")

	type CustomClaims struct {
		jwt.StandardClaims
		GuestID string `json:"guestid,omitempty"`
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		guestID := claims.GuestID
		if guestID == "" {
			return "", errors.New("guest_id not found in token")
		}
		fmt.Println("this is the guest id ", guestID)
		return guestID, nil
	} else {
		return "", err
	}
}

func ExtractGuestIDFromTokenn(tokenString string) (string, *jwt.Token, error) {
	var jwtSecretKey = []byte("your_secret_key")

	type CustomClaims struct {
		jwt.StandardClaims
		GuestID string `json:"guestid,omitempty"`
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", nil, fmt.Errorf("error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		guestID := claims.GuestID
		if guestID == "" {
			return "", nil, errors.New("guest_id not found in token")
		}
		fmt.Println("this is the guest id ", guestID)
		return guestID, token, nil
	} else {
		return "", nil, errors.New("invalid token or claims")
	}
}
