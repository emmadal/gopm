package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emmadal/gopm/pkg"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var InitCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("Only one or zero argument expected\n")
		}
		return nil
	},
	Example: strings.Join([]string{
		"$ gopm init",
		"$ gopm init my-module",
	}, "\n"),
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		if len(args) == 1 {
			if err := createPackageJSONFile(args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "failed to initialize module: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := createPackageJSONFile(""); err != nil {
				fmt.Fprintf(os.Stderr, "failed to initialize module: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("module initialized successfully in %s\n", time.Since(start))
	},
	Short: "Initialize a new module",
	Long:  `Initialize a new module with a package.json file.`,
}

// createPackageJSONFile creates a package.json file
func createPackageJSONFile(moduleName string) error {
	packageName := ""

	if moduleName != "" || len(moduleName) > 0 {
		packageName = moduleName
	} else {
		name, err := getFolderName()
		if err != nil {
			return fmt.Errorf("failed to get folder name: %v\n", err)
		}
		packageName = name
	}

	// initialize package.json file
	jsonFile := pkg.PackageJSON{
		Name:    strings.TrimSpace(packageName),
		Version: "1.0.0",
		Main:    "index.js",
		Scripts: map[string]string{
			"test": "Echo \"Error: no test specified\" && exit 1",
		},
		Description:     "",
		Dependencies:    map[string]string{},
		DevDependencies: map[string]string{},
	}

	file, err := os.Create("package.json")

	if err != nil {
		return fmt.Errorf("failed to create package.json file: %v\n", err)
	}
	defer file.Close()

	// write package.json file
	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jsonFile); err != nil {
		return fmt.Errorf("failed to encode package.json file: %v\n", err)
	}
	return nil

}

// getFolderName returns the current folder name
func getFolderName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	arr := strings.Split(dir, string(os.PathSeparator))
	return arr[len(arr)-1], nil
}
