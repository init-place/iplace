package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	minitiaapp "github.com/init-place/iplace/app"
)

func main() {
	rootCmd, _ := NewRootCmd()

	if err := svrcmd.Execute(rootCmd, minitiaapp.EnvPrefix, minitiaapp.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
