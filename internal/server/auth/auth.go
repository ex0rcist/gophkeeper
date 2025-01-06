// Provides auth functions, jwt vaildation and checks
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Hash password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Check if provided password matches user's password
func CheckPassword(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

// Creates JWT token
func CreateToken(userID int, expireDate time.Time, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"iss":     "gophkeeper",
		"exp":     expireDate.Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Verifies validity of token and returns claims
func VerifyToken(tokenText string, secretKey []byte) (jwt.MapClaims, error) {
	token, err := ParseToken(tokenText, secretKey)
	if err != nil {
		return nil, err
	}

	claims, err := GetClaims(token)
	if err != nil {
		return nil, err
	}

	if IsExpired(claims) {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

// Parses token and returns pointer to jwt.Token
func ParseToken(tokenText string, secretKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenText, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

// Checks if provided token is expired
func IsExpired(claims jwt.MapClaims) bool {
	return float64(time.Now().Unix()) > claims["exp"].(float64)
}

// Extracts claims from provided token
func GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}

	return claims, nil
}
