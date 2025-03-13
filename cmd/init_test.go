package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/emmadal/gopm/pkg"
	"github.com/stretchr/testify/assert"
)

// TestInitCmd ensures the init command executes without error
func TestInitCmd(t *testing.T) {
	// Change working directory to tempDir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err, "should get working directory")

	assert.NoError(t, err, "should change working directory to tempDir")

	// Execute InitCmd
	cmd := InitCmd
	cmd.SetArgs([]string{}) // No arguments expected
	err = cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err, "InitCmd should execute without error")

	// Verify package.json was created
	packageJSONPath := filepath.Join(oldWd, "package.json")
	_, err = os.Stat(packageJSONPath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err, "package.json file should be created")

	// Verify package.json contents
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err, "should be able to read package.json file")

	var pkgJSON pkg.PackageJSON
	err = json.Unmarshal(data, &pkgJSON)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err, "should parse package.json successfully")
	assert.Equal(t, filepath.Base(oldWd), pkgJSON.Name, "package.json name should match directory name")
	assert.Equal(t, "1.0.0", pkgJSON.Version, "package.json should have default version")
	assert.Equal(t, "index.js", pkgJSON.Main, "package.json should have correct main file")

	if err := os.Remove("package.json"); err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err, "should remove package.json file")
}

// TestGetFolderName ensures folder name extraction works correctly
func TestGetFolderName(t *testing.T) {
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err, "should get working directory")
	assert.NoError(t, err, "should change working directory to tempDir")
	folderName, err := getFolderName()
	assert.NoError(t, err, "getFolderName should not return an error")
	assert.Equal(t, filepath.Base(oldWd), folderName, "getFolderName should return correct folder name")
}
