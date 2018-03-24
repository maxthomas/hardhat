package main

import (
	"fmt"
	"os"

	"github.com/maxthomas/hardhat/cmd/hardhat/app"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
