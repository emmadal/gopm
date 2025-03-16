package pkg

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

const (
	NPM_REGISTRY             = "https://registry.npmjs.org/"
	NODE_MODULE              = "node_modules"
	PACKAGE_JSON             = "package.json"
	INDENT                   = "  "
	MAX_CONCURRENT_DOWNLOADS = 20
)

// PackageJSON is a representation of a package.json file.
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// BodyRegistery is a representation of a response from the npm registry.
type BodyRegistery struct {
	Name     string `json:"name"`
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
}

type Dependency struct {
	Url     string `json:"url"`
	Version string `json:"version"`
}

// Tarball returns the tarball URL for a given package.
func Tarball(dependency, version string) string {
	if strings.HasPrefix(dependency, "@") {
		var module = strings.Split(dependency, "/")
		if err := NewDirectory(module[0]); err != nil {
			logrus.Errorf("Failed to create directory for %s: %v", dependency, err)
		}
		return fmt.Sprintf("%s%s/-/%s-%s.tgz", NPM_REGISTRY, dependency, module[1], version)
	}
	return fmt.Sprintf("%s%s/-/%s-%s.tgz", NPM_REGISTRY, dependency, dependency, version)
}

// AddDependency adds dependency to the package.json file.
func (p *PackageJSON) AddDependency(dependencies map[string]string) {
	if len(dependencies) == 0 {
		return
	}
	if p.Dependencies == nil {
		p.Dependencies = make(map[string]string)
	}
	maps.Copy(p.Dependencies, dependencies)
}

// DownloadPackage downloads a package from the npm registry.
func (b *BodyRegistery) DownloadPackage(dependency, version string) bool {
	tarball := Tarball(dependency, version)

	// Set timeout for HTTP request (e.g., 20 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tarball, nil)
	if err != nil {
		logrus.Errorf("Failed to create request for %s: %v", dependency, err)
		return false
	}

	// Send HTTP request
	client, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("Failed to download %s: %v", dependency, err)
		return false
	}
	defer client.Body.Close()

	// Get current working directory
	cwd := GetCwd()
	fileName := fmt.Sprintf("%s.tgz", dependency)
	filePath := filepath.Join(cwd, NODE_MODULE, fileName)

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("Failed to create file %s", filePath)
		return false
	}
	defer f.Close()

	if client.ContentLength > 0 {
		// Create progress bar
		bar := progressbar.DefaultBytes(client.ContentLength, fmt.Sprintf("Downloading %s...", dependency))
		// Copy the package to the node_modules folder
		if _, err := io.Copy(io.MultiWriter(bar, f), client.Body); err != nil {
			logrus.Errorf("Failed to copy %s: %v", dependency, err)
			return false
		}
	}

	logrus.Infof("Successfully downloaded %s\n\n", dependency)
	return true
}

// NewDirectory creates a new directory.
func NewDirectory(dirname string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	dirPath := filepath.Join(cwd, NODE_MODULE, dirname)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.Mkdir(dirPath, 0755)
	}
	return nil
}

// GetCwd returns the current working directory.
func GetCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("failed to get current working directory: %v", err)
		os.Exit(0)
	}
	return cwd
}
