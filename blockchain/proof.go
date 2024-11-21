package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"sync"

	"github.com/Kim-DaeHan/mining-chain/utils"
)

const Difficulty = 30

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) InitData(nonce []byte) []byte {

	data := bytes.Join(
		[][]byte{
			ToHex(pow.Block.Timestamp),
			nonce},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() []byte {
	numThreads := runtime.NumCPU()
	results := make(chan []byte, numThreads)
	done := make(chan struct{})
	var once sync.Once

	hashLimit, err := HashLimit(pow.Block)
	Handle(err)

	for i := 0; i < numThreads; i++ {
		go func(threadID int) {
			for {
				select {
				case <-done:
					return // 다른 고루틴에서 이미 작업이 완료되었으므로 종료
				default:
					nonce := utils.GenerateRandomHex64bit()
					blockRoot := pow.BlockRoot(pow.Block, nonce)

					nonceBytes, err := hex.DecodeString(nonce)
					Handle(err)

					if hashLimit >= blockRoot {
						results <- nonceBytes
						// 한 번만 done 채널 닫기
						once.Do(func() { close(done) })
						return
					}
				}
			}
		}(i)
	}

	select {
	case result := <-results:
		return result
	case <-done:
		// results 채널을 닫지 않음, 필요시 go 루틴에서 닫도록 처리
	}

	return nil
}

func (pow *ProofOfWork) GetHash(block *Block) []byte {
	data := pow.InitData(block.Nonce)
	hash := sha256.Sum256(data)
	return hash[:]
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// hashLimit 함수는 난이도 값을 받아서 큰 수 2^256을 난이도 값으로 나눈 결과를 64자리 16진수 문자열로 반환합니다.
func HashLimit(b *Block) (string, error) {
	diff := b.Difficulty

	a := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

	// 난이도 값이 0이거나 음수인 경우 오류 반환
	if diff.Cmp(big.NewInt(1)) < 0 {
		return "", fmt.Errorf("invalid diff value")
	}

	result := new(big.Int).Div(a, diff)
	hexResult := result.Text(16)
	paddedHexResult := fmt.Sprintf("%064s", hexResult)

	return paddedHexResult, nil
}

// ComputeSHA256은 주어진 문자열을 SHA-256으로 변환하여 반환합니다.
func computeSHA256(input string) string {
	hash := sha256.Sum256([]byte(input)) // 한번에 해시 계산
	return hex.EncodeToString(hash[:])   // 해시를 헥스 문자열로 변환
}

func (pow *ProofOfWork) BlockRoot(block *Block, nonce string) string {

	blockInfo, err := utils.ToJSONString(block)
	if err != nil {
		fmt.Println("Error converting block to JSON string:", err)
		Handle(err)
	}
	// Handle(err)
	combinedString := string(block.PrevHash) + blockInfo + nonce
	sha256Hash := utils.ComputeSHA256(combinedString)
	return sha256Hash
}
