package core

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Bcrypt(text string) (string) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		errors.New("hashing error")
	}
	return string(hashedBytes)
}

func CheckBcrypt(text string, hashedValue string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(text)) == nil
}

func JsonResponse(w http.ResponseWriter, data map[string]interface{}, statusCode int)  {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

type RequestBody struct {
	Key string `json:"key"`
}

func JsonBody(r *http.Request) map[string]string {
	var requestBody map[string]string

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		errors.New("error decoding JSON")
	}

	return requestBody

	// Access the dynamic keys and their values
	//for _key, value := range requestBody {
	//	if _key == key {
	//		return value
	//	}
	//}
	//
	//// Access the value by key
	//errors.New("key not found")
	//return ""
}
