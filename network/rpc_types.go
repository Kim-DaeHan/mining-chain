package network

import (
	"math/big"

	"github.com/Kim-DaeHan/mining-chain/blockchain"
)

// GetBlockNumber
type GetBlockNumberArgs struct {
	ChainId string
}

type GetBlockNumberRes struct {
	Height int64
}

// GetBestHeight
type GetBestHeightArgs struct{}

type GetBestHeightRes struct {
	Height int64
}

// // GetBlock
type GetBlockArgs struct {
	Hash string
}

type GetBlockRes struct {
	Block blockchain.Block
}

// PrintChain
type PrintChainArgs struct {
	Block blockchain.Block
	Pow   string
}

// GetLastBlockHash
type GetLastBlockHashArgs struct{}

type GetLastBlockHashRes struct {
	Hash string
}

// GetBlockHashes
type GetBlockHashesArgs struct{}

type GetBlockHashesRes struct {
	Hash []string
}

// GetBlockList
type GetBlockListArgs struct{}

type GetBlockListRes struct {
	Block []blockchain.Block
}

// GetWork
type GetWorkArgs struct{}

type GetWorkRes struct {
	CurrentPowHash  string `json:"currentPowHash"`
	SeedHash        string `json:"seedHash"`
	TargetThreshold string `json:"targetThreshold"`
}

// GetHashRate
type GetHashRateArgs struct{}

type GetHashRateRes struct {
	Hashrate int `json:"hashrate"`
}

// Coinbase
type CoinbaseArgs struct{}

type CoinbaseRes struct {
	CoinbaseAddress string `json:"coinbaseAddress"`
}

// Mining
type MiningArgs struct{}

type MiningRes struct {
	IsMining bool `json:"isMining"`
}

// AddPeer
type AddPeerArgs struct {
	PeerAddress string `json:"peerAddress"`
}

type AddPeerRes struct {
	Success bool `json:"success"`
}

// GetDataDir
type GetDataDirArgs struct{}

type GetDataDirRes struct {
	DataDirectory string `json:"dataDirectory"`
}

// GetNodeInfo
type GetNodeInfoArgs struct{}

type GetNodeInfoRes struct {
	Enode      string `json:"enode"`
	ID         string `json:"id"`
	IP         string `json:"ip"`
	ListenPort int    `json:"listenPort"`
	Validator  string `json:"validator"`
	Name       string `json:"name"`
}

// GetPeer
type GetPeerArgs struct{}

type Peer struct {
	ID         string `json:"id"`
	IP         string `json:"ip"`
	Name       string `json:"name"`
	ListenPort int    `json:"listenPort"`
}
type GetPeerRes struct {
	Peers []Peer `json:"peers"`
}

// RemovePeer
type RemovePeerArgs struct {
	PeerAddress string `json:"peerAddress"`
}

type RemovePeerRes struct {
	Success bool `json:"success"`
}

// SetXpbase
type SetXpbaseArgs struct {
	Address string `json:"address"`
}

type SetXpbaseRes struct {
	Success bool `json:"success"`
}

// GetNodeHashRate
type GetNodeHashRateArgs struct{}

type GetNodeHashRateRes struct {
	Hashrate int `json:"hashrate"`
}

type GetDifficultyArgs struct {
	Height int64
}

type GetDifficultyRes struct {
	Difficulty *big.Int
}
