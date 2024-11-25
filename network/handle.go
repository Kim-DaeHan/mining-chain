package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/Kim-DaeHan/mining-chain/blockchain"
)

func HandleBlockList(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload BlockList

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Received blocklist with %d blocks\n", len(payload.Blocks))

	isAppendBlockList = true
	isSync = true
	syncChan <- true

	// 낮은 높이부터 순서대로 blocksInTransit에 블록을 추가
	for _, blockData := range payload.Blocks {
		block := blockchain.Deserialize(blockData)

		// block.Hash가 blocksInTransit에 이미 존재하는지 확인
		isDuplicate := false
		for _, existingBlock := range blocksInTransit {
			if bytes.Equal(existingBlock.Hash, block.Hash) {
				isDuplicate = true
				break
			}
		}

		// 중복되지 않으면 blocksInTransit에 추가
		if !isDuplicate {
			blocksInTransit = append(blocksInTransit, block)
			tempBlockList = append(tempBlockList, block)
			fmt.Printf("Added block %x at height %d to blocksInTransit\n", block.Hash, block.Height)
		} else {
			fmt.Printf("Block %x at height %d is a duplicate and was not added\n", block.Hash, block.Height)
		}
	}

	// 모든 블록이 체인에 추가되었다는 메시지 출력
	fmt.Println("blocksInTransit 에 리스트가 추가됨 길이는:: ", len(blocksInTransit))

	if len(tempBlockList) == payload.Length {
		tempBlockList = nil
		isAppendBlockList = false
		newBlockListChan <- true
	}

	SyncKnownNodes(payload.AddrFrom)

}

// addr 요청을 처리하는 함수
func HandleKnownNodes(request []byte) {
	var buff bytes.Buffer
	var payload Addr

	// 요청 데이터에서 명령어 부분을 제외한 데이터를 버퍼에 저장
	buff.Write(request[commandLength:])
	// 버퍼를 GOB 디코더로 디코딩
	dec := gob.NewDecoder(&buff)
	// Addr 구조체로 디코딩
	err := dec.Decode(&payload)
	if err != nil {
		// 에러 발생 시 패닉
		log.Panic(err)

	}

	// KnownNodes에 새로운 노드 주소 추가
	KnownNodes = append(KnownNodes, payload.AddrList...)
	KnownNodes = RemoveDuplicatesNodes(KnownNodes)
	// 알려진 노드 개수 출력
	fmt.Printf("there are %d known nodes\n", len(KnownNodes))

}

// block 요청을 처리하는 함수
func HandleBlock(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Block
	var blockHeight int64

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blockchain.Deserialize(blockData)

	if chain.CurrentBlock != nil {
		blockHeight = chain.CurrentBlock.Height + 1
	} else {
		blockHeight = 0
	}

	if block.Height > blockHeight {
		otherHeight := block.Height
		SyncWithLongestChain(chain, otherHeight, payload.AddrFrom)
	} else {
		blocksInTransit = append(blocksInTransit, block)

		fmt.Println("blocksInTransit 에 추가됨 길이는:: ", len(blocksInTransit))

		if !isSync {
			isSync = true
			syncChan <- true
		}

		if !isAppendBlockList {
			newBlockListChan <- true
		}
	}

	SyncKnownNodes(payload.AddrFrom)
}

// 특정 데이터 요청을 처리하는 함수
func HandleLatestBlockHeight(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload LatestBlockHeight

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	rangeStr := string(payload.ID)
	var startHeight, endHeight int64
	_, err = fmt.Sscanf(rangeStr, "%d-%d", &startHeight, &endHeight)
	if err != nil {
		fmt.Printf("Invalid block range requested: %s\n", payload.ID)
		return
	}

	blocks := chain.GetBlocksInRange(startHeight, endHeight)

	// 오래된 블록부터 처리하기 위해 슬라이스를 뒤집음
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}

	fmt.Printf("Sending blocks from height %d to %d to %s\n", startHeight, endHeight, payload.AddrFrom)

	// 100개씩 나누어 SendBlockList 호출
	batchSize := 100
	for i := 0; i < len(blocks); i += batchSize {
		end := i + batchSize
		if end > len(blocks) {
			end = len(blocks) // 마지막 batch 처리
		}

		// 슬라이스의 부분 집합을 전달
		batch := blocks[i:end]
		fmt.Printf("Sending batch of blocks from index %d to %d to %s\n", i, end-1, payload.AddrFrom)
		SendBlockList(payload.AddrFrom, batch, len(blocks))

		time.Sleep(1 * time.Second)
	}

}

// 버전 정보를 처리하는 함수
func HandleVersion(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	bestHeight := chain.GetBestHeight()
	otherHeight := payload.BestHeight

	fmt.Printf("Received version from %s with bestHeight: %d. Local bestHeight: %d\n", payload.AddrFrom, otherHeight, bestHeight)

	if bestHeight < otherHeight && payload.AddrFrom != "" {
		fmt.Printf("Node height is lower. Starting sync from height 1 with %s\n", payload.AddrFrom)

		if len(chain.LastHash) > 0 {
			SendLatestBlockHeight(payload.AddrFrom, bestHeight+1, otherHeight)
		} else {
			SendLatestBlockHeight(payload.AddrFrom, 0, otherHeight)
		}

		isSync = true
		syncChan <- true
	} else if bestHeight > otherHeight && payload.AddrFrom != "" {
		fmt.Printf("Node height is higher. Sending version to %s\n", payload.AddrFrom)
		SendVersion(payload.AddrFrom, chain)
	}

	SyncKnownNodes(payload.AddrFrom)
}

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	req, err := io.ReadAll(conn)
	defer conn.Close()

	if err != nil {
		log.Panic(err)
	}

	command := BytesToCmd(req[:commandLength])
	fmt.Printf("Received command: %s\n", command)

	switch command {
	case "knownNodes":
		HandleKnownNodes(req)
	case "block":
		HandleBlock(req, chain)
	case "latestBlockHeight":
		HandleLatestBlockHeight(req, chain)
	case "version":
		HandleVersion(req, chain)
	case "blocklist":
		HandleBlockList(req, chain)
	default:
		fmt.Println("Unknown command")
	}
}
