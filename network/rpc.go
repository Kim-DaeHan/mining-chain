package network

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/Kim-DaeHan/mining-chain/blockchain"
	"github.com/Kim-DaeHan/mining-chain/config"
)

type RPCServer struct {
	Port  int
	chain *blockchain.BlockChain
	node  *blockchain.Node
}

func (r *RPCServer) GetBlockNumber(req *GetBlockNumberArgs, res *GetBlockNumberRes) error {
	// chain := blockchain.ContinueBlockChain(req.ChainId)
	// defer chain.Database.Close()

	bestHeight := r.chain.GetBestHeight()
	res.Height = bestHeight
	return nil
}

func (r *RPCServer) GetBestHeight(req *GetBestHeightArgs, res *GetBestHeightRes) error {
	bestHeight := r.chain.GetBestHeight()
	res.Height = bestHeight
	return nil
}

func (r *RPCServer) GetLastBlockHash(req *GetLastBlockHashArgs, res *GetLastBlockHashRes) error {
	blockHash := r.chain.GetLastBlockHash()
	res.Hash = fmt.Sprintf("%x", blockHash)
	return nil
}

func (r *RPCServer) GetBlock(req *GetBlockArgs, res *GetBlockRes) error {
	if req.Hash == "" {
		log.Panic("-hash option")
	}

	hashBytes, err := hex.DecodeString(req.Hash)

	if err != nil {
		log.Panic(err)
	}

	block, err := r.chain.GetBlock(hashBytes)

	res.Block = block

	if err != nil {
		fmt.Printf("여기구나!")
		log.Panic(err)
	}

	return nil
}

func (r *RPCServer) GetBlockHashes(req *GetBlockHashesArgs, res *GetBlockHashesRes) error {
	blockHashes := r.chain.GetBlockHashes()

	for _, hash := range blockHashes {
		hashString := hex.EncodeToString(hash)
		res.Hash = append(res.Hash, hashString)
	}

	return nil
}

func (r *RPCServer) GetBlockList(req *GetBlockListArgs, res *GetBlockListRes) error {

	blockList := r.chain.GetBlockList()

	for _, block := range blockList {
		fmt.Printf("block height GetBlockList :: >> %d\n", block.Height)
		res.Block = append(res.Block, block)
	}

	return nil
}

// 작업증명 기반 블록체인에서 작업의 난이도와 작업을 완료하기 위한 해시값을 제공하는 JSON-RPC 메서드(마이너가 다음 블록을 채굴하기 위해 필요한 정보 반환)
func (r *RPCServer) GetWork(req *GetWorkArgs, res *GetWorkRes) error {
	// 현재 작업해야 하는 블록의 상태를 나타내는 해시(마이너가 이 값을 대상으로 작업 수행)
	res.CurrentPowHash = "0x5eab5d6d3e47adf8d0d4ae9a4b96be1f24acb8b3f1b3cfa7cde7d77eae75a4a8"
	// 난이도를 설정하는데 필요한 시드 값
	res.SeedHash = "0x7e3f9c19b0d3ea84363e7fd32158ef87688b1bb1a5ed0a5d003a3f54e9a2a5a3"
	// 네트워크 난이도에 따라 작업이 완료되기 위해 충족해야 할 해시 목표값
	res.TargetThreshold = "0x0000000000000000000000000000000000000000000000000000000000001abc"
	return nil
}

// 전체 네트워크의 평균적인 해시레이트(초당 해시 계산 속도)를 조회하는 JSON-RPC 메서드
func (r *RPCServer) GetHashRate(req *GetHashRateArgs, res *GetHashRateRes) error {
	// 초당 1000개의 해시를 계산
	res.Hashrate = 1000
	return nil
}

// 마이너가 블록을 채굴했을 때 보상을 받을 계정을 조회하는 JSON-RPC 메서드
func (r *RPCServer) Coinbase(req *CoinbaseArgs, res *CoinbaseRes) error {
	res.CoinbaseAddress = "0x742d35cc6634c0532925a3b844bc454e4438f44e"
	return nil
}

// 노드 마이닝 여부를 조회하는 JSON-RPC 메서드
func (r *RPCServer) Mining(req *MiningArgs, res *MiningRes) error {
	res.IsMining = true
	return nil
}

// peer 추가하는 JSON-RPC 메서드
func (r *RPCServer) AddPeer(req *AddPeerArgs, res *AddPeerRes) error {
	fmt.Println("peer address: ", req.PeerAddress)
	res.Success = true
	return nil
}

// datadir를 조회하는 JSON-RPC 메서드
func (r *RPCServer) GetDataDir(req *GetDataDirArgs, res *GetDataDirRes) error {
	res.DataDirectory = "/tmp/blocks_1001"
	return nil
}

// 노드 정보를 조회하는 JSON-RPC 메서드
func (r *RPCServer) GetNodeInfo(req *GetNodeInfoArgs, res *GetNodeInfoRes) error {
	res.Enode = "enode://abcdef1234567890@127.0.0.1:30303"
	res.ID = "abcdef1234567890"
	res.IP = r.node.IP
	res.ListenPort = r.node.ListenPort
	res.Validator = r.node.Validator
	res.Name = "Geth/v1.9.0-stable/linux-amd64/go1.12"
	return nil
}

// peer 정보를 조회하는 JSON-RPC 메서드
func (r *RPCServer) GetPeer(req *GetPeerArgs, res *GetPeerRes) error {
	peer1 := Peer{
		ID:         "abcdef1234567890",
		IP:         "192.168.1.10",
		Name:       "Geth/v1.9.0-stable/linux-amd64/go1.12",
		ListenPort: 30303,
	}
	res.Peers = append(res.Peers, peer1)

	peer2 := Peer{
		ID:         "abcdef23456",
		IP:         "192.168.1.11",
		Name:       "Geth/v1.9.0-stable1/linux-amd64/go1.12",
		ListenPort: 30303,
	}
	res.Peers = append(res.Peers, peer2)

	return nil
}

// peer 제거하는 JSON-RPC 메서드
func (r *RPCServer) RemovePeer(req *RemovePeerArgs, res *RemovePeerRes) error {
	fmt.Println("peer address: ", req.PeerAddress)
	res.Success = true
	return nil
}

// xp 보상 얻을 주소 설정하는 JSON-RPC 메서드
func (r *RPCServer) SetXpbase(req *SetXpbaseArgs, res *SetXpbaseRes) error {
	fmt.Println("xpbase address: ", req.Address)
	res.Success = true
	return nil
}

// 현재 노드의 해시레이트(초당 해시 계산 속도)를 조회하는 JSON-RPC 메서드
func (r *RPCServer) GetNodeHashRate(req *GetNodeHashRateArgs, res *GetNodeHashRateRes) error {
	// 초당 1000개의 해시를 계산
	res.Hashrate = 1000
	return nil
}

// 현재 노드의 해시레이트(초당 해시 계산 속도)를 조회하는 JSON-RPC 메서드
func (r *RPCServer) GetDifficulty(req *GetDifficultyArgs, res *GetDifficultyRes) error {
	// 초당 1000개의 해시를 계산
	diff := r.chain.Difficulty(req.Height)
	res.Difficulty = diff

	return nil
}

func StartRPCServer(chain *blockchain.BlockChain, rpcErrorChan chan error, node *blockchain.Node) {
	rpcServer := &RPCServer{config.GlobalConfig.RPCPort, chain, node}

	err := rpc.Register(rpcServer)
	if err != nil {
		rpcErrorChan <- err
		return
	}

	rpchost := fmt.Sprintf("localhost:%d", rpcServer.Port)

	rpcListener, err := net.Listen("tcp", rpchost)
	if err != nil {
		rpcErrorChan <- err
	}

	rpc.HandleHTTP()
	log.Printf("Serving RPC server on %s", rpchost)

	err = http.Serve(rpcListener, nil)
	if err != nil {
		rpcErrorChan <- err
	}
}
