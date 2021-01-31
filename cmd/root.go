package main

import (
	"github.com/go-clix/cli"

	"fmt"
	"os"
)

func main() {
	rootCmd := &cli.Command{
		Use:   "crypt",
		Short: "Implementation of https://ieeexplore.ieee.org/document/7420966",
	}

	// add the child command
	rootCmd.AddCommand(
		encryptCmd(),
		decryptCmd(),
	)

	// run and check for errors
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
