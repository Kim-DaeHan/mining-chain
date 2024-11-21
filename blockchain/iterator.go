package blockchain

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type BlockchainIterator struct {
	currentHash []byte
	Database    *leveldb.DB
}

func (chain *BlockChain) NewIterator() *BlockchainIterator {
	return &BlockchainIterator{currentHash: chain.LastHash, Database: chain.Database}
}

func (iter *BlockchainIterator) NextBlock() *Block {
	blockData, err := iter.Database.Get(iter.currentHash, nil)
	if err != nil {
		fmt.Println("End of chain reached or error occurred:", err)
		return nil
	}

	block := Deserialize(blockData)
	iter.currentHash = block.PrevHash

	return block
}
