package mining

import (
	"context"
	"fmt"
	"time"

	"github.com/Kim-DaeHan/mining-chain/blockchain"
)

type Mining struct {
	chainId string
}

func Run(ctx context.Context,
	chain *blockchain.BlockChain,
	validator string,
	miningBlockChan chan *blockchain.Block) {

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Mining stopped.")
			return
		default:
			lastBlock := chain.GetLastBlock()

			chain.Mu.Lock()
			block := &blockchain.Block{
				Timestamp:  time.Now().Unix(),
				PrevHash:   lastBlock.Hash,
				Height:     lastBlock.Height + 1,
				Difficulty: chain.Difficulty(lastBlock.Height + 1),
				Miner:      blockchain.HexBytes(validator),
				Validator:  blockchain.HexBytes(validator),
			}
			chain.Mu.Unlock()

			// 경합으로 인한 분기 최소화
			time.Sleep(1 * time.Second)

			// Mining work
			pow := blockchain.NewProof(block)
			nonceByte := pow.Run()

			select {
			case <-ctx.Done():
				// 채굴 도중 중단 요청 확인
				fmt.Println("Mining interrupted before block completion.")
				return
			default:
				block.Nonce = nonceByte
				block.Hash = pow.GetHash(block)
				fmt.Println("miningBlock: ", block)

				miningBlockChan <- block
			}
		}
	}
}
