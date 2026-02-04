package e2e

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cucumber/godog"
)

func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// Build the binary before running tests
		cmd := exec.Command("go", "build", "-o", "mux-session", "..")
		if output, err := cmd.CombinedOutput(); err != nil {
			panic(fmt.Sprintf("Failed to build mux-session: %v\nOutput: %s", err, string(output)))
		}
	})

	ctx.AfterSuite(func() {
		// Cleanup any leftover binaries if they exist
		if _, err := os.Stat("mux-session"); err == nil {
			exec.Command("rm", "mux-session").Run()
		}
	})
}
