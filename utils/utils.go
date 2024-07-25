package Utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
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

var jwtSecretKey = []byte("your_secret_key")

func ExtractGuestIDFromToken(tokenString string) (int, error) {

	type CustomClaims struct {
		jwt.StandardClaims
		GuestID string `json:"guestid"`
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		guestID, err := strconv.Atoi(claims.GuestID)
		if err != nil {
			return 0, errors.New("guest_id not found in token")
		}
		fmt.Println("this is the guest id ", guestID)
		return int(guestID), nil
	} else {
		return 0, err
	}
}
