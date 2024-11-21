package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Kim-DaeHan/mining-chain/cli"
	"github.com/Kim-DaeHan/mining-chain/config"
)

func main() {

	defer os.Exit(0)
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Config file is not exist tmp/config.json")
	}

	app := cli.InitializeApp()
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
