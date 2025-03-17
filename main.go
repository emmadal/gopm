/*
Copyright © 2025 Emmanuel Dalougou <emmanueldalougou@gmail.com>
*/
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/emmadal/gopm/cmd"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var root = &cobra.Command{
	Use:   "gopm",
	Short: "Fast, efficient package manager for JavaScript and TypeScript",
	Long:  "gopm is a package manager for JavaScript and TypeScript projects.\nIt is designed to be fast, efficient and easy to use and inspired by npm and yarn.",
	Version: strings.Join([]string{
		"v1.0.0",
		fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		"https://github.com/emmadal/gopm",
	}, "\n"),
}

func main() {
	root.AddCommand(cmd.InitCmd, cmd.AddCmd)

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute command: %v\n", err)
		os.Exit(1)
	}
}
