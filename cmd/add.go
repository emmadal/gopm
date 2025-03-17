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

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new dependency",
	Long:  "Add a new dependency to the project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Expect one or more dependencies\n")
		}
		return nil
	},
	Example: strings.Join([]string{
		"$ gopm add lodash",
		"$ gopm add react react-dom",
	}, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		exists, err := verifyJsonFile()
		if err != nil || !exists {
			logrus.Errorln("package.json file not found. Run 'gopm init' or 'gopm init my-module' to create one")
			os.Exit(0)
		} else {
			if err := createNodeModulesFolder(); err != nil {
				logrus.Errorln("failed to create node_modules folder")
				os.Exit(0)
			} else {
				if err := fetchDependencies(args); err != nil {
					logrus.Errorln("failed to fetch dependencies")
					os.Exit(0)
				} else {
					fmt.Printf("🍺 dependencies added successfully in %s\n\n", time.Since(start))
				}
			}
		}
	},
}

// createNodeModulesFolder creates a node_modules folder
func createNodeModulesFolder() error {
	cwd := pkg.GetCwd()
	modulesPath := filepath.Join(cwd, pkg.NODE_MODULE)

	if _, err := os.Stat(modulesPath); os.IsNotExist(err) {
		if err := os.Mkdir(modulesPath, 0755); err != nil {
			return fmt.Errorf("Unable to create dependencies folder: %v\n", err)
		}
	} else {
		logrus.Infoln("dependencies folder already exists. Skipping...")
	}
	return nil
}

// verifyJsonFile verifies if package.json file exists
func verifyJsonFile() (bool, error) {
	cwd := pkg.GetCwd()
	packageJsonPath := filepath.Join(cwd, pkg.PACKAGE_JSON)
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return false, fmt.Errorf("package.json file not found")
	}
	return true, nil
}

// fetchDependencies fetches dependencies from the npm registry
func fetchDependencies(args []string) error {
	logrus.Infof("Ready to download %d dependencies\n\n", len(args))
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
			version, err := body.GetDependencyLatest(dependency)
			if err != nil || version == "" {
				return err
			}
			if !body.DownloadPackage(dependency, version) {
				return err
			}
			mu.Lock()
			packageJson.AddDependency(map[string]string{dependency: version})
			added.Store(true) // Atomic write
			mu.Unlock()
			return nil
		})
	}

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error downloading dependencies: %w", err)
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
