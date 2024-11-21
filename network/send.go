package network

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/Kim-DaeHan/mining-chain/blockchain"
)

// 알려진 노드 주소를 전송
func SendKnownNodes(addr string) {
	// 현재 알려진 노드 리스트를 Addr 구조체에 저장
	nodes := Addr{KnownNodes}
	// Addr 구조체를 GOB 인코딩하여 바이트 배열로 변환
	payload := GobEncode(nodes)
	// 'addr' 명령어와 인코딩된 데이터를 결합하여 요청 생성
	request := append(CmdToBytes("knownNodes"), payload...)

	// 특정 주소로 요청 데이터 전송
	SendData(addr, request)
}

// 블록 데이터를 전송
func SendBlock(addr string, b *blockchain.Block) {
	// 현재 노드 주소와 직렬화된 블록 데이터를 Block 구조체에 담음
	data := Block{nodeAddress, b.Serialize()}
	// Block 구조체를 GOB 인코딩하여 바이트 배열로 변환
	payload := GobEncode(data)
	// "block" 명령어와 인코딩된 데이터를 결합하여 요청 생성
	request := append(CmdToBytes("block"), payload...)

	// 특정 주소로 요청 데이터 전송
	SendData(addr, request)
}

// 데이터를 특정 주소로 전송
func SendData(addr string, data []byte) {
	if addr == "" {
		log.Println("Error: Target address is empty, cannot send data.")
		return
	}

	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", addr, err)
		var updatedNodes []string
		for _, node := range KnownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}
		KnownNodes = updatedNodes
		return
	}
	defer conn.Close()

	fmt.Printf("Sending data to %s\n", addr)
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Printf("Error while sending data to %s: %v\n", addr, err)
	} else {
		fmt.Printf("Data successfully sent to %s\n", addr)
	}
}

// 특정 데이터(블록 또는 트랜잭션)를 요청하는 함수
// Sends a request for a range of blocks by height
func SendLatestBlockHeight(addr string, startHeight, endHeight int64) {
	payload := GobEncode(LatestBlockHeight{AddrFrom: nodeAddress, ID: []byte(fmt.Sprintf("%d-%d", startHeight, endHeight))})
	request := append(CmdToBytes("latestBlockHeight"), payload...)

	fmt.Printf("Requesting blocks from height %d to %d from %s\n", startHeight, endHeight, addr)
	SendData(addr, request)
}

func SendBlockList(addr string, blocks [][]byte, length int) {
	if addr == "" {
		log.Println("Error: Target address for block list is empty.")
		return
	}

	data := BlockList{
		AddrFrom: nodeAddress,
		Blocks:   blocks,
		Length:   length,
	}

	payload := GobEncode(data)
	request := append(CmdToBytes("blocklist"), payload...)

	fmt.Printf("Sending blocklist to %s with %d blocks\n", addr, len(blocks))

	SendData(addr, request)
}

func SendVersion(addr string, chain *blockchain.BlockChain) {
	if nodeAddress == "" {
		log.Println("Error: nodeAddress is empty, cannot send version message.")
		return
	}

	bestHeight := chain.GetBestHeight()
	fmt.Println("bestHeight: ", bestHeight)

	payload := GobEncode(Version{
		Version:    version,
		BestHeight: bestHeight,
		AddrFrom:   nodeAddress, // Ensure AddrFrom is set here
	})
	request := append(CmdToBytes("version"), payload...)

	fmt.Printf("Sending version to %s with bestHeight: %d from %s\n", addr, bestHeight, nodeAddress)

	SendData(addr, request)
}
