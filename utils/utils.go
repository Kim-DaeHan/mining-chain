package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
)

// ToJSONString은 BlockInfo 구조체를 JSON 문자열로 변환하는 함수입니다.
func ToJSONString(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("input is nil")
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Printf("Failed to marshal JSON for value: %v, error: %v", v, err)
		return "", err
	}
	return string(jsonBytes), nil
}

// generateRandomHex64bit는 64비트(8바이트) 크기의 랜덤 16진수 문자열을 생성합니다.
func GenerateRandomHex64bit() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	// SHA-256 해시를 계산
	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:])
}
