package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
)

type Config struct {
	ChainId                 int     `json:"chainId"`
	Port                    int     `json:"port"`
	RPCPort                 int     `json:"rpcPort"`
	NodeType                string  `json:"nodeType"`
	Mining                  bool    `json:"mining"`
	DEFAULT_DIFFICULTY      big.Int // 하드코딩된 값
	DIFFICULTY_CHANGE_CYCLE int64   // 하드코딩된 값
	RESOURCE_INTERVAL       int64   // 하드코딩된 값
	MAX_DIFFICULTY_WEIGHT   float64 // 하드코딩된 값
	MIN_DIFFICULTY_WEIGHT   float64 // 하드코딩된 값
}

// 전역 설정 파일 경로
const configPath = "./tmp/config.json"

// 기본값을 가진 전역 Config 변수
var GlobalConfig = Config{
	ChainId:                 1,                   // 하드 코딩된 기본값
	Port:                    8080,                // 하드 코딩된 기본값
	RPCPort:                 8545,                // 하드 코딩된 기본값
	NodeType:                "full-node",         // 하드 코딩된 기본값
	Mining:                  false,               // 하드 코딩된 기본값
	DEFAULT_DIFFICULTY:      *big.NewInt(500000), // 하드 코딩된 값
	DIFFICULTY_CHANGE_CYCLE: 20,                  // 하드 코딩된 값
	RESOURCE_INTERVAL:       20,                  // 하드 코딩된 값 (초 단위)
	MAX_DIFFICULTY_WEIGHT:   4.0,                 // 하드 코딩된 값
	MIN_DIFFICULTY_WEIGHT:   0.25,                // 하드 코딩된 값
}

// 설정 파일을 로드하는 함수
func LoadConfig() error {
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&GlobalConfig)
	if err != nil {
		return fmt.Errorf("could not decode config file: %v", err)
	}
	return nil
}
