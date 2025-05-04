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
		fmt.Println(err)
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
		err := fmt.Errorf("invalid token")
		fmt.Println(err)
		return token, err
	}

	payload, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		fmt.Println(err)
		return token, fmt.Errorf("invalid payload encoding")
	}

	h := hmac.New(sha256.New, []byte("secret"))
	h.Write(payload)
	expectedSig := h.Sum(nil)

	actualSig, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		err := fmt.Errorf("invalid signature encoding: %w", err)
		fmt.Println(err)
		return token, err
	}

	if !hmac.Equal(expectedSig, actualSig) {
		err := fmt.Errorf("invalid signature")
		fmt.Println(err)
		return token, err
	}

	if err := json.Unmarshal(payload, &token); err != nil {
		err := fmt.Errorf("invalid token data: %w", err)
		fmt.Println(err)
		return token, err
	}

	if time.Now().After(token.ExpiresAt) {
		err := fmt.Errorf("token expired")
		fmt.Println(err)
		return token, err
	}

	return token, nil
}
