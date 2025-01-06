package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswords(t *testing.T) {
	t.Run("hash password", func(t *testing.T) {
		pwHash, err := HashPassword("test")

		assert.NoError(t, err)
		assert.NotEmpty(t, pwHash)
	})

	t.Run("long password", func(t *testing.T) {
		pwHash, err := HashPassword(strings.Repeat("A", 100))

		assert.Error(t, err)
		assert.Empty(t, pwHash)
	})

	t.Run("check correct password", func(t *testing.T) {
		pw := "test"
		pwHash, err := HashPassword(pw)

		assert.NoError(t, err)
		assert.NotEmpty(t, pwHash)

		checkResult := CheckPassword(pwHash, pw)

		assert.True(t, checkResult)
	})
	t.Run("check incorrect password", func(t *testing.T) {
		pw := "test"
		pwHash, err := HashPassword(pw)

		assert.NoError(t, err)
		assert.NotEmpty(t, pwHash)

		checkResult := CheckPassword(pwHash, "invalid")

		assert.False(t, checkResult)
	})
}

func TestTokens(t *testing.T) {
	secret := []byte("test_secret_key")
	testUserID := uint64(1337)
	expireDate := time.Now().Add(time.Hour)

	t.Run("create token", func(t *testing.T) {
		token, err := CreateToken(int(testUserID), expireDate, secret)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("successful token verification", func(t *testing.T) {
		tokenString, err := CreateToken(int(testUserID), expireDate, secret)

		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		claims, err := VerifyToken(tokenString, secret)

		assert.NoError(t, err)

		assert.Contains(t, claims, "user_id", "No user_id in claims")

		claimedUserID := uint64(claims["user_id"].(float64))

		assert.Equal(t, testUserID, claimedUserID)
	})

	t.Run("verify invalid token", func(t *testing.T) {
		invalidToken := "invalid"

		token, err := VerifyToken(invalidToken, secret)

		assert.Nil(t, token)
		assert.Error(t, err)
	})

	t.Run("verify token with incorrect alg", func(t *testing.T) {
		tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiO" +
			"iI4NWEwMzg2Ny1kY2NmLTQ4ODItYWRkZS0xYTc5YWVlYzUwZGYiLCJleHAiOjE2NDQ4ODQ" +
			"xODUsImlhdCI6MTY0NDg4MDU4NSwiaXNzIjoiYWNtZS5jb20iLCJzdWIiOiIwMDAwMDAwM" +
			"C0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDEiLCJqdGkiOiIzZGQ2NDM0ZC03OWE5LTR" +
			"kMTUtOThiNS03YjUxZGJiMmNkMzEiLCJhdXRoZW50aWNhdGlvblR5cGUiOiJQQVNTV09SR" +
			"CIsImVtYWlsIjoiYWRtaW5AZnVzaW9uYXV0aC5pbyIsImVtYWlsX3ZlcmlmaWVkIjp0cnV" +
			"lLCJhcHBsaWNhdGlvbklkIjoiODVhMDM4NjctZGNjZi00ODgyLWFkZGUtMWE3OWFlZWM1M" +
			"GRmIiwicm9sZXMiOlsiY2VvIl19.dee-Ke6RzR0G9avaLNRZf1GUCDfe8Zbk9L2c7yaqKME"

		token, err := ParseToken(tokenString, secret)

		assert.Nil(t, token)
		assert.Error(t, err)
	})
}
