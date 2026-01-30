package e2e

import (
	"os"
	"os/exec"

	"github.com/cucumber/godog"
)

func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// Setup before running the test suite
	})

	ctx.AfterSuite(func() {
		// Cleanup any leftover binaries if they exist
		if _, err := os.Stat("mux-session"); err == nil {
			exec.Command("rm", "mux-session").Run()
		}
	})
}
