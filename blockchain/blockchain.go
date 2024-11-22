package blockchain

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime"
	"sort"
	"sync"

	"github.com/Kim-DaeHan/mining-chain/config"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	dbPath      = "./tmp/blocks_%s"
	genesisData = "First Proof from Genesis"
)

type BlockChain struct {
	ChainId      string
	LastHash     []byte
	Database     *leveldb.DB
	CurrentBlock *Block
	Mu           sync.Mutex
}

func (chain *BlockChain) GetBlocksInRange(startHeight, endHeight int64) [][]byte {
	var blocks [][]byte
	iter := chain.NewIterator()

	for {
		block := iter.NextBlock()
		if block == nil || block.Height < startHeight {
			break
		}
		if block.Height <= endHeight {
			blocks = append(blocks, block.Serialize())
		}
	}
	return blocks
}

func (chain *BlockChain) GetExpectedBestHeight() int64 {
	// Retrieve the best (latest) block height from the chain
	lastBlock := chain.GetLastBlock()
	return lastBlock.Height
}

func (chain *BlockChain) AddBlock(block *Block) {
	db := chain.Database

	_, err := db.Get(block.Hash, nil)

	if err == nil {
		return
	}

	lastHash, err := db.Get([]byte("lh"), nil)
	Handle(err)

	item, err := db.Get(lastHash, nil)
	Handle(err)
	lastBlock := Deserialize(item)

	// 블록 높이 검증
	if block.Height <= lastBlock.Height {
		fmt.Println("블록 높이 검증 실패")
		return
	}

	// 추가할려는 블록의 이전해시랑 현재 체인의 마지막 블록 해시랑 같은지 검증
	if string(block.PrevHash) != string(lastBlock.Hash) {
		fmt.Println("블록 해시 검증 실패")
		return
	}

	batch := new(leveldb.Batch)
	blockData := block.Serialize()

	fmt.Println("Serialize block is added", string(blockData))

	batch.Put(block.Hash, blockData)
	batch.Put([]byte("lh"), block.Hash)

	// 블록 높이를 키로 블록 해시 저장
	heightKey := []byte(fmt.Sprintf("height-%d", block.Height))
	batch.Put(heightKey, block.Hash)

	err = db.Write(batch, nil)
	Handle(err)

	chain.CurrentBlock = block
	chain.LastHash = block.Hash
}

func (chain *BlockChain) GetBlockByHeight(height int64) (*Block, error) {
	db := chain.Database

	// 블록 높이로 해시를 가져오기
	heightKey := []byte(fmt.Sprintf("height-%d", height))
	blockHash, err := db.Get(heightKey, nil)
	if err != nil {
		return nil, fmt.Errorf("block with height %d not found: %v", height, err)
	}

	// 해시로 블록 데이터를 가져오기
	blockData, err := db.Get(blockHash, nil)
	if err != nil {
		return nil, fmt.Errorf("block data for hash %x not found: %v", blockHash, err)
	}

	block := Deserialize(blockData)
	return block, nil
}

func (chain *BlockChain) GetLastBlockHash() []byte {
	db := chain.Database
	lasthash, err := db.Get([]byte("lh"), nil)
	if err != nil {
		log.Panic(err)
	}

	return lasthash
}

func (chain *BlockChain) GetLastBlock() *Block {
	db := chain.Database
	var lastBlock *Block
	lasthash, err := db.Get([]byte("lh"), nil)
	if err != nil {
		return DefaultBlock()
	}

	item, err := db.Get(lasthash, nil)
	if err != nil {
		return DefaultBlock()
	}

	lastBlock = Deserialize(item)
	return lastBlock
}

func (chain *BlockChain) GetBestHeight() int64 {
	return chain.GetLastBlock().Height
}

func (chain *BlockChain) GetBlock(blockhash []byte) (Block, error) {
	db := chain.Database
	var block Block

	item, err := db.Get(blockhash, nil)
	block = *Deserialize(item)

	return block, err
}

func (chain *BlockChain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	height := int64(0)

	for {
		block, err := chain.GetBlockByHeight(height)
		if err != nil {
			break // 더 이상 블록이 없으면 종료
		}
		blockCopy := make([]byte, len(block.Hash))
		copy(blockCopy, block.Hash)
		blocks = append(blocks, blockCopy)
		height++
	}

	return blocks
}

func (chain *BlockChain) GetBlockList() []Block {
	var blocks []Block
	iter := chain.NewIterator()

	for {
		block := iter.NextBlock()
		blocks = append(blocks, *block)

		// 마지막 블록(이전 해시가 없는 블록)에 도달하면 반복 종료
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return blocks
}

func (chain *BlockChain) Difficulty(height int64) *big.Int {
	Config := config.GlobalConfig

	if height < (Config.DIFFICULTY_CHANGE_CYCLE + 1) {
		return new(big.Int).Set(&Config.DEFAULT_DIFFICULTY)
	} else if height%Config.DIFFICULTY_CHANGE_CYCLE != 1 {
		block, err := chain.GetBlockByHeight(height - 1)
		Handle(err)
		return block.Difficulty
	}

	endBlock, err := chain.GetBlockByHeight(height - 1)
	Handle(err)
	startBlock, err := chain.GetBlockByHeight(height - Config.DIFFICULTY_CHANGE_CYCLE - 1)
	Handle(err)
	lastDifficulty := endBlock.Difficulty

	gap := endBlock.Timestamp - startBlock.Timestamp
	standardGap := Config.RESOURCE_INTERVAL * Config.DIFFICULTY_CHANGE_CYCLE

	weight := float64(standardGap) / float64(gap)
	if weight > Config.MAX_DIFFICULTY_WEIGHT {
		weight = Config.MAX_DIFFICULTY_WEIGHT
	}
	if weight < Config.MIN_DIFFICULTY_WEIGHT {
		weight = Config.MIN_DIFFICULTY_WEIGHT
	}

	// 난이도 계산
	weightBig := big.NewFloat(weight)
	difficultyFloat := new(big.Float).SetInt(lastDifficulty)
	difficultyFloat.Mul(difficultyFloat, weightBig)

	resultDifficulty, _ := difficultyFloat.Int(nil)

	if resultDifficulty.Cmp(big.NewInt(1)) <= 0 {
		return big.NewInt(1)
	}

	return resultDifficulty
}

func InitBlockChain(address, chainId string) *BlockChain {
	fmt.Printf("init blockchain path : %s\n", chainId)

	path := fmt.Sprintf(dbPath, chainId)
	fmt.Printf("init blockchain path : %s\n", path)
	if DBexists(path) {
		// err := os.RemoveAll(path)
		// if err != nil {
		// 	log.Fatalf("Failed to delete existing blockchain data: %v", err)
		// }
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := leveldb.OpenFile(path, nil)
	Handle(err)
	batch := new(leveldb.Batch)

	genesis := Genesis(address)

	fmt.Printf("genesis hash :%v\n", string(genesis.Serialize()))
	batch.Put(genesis.Hash, genesis.Serialize())
	batch.Put([]byte("lh"), genesis.Hash)
	heightKey := []byte(fmt.Sprintf("height-%d", genesis.Height))
	batch.Put(heightKey, genesis.Hash)

	err = db.Write(batch, nil)
	Handle(err)
	lastHash = genesis.Hash

	chain := BlockChain{
		LastHash: lastHash,
		Database: db,
	}
	return &chain
}

func DBexists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockChain(chainId string) *BlockChain {
	path := fmt.Sprintf(dbPath, chainId)
	fmt.Printf("blockchain Path : %s\n", path)
	if !DBexists(path) {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := leveldb.OpenFile(path, nil)
	Handle(err)

	lastHash, _ = db.Get([]byte("lh"), nil)
	Handle(err)

	chain := BlockChain{
		ChainId:  chainId,
		LastHash: lastHash,
		Database: db,
	}

	return &chain
}

func (chain *BlockChain) ResetDatabase() {
	db := chain.Database
	iter := db.NewIterator(nil, nil)

	for iter.Next() {
		key := iter.Key()
		db.Delete(key, nil)
	}

	iter.Release()
	fmt.Println("로컬 데이터베이스 초기화 완료")
}

func SortBlocksByHeight(blocks []*Block) {
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Height < blocks[j].Height // Height 기준 오름차순 정렬
	})
}
