package main

import (
	"os"
	"os/exec"

	"github.com/gitflush/utils"
)

var logger = utils.SetupLogger()

func main() {
	_, err := os.Getwd()
	if err != nil {
		logger.Panic("Failed to read working directory")
	}
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		logger.Panic("Failed to execute git diff --cached")
	}
	print("Git diff is", string(output))
}
