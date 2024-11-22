package nodecmd

import (
	"fmt"
	"strconv"

	"github.com/Kim-DaeHan/mining-chain/blockchain"
	"github.com/Kim-DaeHan/mining-chain/config"
	"github.com/Kim-DaeHan/mining-chain/network"
	"github.com/urfave/cli/v2"
)

var (
	InitDB = &cli.Command{
		Name:  "initDB",
		Usage: "Initialize database",
		Action: func(c *cli.Context) error {
			chainId := strconv.Itoa(config.GlobalConfig.ChainId)
			validatorAddress := c.String("validator")
			blockchain.InitBlockChain(validatorAddress, chainId)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "validator", Usage: "Set validator address"},
		},
	}
	Start = &cli.Command{
		Name:  "start",
		Usage: "Start the xphere node",
		Action: func(c *cli.Context) error {
			chainId := strconv.Itoa(config.GlobalConfig.ChainId)
			validatorAddress := c.String("validator")
			chain := blockchain.ContinueBlockChain(chainId)
			defer chain.Database.Close()
			network.StartServer(chain, validatorAddress)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "validator", Usage: "Set validator address"},
		},
	}

	// Define other commands similarly
	CreateBlockchain = &cli.Command{
		Name:  "createBlockchain",
		Usage: "Create blockchain with a validator",
		Action: func(c *cli.Context) error {
			chainId := strconv.Itoa(config.GlobalConfig.ChainId)
			validatorAddress := c.String("address")
			if validatorAddress == "" {
				fmt.Println("Validator address is required")
				return nil
			}
			chain := blockchain.InitBlockChain(validatorAddress, chainId)
			defer chain.Database.Close()
			fmt.Println("Blockchain created successfully")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address", Usage: "Set validator address"},
		},
	}

	GenesisProofBlock = &cli.Command{
		Name:  "genesisProofBlock",
		Usage: "Create a genesis proof block",
		Action: func(c *cli.Context) error {
			chainId := strconv.Itoa(config.GlobalConfig.ChainId)
			validatorAddress := c.String("address")
			chain := blockchain.ContinueBlockChain(chainId)
			defer chain.Database.Close()
			block := blockchain.Genesis(validatorAddress)
			chain.AddBlock(block)
			fmt.Println("Genesis block created")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address", Usage: "Set validator for genesis block"},
		},
	}
)
