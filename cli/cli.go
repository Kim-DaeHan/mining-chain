package cli

import (
	"sort"

	"github.com/Kim-DaeHan/mining-chain/cli/utils/nodecmd"
	"github.com/urfave/cli/v2"
)

func InitializeApp() *cli.App {
	app := &cli.App{
		Name:    "xphere-proofnode",
		Usage:   "CLI for managing the xphere-proofnode blockchain",
		Version: "1.0.0",
		Commands: []*cli.Command{
			nodecmd.InitDB,
			nodecmd.Start,
			nodecmd.CreateBlockchain,
			nodecmd.GenesisProofBlock,
			nodecmd.RPCCommands,
		},
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	return app
}
