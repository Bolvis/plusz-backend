package authorization

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Token struct {
	UserId    string    `json:"userId"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func GenerateToken(userId string) (string, error) {
	token := Token{
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	payload, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte("secret"))
	h.Write(payload)
	signature := h.Sum(nil)

	tokenString := strings.Join([]string{
		base64.URLEncoding.EncodeToString(payload),
		base64.StdEncoding.EncodeToString(signature)},
		".",
	)

	return tokenString, nil
}

func VerifyToken(tokenString string) (Token, error) {
	token := Token{}
	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 {
		return token, fmt.Errorf("invalid token")
	}

	payload, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return token, fmt.Errorf("invalid payload encoding")
	}

	h := hmac.New(sha256.New, []byte("secret"))
	h.Write(payload)
	expectedSig := h.Sum(nil)

	actualSig, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return token, fmt.Errorf("invalid signature encoding: %w", err)
	}

	if !hmac.Equal(expectedSig, actualSig) {
		return token, fmt.Errorf("invalid signature")
	}

	if err := json.Unmarshal(payload, &token); err != nil {
		return token, fmt.Errorf("invalid token data: %w", err)
	}

	if time.Now().After(token.ExpiresAt) {
		return token, fmt.Errorf("token expired")
	}

	return token, nil
}
