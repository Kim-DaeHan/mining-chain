package nodecmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"

	"github.com/Kim-DaeHan/mining-chain/config"
	"github.com/Kim-DaeHan/mining-chain/network"
	"github.com/urfave/cli/v2"
)

var RPCCommands = &cli.Command{
	Name:  "rpc",
	Usage: "RPC commands for managing the blockchain node",
	Subcommands: []*cli.Command{
		GetBlockNumber, GetBlockList, GetLastBlockHash,
		GetBlock, GetBlockHashes,
		GetWork, GetHashRate, Coinbase, IsMining, AddPeer,
		GetDataDir, GetNodeInfo, GetPeer, RemovePeer,
		SetXpbase, GetNodeHashRate, GetDifficulty,
	},
}
var (
	GetBlockNumber = &cli.Command{
		Name:  "getBlockNumber",
		Usage: "Get the current block number",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Fatalf("Error connecting to RPC: %v", err)
			}
			defer client.Close()

			req := network.GetBlockNumberArgs{}
			var res network.GetBlockNumberRes
			if err = client.Call("RPCServer.GetBlockNumber", req, &res); err != nil {
				log.Fatalf("Error calling GetBlockNumber: %v", err)
			}
			fmt.Printf("Block Number: %d\n", res.Height)
			return nil
		},
	}
	GetBlockList = &cli.Command{
		Name:  "getBlockList",
		Usage: "Retrieve and display the list of all blocks",
		Action: func(c *cli.Context) error {
			fmt.Println("Get All Block List")

			rpchost := fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort)
			fmt.Printf("cli request getBlockList rpchost : %s\n", rpchost)

			client, err := rpc.DialHTTP("tcp", rpchost)
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}

			req := network.GetBlockListArgs{}
			var res network.GetBlockListRes

			err = client.Call("RPCServer.GetBlockList", req, &res)
			if err != nil {
				log.Panic("Error calling RPC:", err)
			}

			for _, block := range res.Block {
				fmt.Println("Block: ", block)
			}

			return nil
		},
	}

	GetLastBlockHash = &cli.Command{
		Name:  "getLastBlockHash",
		Usage: "Get the hash of the last block",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "block hash", Usage: "Provide block hash you want block"},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Fatalf("Error connecting to RPC: %v", err)
			}
			defer client.Close()

			req := network.GetLastBlockHashArgs{}
			var res network.GetLastBlockHashRes
			if err = client.Call("RPCServer.GetLastBlockHash", req, &res); err != nil {
				log.Fatalf("Error calling GetLastBlockHash: %v", err)
			}
			fmt.Printf("Last Block Hash: %s\n", res.Hash)
			return nil
		},
	}

	GetBestHeight = &cli.Command{
		Name:  "getBestHeight",
		Usage: "Get the best block height",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "block hash", Usage: "Provide block hash you want block"},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Fatalf("Error connecting to RPC: %v", err)
			}
			defer client.Close()

			req := network.GetBestHeightArgs{}
			var res network.GetBestHeightRes
			if err = client.Call("RPCServer.GetBestHeight", req, &res); err != nil {
				log.Fatalf("Error calling GetBestHeight: %v", err)
			}
			fmt.Printf("Best Height: %d\n", res.Height)
			return nil
		},
	}

	GetBlock = &cli.Command{
		Name:  "getBlock",
		Usage: "Get block details by hash",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "hash", Usage: "Hash of the block to retrieve"},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Fatalf("Error connecting to RPC: %v", err)
			}
			defer client.Close()

			req := network.GetBlockArgs{Hash: c.String("hash")}
			var res network.GetBlockRes
			if err = client.Call("RPCServer.GetBlock", req, &res); err != nil {
				log.Fatalf("Error calling GetBlock: %v", err)
			}
			fmt.Printf("Block Hash: %x\nNonce: %d\n", res.Block.Hash, res.Block.Nonce)
			return nil
		},
	}

	GetBlockHashes = &cli.Command{
		Name:  "getBlockHashes",
		Usage: "Get all block hashes",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Fatalf("Error connecting to RPC: %v", err)
			}
			defer client.Close()

			req := network.GetBlockHashesArgs{}
			var res network.GetBlockHashesRes
			if err = client.Call("RPCServer.GetBlockHashes", req, &res); err != nil {
				log.Fatalf("Error calling GetBlockHashes: %v", err)
			}

			for i, hash := range res.Hash {
				fmt.Printf("Block %d: %s\n", i, hash)
			}
			return nil
		},
	}

	GetWork = &cli.Command{
		Name:  "getWork",
		Usage: "Retrieve mining work info",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Fatalf("Error connecting to RPC: %v", err)
			}
			defer client.Close()

			req := network.GetWorkArgs{}
			var res network.GetWorkRes
			if err = client.Call("RPCServer.GetWork", req, &res); err != nil {
				log.Fatalf("Error calling GetWork: %v", err)
			}
			workInfoJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("Work Info:", string(workInfoJSON))
			return nil
		},
	}

	GetHashRate = &cli.Command{
		Name:  "getHashRate",
		Usage: "Get current hash rate",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Fatalf("Error connecting to RPC: %v", err)
			}
			defer client.Close()

			req := network.GetHashRateArgs{}
			var res network.GetHashRateRes
			if err = client.Call("RPCServer.GetHashRate", req, &res); err != nil {
				log.Fatalf("Error calling GetHashRate: %v", err)
			}
			hashRateJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("HashRate:", string(hashRateJSON))
			return nil
		},
	}
	Coinbase = &cli.Command{
		Name:  "coinbase",
		Usage: "Get Coinbase address",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.CoinbaseArgs{}
			var res network.CoinbaseRes
			if err = client.Call("RPCServer.Coinbase", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			coinbaseJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("Coinbase:", string(coinbaseJSON))
			return nil
		},
	}

	IsMining = &cli.Command{
		Name:  "isMining",
		Usage: "Check if node is currently mining",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.MiningArgs{}
			var res network.MiningRes
			if err = client.Call("RPCServer.Mining", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			miningJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("isMining:", string(miningJSON))
			return nil
		},
	}

	AddPeer = &cli.Command{
		Name:  "addPeer",
		Usage: "Add a new peer to the network",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address", Usage: "Address of peer to add", Required: true},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.AddPeerArgs{PeerAddress: c.String("address")}
			var res network.AddPeerRes
			if err = client.Call("RPCServer.AddPeer", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			addPeerJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("add peer result:", string(addPeerJSON))
			return nil
		},
	}

	GetDataDir = &cli.Command{
		Name:  "getDataDir",
		Usage: "Get the data directory path",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.GetDataDirArgs{}
			var res network.GetDataDirRes
			if err = client.Call("RPCServer.GetDataDir", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			getDataDirJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("dataDirectory:", string(getDataDirJSON))
			return nil
		},
	}

	GetNodeInfo = &cli.Command{
		Name:  "getNodeInfo",
		Usage: "Get information about the node",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.GetNodeInfoArgs{}
			var res network.GetNodeInfoRes
			if err = client.Call("RPCServer.GetNodeInfo", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			getNodeInfoJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("nodeInfo:", string(getNodeInfoJSON))
			return nil
		},
	}

	GetPeer = &cli.Command{
		Name:  "getPeer",
		Usage: "Get information about connected peers",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.GetPeerArgs{}
			var res network.GetPeerRes
			if err = client.Call("RPCServer.GetPeer", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			getPeerJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("peers:", string(getPeerJSON))
			return nil
		},
	}

	RemovePeer = &cli.Command{
		Name:  "removePeer",
		Usage: "Remove a peer from the network",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address", Usage: "Address of peer to remove", Required: true},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.RemovePeerArgs{PeerAddress: c.String("address")}
			var res network.RemovePeerRes
			if err = client.Call("RPCServer.RemovePeer", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			removePeerJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("remove peer result:", string(removePeerJSON))
			return nil
		},
	}

	SetXpbase = &cli.Command{
		Name:  "setXpbase",
		Usage: "Set the XPBase address",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address", Usage: "Address to set as XPBase", Required: true},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.SetXpbaseArgs{Address: c.String("address")}
			var res network.SetXpbaseRes
			if err = client.Call("RPCServer.SetXpbase", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			setXpbaseJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("set xpbase result:", string(setXpbaseJSON))
			return nil
		},
	}

	GetNodeHashRate = &cli.Command{
		Name:  "getNodeHashRate",
		Usage: "Get the node's hash rate",
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.GetNodeHashRateArgs{}
			var res network.GetNodeHashRateRes
			if err = client.Call("RPCServer.GetNodeHashRate", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			nodeHashRateJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("NodeHashRate:", string(nodeHashRateJSON))
			return nil
		},
	}

	GetDifficulty = &cli.Command{
		Name:  "getDifficulty",
		Usage: "Get the difficulty at a specific height",
		Flags: []cli.Flag{
			&cli.Int64Flag{Name: "height", Usage: "Height to check difficulty", Required: true},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", config.GlobalConfig.RPCPort))
			if err != nil {
				log.Panic("Error dialing RPC:", err)
			}
			defer client.Close()

			req := network.GetDifficultyArgs{Height: c.Int64("height")}
			var res network.GetDifficultyRes
			if err = client.Call("RPCServer.GetDifficulty", req, &res); err != nil {
				log.Panic("Error calling RPC:", err)
			}

			difficultyJSON, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println("Difficulty:", string(difficultyJSON))
			return nil
		},
	}
)
