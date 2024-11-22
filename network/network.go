package network

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"syscall"

	"github.com/Kim-DaeHan/mining-chain/blockchain"
	"github.com/Kim-DaeHan/mining-chain/config"
	"github.com/Kim-DaeHan/mining-chain/mining"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/vrecan/death/v3"
)

const (
	protocol      = "tcp"
	version       = 1
	commandLength = 20
)

var (
	nodeAddress      string
	validatorAddress string
	KnownNodes       = []string{"localhost:3000"}
	blocksInTransit  []*blockchain.Block // []Block 타입으로 선언
	tempBlockList    []*blockchain.Block
	globalChainId    string

	// sync block list
	newBlockListChan = make(chan bool) // blocksInTransit에 새 블록이 추가되면 알림
	// mining block
	miningBlockChan = make(chan *blockchain.Block) // KnownNodes 새 노드가 추가되면 알림
	// syunc
	syncChan = make(chan bool, 1)

	isSync            = false
	isAppendBlockList = false
	mu                sync.Mutex // blocksInTransit 배열 동시 접근 제어
)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block    []byte
}
type BlockList struct {
	AddrFrom string
	Blocks   [][]byte
	Length   int
}

// 블록 요청을 위한 데이터 구조
type GetBlocks struct {
	AddrFrom string
}

// 특정 데이터 요청을 위한 데이터 구조
type LatestBlockHeight struct {
	AddrFrom string
	Type     string
	ID       []byte
}

// 인벤토리(블록/트랜잭션) 정보를 저장
type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type Version struct {
	Version    int
	BestHeight int64
	AddrFrom   string
}

// 명령어를 바이트 배열로 변환
func CmdToBytes(cmd string) []byte {
	// 명령어 길이에 맞는 바이트 배열 생성
	var bytes [commandLength]byte

	// 명령어 문자열을 바이트 배열로 전환
	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// 바이트 배열을 명령어 문자열로 변환
func BytesToCmd(bytes []byte) string {
	var cmd []byte

	// 바이트 배열에서 명령어 부분만 추출
	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	// 추출한 명령어를 문자열로 변환
	return fmt.Sprintf("%s", cmd)
}

// KnownNode에 중복된 주소 제거하는 함수
func RemoveDuplicatesNodes(nodes []string) []string {
	unique := make(map[string]bool)
	result := []string{}

	for _, node := range nodes {
		if !unique[node] {
			unique[node] = true
			result = append(result, node)
		}
	}

	return result
}

// 요청에서 명령어 부분을 추출
func ExtractCmd(request []byte) []byte {
	// 요청 데이터에서 명령어 부분만 추출
	return request[:commandLength]
}

func StartServer(chain *blockchain.BlockChain, valiAddress string) {
	var bcNode blockchain.Node

	globalChainId = chain.ChainId
	validatorAddress = valiAddress

	newNode := bcNode.NewNode(valiAddress, config.GlobalConfig.Port)
	nodeAddress = newNode.GetIP()
	if !newNode.IsPublicIP(newNode.IP) {
		nodeAddress = fmt.Sprintf("localhost:%d", newNode.ListenPort)
	}
	fmt.Printf("Starting node server with nodeAddress: %s\n", nodeAddress)

	conn, err := net.Dial("tcp", nodeAddress)

	if err == nil {
		conn.Close()
		log.Printf("Error: Port %s is already in use.", nodeAddress)
		os.Exit(1)
	}

	log.Printf("Starting node server on %s", nodeAddress)

	ln, err := net.Listen("tcp", nodeAddress)
	if err != nil {
		log.Printf("Error occurred while starting server: %v", err)
		os.Exit(1)
	}

	defer ln.Close()
	defer chain.Database.Close()
	go CloseDB(chain)

	log.Printf("Node server successfully started on %s", nodeAddress)
	rpcErrorChan := make(chan error)

	go StartRPCServer(chain, rpcErrorChan, newNode)

	if nodeAddress != KnownNodes[0] {
		fmt.Printf("Sending version from node %s to master node %s\n", nodeAddress, KnownNodes[0])
		SendVersion(KnownNodes[0], chain)
	}

	go listenForNewBlocks(chain)

	// Main server loop to handle incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic()
		}
		go HandleConnection(conn, chain)

		select {
		case rpcErr := <-rpcErrorChan:
			if rpcErr != nil {
				log.Println("Error in RPC server:", rpcErr)
			}
		default:
			// No action needed in default case
		}
	}
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)

	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// 노드가 알려진 노드 목록에 있는지 확인하는 함수
func NodeIsKnown(addr string) bool {
	// 알려진 노드 목록을 순회
	for _, node := range KnownNodes {
		// 주소가 있는지 확인
		if node == addr {
			// 있으면 true 반환
			return true
		}
	}

	// 없으면 false 반환
	return false
}

func CloseDB(chain *blockchain.BlockChain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}

// sync
func monitorBlocksInTransit(chain *blockchain.BlockChain) {
	mu.Lock()
	defer mu.Unlock()
	isSync = true
	syncChan <- true

	blockchain.SortBlocksByHeight(blocksInTransit)

	for len(blocksInTransit) > 0 {
		block := blocksInTransit[0]
		blocksInTransit = blocksInTransit[1:] // 첫 번째 블록 제거
		fmt.Printf("블록 추가 작업 중... 블록높이: %d, 해시: %x\n", block.Height, block.Hash)

		chain.Mu.Lock()

		blockHeight := block.Height
		blockHash := block.Hash

		if blockHeight == 0 {
			batch := new(leveldb.Batch)
			batch.Put(blockHash, block.Serialize())
			batch.Put([]byte("lh"), blockHash)
			heightKey := []byte(fmt.Sprintf("height-%d", blockHeight))
			batch.Put(heightKey, blockHash)

			err := chain.Database.Write(batch, nil)
			blockchain.Handle(err)

			chain.LastHash = blockHash
		} else {
			chain.AddBlock(block)
		}

		chain.Mu.Unlock()
	}

	isSync = false
}

func listenForNewBlocks(chain *blockchain.BlockChain) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !isSync && len(chain.LastHash) > 0 {
		// 초기 mining 시작
		go mining.Run(ctx, chain, validatorAddress, miningBlockChan)
	}

	for {

		select {
		case <-syncChan:
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
		case miningBlock := <-miningBlockChan:
			cancel()
			ctx, cancel = context.WithCancel(context.Background())

			// 블록 전파
			for _, node := range KnownNodes {
				if node == nodeAddress {
					continue
				}
				fmt.Println("블록 전파 node: ", node)
				fmt.Println("Propagating block to node:", node)
				SendBlock(node, miningBlock)
			}

			chain.Mu.Lock()
			chain.AddBlock(miningBlock)
			chain.Mu.Unlock()

			if !isSync && len(chain.LastHash) > 0 {
				go mining.Run(ctx, chain, validatorAddress, miningBlockChan)
			}

		case <-newBlockListChan:
			// 새 블록 알림 수신 시 기존 mining 중단
			cancel()
			ctx, cancel = context.WithCancel(context.Background())

			// 새 블록 알림 수신 시 즉시 blocksInTransit 확인 및 add 블록 작업
			monitorBlocksInTransit(chain)

			// // 새로운 컨텍스트 생성 및 mining 재시작

			if !isSync && len(chain.LastHash) > 0 {
				go mining.Run(ctx, chain, validatorAddress, miningBlockChan)
			}

		}

	}
}

func SyncWithLongestChain(chain *blockchain.BlockChain, otherHeight int64, addr string) {
	if !isSync {
		isSync = true
		syncChan <- true

		chain.ResetDatabase()
		chain.LastHash = []byte{}
		chain.CurrentBlock = nil
		SendLatestBlockHeight(addr, 0, otherHeight)
	}
}

func SyncKnownNodes(addr string) {
	if !NodeIsKnown(addr) && addr != "" {
		KnownNodes = append(KnownNodes, addr)
		KnownNodes = RemoveDuplicatesNodes(KnownNodes)

		for _, node := range KnownNodes {
			if node == nodeAddress {
				continue
			}

			SendKnownNodes(node)
		}
		fmt.Println("KnownNodes: ", KnownNodes)
	}
}
