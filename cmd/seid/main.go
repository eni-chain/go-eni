package main

import (
	"os"

	"github.com/eni-chain/go-eni/app/params"
	"github.com/eni-chain/go-eni/cmd/seid/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/eni-chain/go-eni/app"
)

func main() {
	params.SetAddressPrefixes()
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
