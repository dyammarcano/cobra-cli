package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestGoldenAddCmd(t *testing.T) {
	viper.Set("useViper", true)
	viper.Set("license", "apache")
	command := &Command{
		CmdName:   "test",
		CmdParent: parentName,
		Project:   getProject(),
	}
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			t.Fatalf("could not remove path %s: %v", path, err)
		}
	}(command.AbsolutePath)

	assertNoErr(t, command.Project.Create())
	assertNoErr(t, command.Create())

	generatedFile := fmt.Sprintf("%s/cmd/%s.go", command.AbsolutePath, command.CmdName)
	goldenFile := fmt.Sprintf("testdata/%s.go.golden", command.CmdName)
	err := compareFiles(generatedFile, goldenFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateCmdName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"cmdName", "cmdName"},
		{"cmd_name", "cmdName"},
		{"cmd-name", "cmdName"},
		{"cmd______Name", "cmdName"},
		{"cmd------Name", "cmdName"},
		{"cmd______name", "cmdName"},
		{"cmd------name", "cmdName"},
		{"cmdName-----", "cmdName"},
		{"cmdname-", "cmdname"},
	}

	for _, testCase := range testCases {
		got := validateCmdName(testCase.input)
		if testCase.expected != got {
			t.Errorf("Expected %q, got %q", testCase.expected, got)
		}
	}
}
