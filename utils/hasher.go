package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeSHA256은 주어진 문자열을 SHA-256으로 변환하여 반환합니다.
func ComputeSHA256(input string) string {
	hash := sha256.Sum256([]byte(input)) // 한번에 해시 계산
	return hex.EncodeToString(hash[:])   // 해시를 헥스 문자열로 변환
}
