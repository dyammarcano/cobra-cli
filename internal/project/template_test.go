package project

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up viper defaults
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache2")
	viper.SetDefault("projectName", "testApp")

	afs = afero.NewMemMapFs()

	// Run the tests
	code := m.Run()

	// Clean up
	viper.Reset()

	os.Exit(code)
}

func TestLicenseTemplate(t *testing.T) {
	if err := licenseTemplate(afs, "LICENSE"); err != nil {
		t.Errorf("Error creating license template: %v", err)
	}

	licenseText, err := afero.ReadFile(afs, "LICENSE")
	if err != nil {
		t.Errorf("Error reading license file: %v", err)
	}

	if string(licenseText) == "" {
		t.Errorf("License file is empty")
	}
}

func TestRootTemplate(t *testing.T) {
	userLicense := viper.GetString("license")
	p := Project{AppName: "testApp", Legal: findLicense(userLicense)}

	if err := rootTemplate(afs, "cmd/root.go", p); err != nil {
		t.Errorf("Error creating root template: %v", err)
	}

	rootText, err := afero.ReadFile(afs, "cmd/root.go")
	if err != nil {
		t.Errorf("Error reading root file: %v", err)
	}

	goldenFile, err := os.ReadFile("testdata/root.go.golden")
	if err != nil {
		t.Errorf("Error reading golden file: %v", err)
	}

	if err := CompareContent(rootText, goldenFile); err != nil {
		t.Errorf("Error comparing files: %v", err)
	}
}

func TestMainTemplate(t *testing.T) {
	userLicense := viper.GetString("license")
	p := Project{AppName: "testApp", PkgName: "github.com/acme/myproject", Legal: findLicense(userLicense)}

	if err := mainTemplate(afs, "main.go", p); err != nil {
		t.Errorf("Error creating main template: %v", err)
	}

	mainText, err := afero.ReadFile(afs, "main.go")
	if err != nil {
		t.Errorf("Error reading main file: %v", err)
	}

	goldenFile, err := os.ReadFile("testdata/main.go.golden")
	if err != nil {
		t.Errorf("Error reading golden file: %v", err)
	}

	if err := CompareContent(mainText, goldenFile); err != nil {
		t.Errorf("Error comparing files: %v", err)
	}
}

func TestAddCommandTemplate(t *testing.T) {
	userLicense := viper.GetString("license")
	p := &Project{AppName: "testApp", PkgName: "github.com/acme/myproject", Legal: findLicense(userLicense)}
	c := Command{
		CmdName:   "add",
		CmdParent: "root",
		Project:   p,
	}

	if err := addCommandTemplate(afs, "cmd/add_command.go", c); err != nil {
		t.Errorf("Error creating add command template: %v", err)
	}

	addCommandText, err := afero.ReadFile(afs, "cmd/add_command.go")
	if err != nil {
		t.Errorf("Error reading add command file: %v", err)
	}

	goldenFile, err := os.ReadFile("testdata/add_command.go.golden")
	if err != nil {
		t.Errorf("Error reading golden file: %v", err)
	}

	if err := CompareContent(addCommandText, goldenFile); err != nil {
		t.Errorf("Error comparing files: %v", err)
	}
}
