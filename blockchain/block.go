package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/Kim-DaeHan/mining-chain/config"
)

type HexBytes []byte

// MarshalJSON implements the json.Marshaler interface for HexBytes.
func (h HexBytes) MarshalJSON() ([]byte, error) {
	if h == nil || len(h) == 0 {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%x"`, h)), nil
}

func (h *HexBytes) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}

	// 유효성 검사: 16진수 형식만 허용
	if len(hexStr)%2 != 0 || !isHexString(hexStr) {
		return fmt.Errorf("invalid hex string: %s", hexStr)
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	*h = bytes
	return nil
}

// 헬퍼 함수: 문자열이 유효한 16진수인지 검사
func isHexString(s string) bool {
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}
	return true
}

type Block struct {
	Timestamp       int64
	Hash            HexBytes
	PrevHash        HexBytes
	MainBlockHeight int
	MainBlockHash   HexBytes
	Nonce           HexBytes
	Height          int64
	Difficulty      *big.Int
	Miner           HexBytes
	Validator       HexBytes
}

func CreateBlock(prevHash []byte, height int64, address string) *Block {
	block := &Block{
		Timestamp:       int64(0),
		Hash:            HexBytes{},
		PrevHash:        HexBytes(prevHash),
		MainBlockHeight: 0,
		MainBlockHash:   HexBytes{},
		Nonce:           HexBytes{},
		Height:          height,
		Difficulty:      big.NewInt(config.GlobalConfig.DEFAULT_DIFFICULTY.Int64()),
		Miner:           HexBytes(address),
		Validator:       HexBytes(address),
	}

	pow := NewProof(block)
	nonce := pow.Run()
	block.Hash = HexBytes(pow.GetHash(block))
	block.Nonce = HexBytes(nonce)

	return block
}

func Genesis(address string) *Block {
	return CreateBlock([]byte{}, 0, address)
}

func (b *Block) Serialize() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		log.Panic(err)
	}
	return data
}

func Deserialize(data []byte) *Block {
	var block Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func DefaultBlock() *Block {
	return &Block{
		Timestamp:       0,
		Hash:            nil,
		PrevHash:        nil,
		MainBlockHeight: 0,
		MainBlockHash:   nil,
		Nonce:           nil,
		Height:          0,
		Difficulty:      big.NewInt(1),
		Miner:           nil,
		Validator:       nil,
	}
}

func (b Block) String() string {
	var lines []string

	// 블록의 기본 정보를 추가
	lines = append(lines, "----- Block -----")
	lines = append(lines, fmt.Sprintf("Height:      %d", b.Height))
	lines = append(lines, fmt.Sprintf("Timestamp:   %d", b.Timestamp))
	lines = append(lines, fmt.Sprintf("Hash:        %x", b.Hash))
	lines = append(lines, fmt.Sprintf("PrevHash:    %x", b.PrevHash))
	lines = append(lines, fmt.Sprintf("Nonce:       %x", b.Nonce))
	lines = append(lines, fmt.Sprintf("Difficulty:  %d", b.Difficulty))
	lines = append(lines, fmt.Sprintf("Miner:  %x", b.Miner))
	lines = append(lines, fmt.Sprintf("Validator:  %x", b.Validator))
	lines = append(lines, fmt.Sprintln())

	// 모든 정보를 개행 문자로 구분하여 하나의 문자열로 결합
	return strings.Join(lines, "\n")
}
