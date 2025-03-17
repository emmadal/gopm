/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/emmadal/gopm/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// DevCmd represents the dev command
var DevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Install a package in development mode",
	Long:  `Install a package in development mode for use in the project`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Expect one or more dependencies\n")
		}
		return nil
	},
	Example: strings.Join([]string{
		"$ gopm dev @types/node",
	}, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		exists, err := pkg.VerifyJsonFile()
		if err != nil || !exists {
			logrus.Errorln("package.json file not found. Run 'gopm init' or 'gopm init my-module' to create one")
			os.Exit(1)
		}
		if err := pkg.CreateNodeModulesFolder(); err != nil {
			logrus.Errorln("failed to create node_modules folder")
			os.Exit(1)
		}
		if err := fetchDevDependencies(args); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("🍺 Dependencies added successfully in %s\n\n", time.Since(start))
	},
}

// fetchDevDependencies fetches dependencies from the npm registry
func fetchDevDependencies(args []string) error {
	logrus.Infof("Ready to download %d dev dependencies\n\n", len(args))
	cwd := pkg.GetCwd()
	packageJsonPath := filepath.Join(cwd, pkg.PACKAGE_JSON)

	// Open package.json file in read mode
	file, err := os.Open(packageJsonPath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Stream JSON decoding (avoids full memory load)
	var packageJson pkg.PackageJSON
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&packageJson); err != nil {
		return fmt.Errorf("error decoding package.json: %w", err)
	}

	// Create errgroup to limit concurrent downloads
	g := errgroup.Group{}
	if len(args) > pkg.MAX_CONCURRENT_DOWNLOADS {
		g.SetLimit(pkg.MAX_CONCURRENT_DOWNLOADS)
	} else {
		g.SetLimit(len(args))
	}

	var mu sync.Mutex // Protect concurrent writes
	var added atomic.Bool

	// Download dependencies concurrently
	for _, dependency := range args {
		body := pkg.BodyRegistery{}
		g.Go(func() error {
			// Get the latest version of the dependency
			version, err := body.GetDependencyLatest(dependency)
			if err != nil || version == "" {
				return err
			}
			// Download the package
			if !body.DownloadPackage(dependency, version) {
				return err
			}
			// Unzip the downloaded
			dependencyPath := filepath.Join(cwd, pkg.NODE_MODULE, dependency)
			if err := pkg.UnzipDependency(dependencyPath); err != nil {
				return err
			}
			// Add dependency to package.json
			mu.Lock()
			packageJson.AddDevDependency(map[string]string{dependency: version})
			added.Store(true) // Atomic write
			mu.Unlock()
			return nil
		})
	}

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		return err
	}

	if !added.Load() {
		logrus.Infoln("❌ No dependencies added. Skipping file write.")
		return nil
	}

	// Write updated dependencies to a temporary file
	tempFilePath := packageJsonPath + ".tmp"
	tempFile, err := os.OpenFile(tempFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening temp file: %w", err)
	}
	defer tempFile.Close()

	// Ensure cleanup if an error occurs
	defer os.Remove(tempFilePath)

	// Stream JSON encoding (efficient memory usage)
	encoder := json.NewEncoder(tempFile)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", pkg.INDENT)
	if err := encoder.Encode(&packageJson); err != nil {
		return fmt.Errorf("error encoding package.json: %w", err)
	}

	// Replace original file (atomic operation to prevent corruption)
	if err := os.Rename(tempFilePath, packageJsonPath); err != nil {
		return fmt.Errorf("error replacing file: %w", err)
	}
	return nil
}
