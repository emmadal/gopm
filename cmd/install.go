package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emmadal/gopm/pkg"
	"github.com/spf13/cobra"
)

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install packages",
	Long:  `Install all packages from package.json`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("Expect no arguments\n")
		}
		return nil
	},
	Example: strings.Join([]string{
		"$ gopm install",
	}, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		err := getPackageJson()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

// getPackageJson gets the package.json file
func getPackageJson() error {
	start := time.Now()
	var p pkg.PackageJSON

	fileContent, err := p.ReadPackageJson()
	if err != nil {
		return err
	}

	if err := pkg.CreateNodeModulesFolder(); err != nil {
		return err
	}

	// Process Dependencies
	if deps, err := extractDependencies(fileContent.Dependencies); err == nil {
		if err := fetchDependencies(deps); err != nil {
			return err
		}
	} else {
		return err
	}

	// Process Dev Dependencies
	if devDeps, err := extractDependencies(fileContent.DevDependencies); err == nil {
		if err := fetchDevDependencies(devDeps); err != nil {
			return err
		}
	} else {
		return err
	}

	fmt.Printf("🍺 Install completed in %v\n", time.Since(start))
	return nil
}

// extractDependencies gets the dependencies from the package.json file
func extractDependencies(dependencies map[string]string) ([]string, error) {
	deps := make([]string, 0, len(dependencies)) // Preallocate slice capacity
	for dep := range dependencies {
		deps = append(deps, dep)
	}
	return deps, nil
}
