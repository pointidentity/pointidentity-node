package main

import (
	"os"

	"github.com/pointidentity/pointidentity-node/app"
	"github.com/pointidentity/pointidentity-node/cmd/pointidentity-noded/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.Name, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
