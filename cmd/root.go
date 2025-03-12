/*
Copyright © 2025 Emmanuel Dalougou <emmanueldalougou@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopm",
	Short: "Package manager for JavaScript and TypeScript",
	Long: `Fast, efficient and easy to use package manager 🚀

Usage: gopm <command>

How to use:

  gopm add <package> - Install a package and any packages that it depends on.
  gopm install <package> - Install all packages from package.json.
  gopm dev <package> - Install a package in development mode.
  gopm rm <package> - Uninstall a package from node_modules and package.json.
  gopm up <package> - Update a package in node_modules.
  gopm list - List installed packages
  gopm init - Initialize a new project
  gopm help - Display help information
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
