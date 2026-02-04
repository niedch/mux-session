package e2e

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestE2E(t *testing.T) {
	suite := godog.TestSuite{
		Name:                 "Mux Session E2E Tests",
		TestSuiteInitializer: InitializeSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"},
			StopOnFailure: true,
			Strict:        true,
			Concurrency: 4,
		},
	}

	status := suite.Run()
	if status > 0 {
		t.Fatalf("E2E tests failed with status: %d", status)
	}
}
