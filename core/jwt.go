package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func GenerateJwt(userId string, jti string, expiresIn time.Duration) string {
	// Define the header and payload for the JWT
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	payload := map[string]interface{}{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * expiresIn).Unix(),
		"jti": jti,
	}

	// Encode the header and payload as base64 URL-safe strings
	headerBytes, _ := json.Marshal(header)
	payloadBytes, _ := json.Marshal(payload)

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Create the signature by concatenating and signing the encoded header and payload
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	signatureInput := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)
	signature := hmacSHA256Encode(secretKey, []byte(signatureInput))

	// Combine the encoded header, payload, and signature to create the JWT
	jwt := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, base64.RawURLEncoding.EncodeToString(signature))

	fmt.Println("Generated JWT:", jwt)

	return jwt
}

func hmacSHA256Encode(key, message []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return h.Sum(nil)
}

type JWTClaims struct {
	UserId string `json:"sub"`
	Jti   string `json:"jti"`
	Exp      int64  `json:"exp"`
}

type JWTHeader struct {
	Alg string `json:"alg"`
}

type JWTDecoded struct {
	Success bool
	Jti string
	UserId string
	Message string
}

func DecodeJwt(jwt string) JWTDecoded {
	// Replace this with your actual JWT

	secretKey := os.Getenv("SECRET_KEY")

	segments := strings.Split(jwt, ".")
	var response JWTDecoded
	response.Success = false
	if len(segments) != 3 {
		response.Message = "Invalid JWT format"
		return response
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(segments[0])
	if err != nil {
		response.Message = "Error decoding header:"
		return response
	}

	var header JWTHeader
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		response.Message = "Error unMarshaling header:"
		return response
	}

	if header.Alg != "HS256" {
		response.Message = "Unsupported algorithm:"
		return response
	}

	// Fetch the secret key based on header.Alg

	// Compute expected signature using HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(segments[0] + "." + segments[1]))
	expectedSignature := mac.Sum(nil)

	// Decode the provided signature
	actualSignature, err := base64.RawURLEncoding.DecodeString(segments[2])
	if err != nil {
		response.Message = "Error decoding signature:"
		return response
	}

	// Compare expected and actual signatures
	if !hmac.Equal(expectedSignature, actualSignature) {
		response.Message = "Invalid JWT"
		return response
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(segments[1])
	if err != nil {
		response.Message = "Error decoding payload:"
		return response
	}

	var claims JWTClaims
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		response.Message = "Error unmarshaling claims:"
		return response
	}

	response.Success = true
	response.Jti = claims.Jti
	response.UserId = claims.UserId

	return response
}
