package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kuberhealthy/kuberhealthy/v3/pkg/checkclient"
	log "github.com/sirupsen/logrus"
)

const (
	reportFailureEnv   = "REPORT_FAILURE"
	reportDelayEnv     = "REPORT_DELAY"
	timeLimitSkew      = time.Second * 5
	defaultTimeLimit   = time.Minute * 10
	defaultReportDelay = time.Second * 5
)

// CheckConfig stores the configuration used to run the test check.
type CheckConfig struct {
	ReportFailure bool
	ReportDelay   time.Duration
	TimeLimit     time.Duration
}

// loadConfig reads the environment variables and builds the check configuration.
func loadConfig() (*CheckConfig, error) {
	// Parse the report failure flag.
	reportFailure, err := parseReportFailure()
	if err != nil {
		return nil, err
	}

	// Parse the report delay duration.
	reportDelay, err := parseReportDelay()
	if err != nil {
		return nil, err
	}

	// Resolve the check time limit from the deadline.
	timeLimit, err := getTimeLimit()
	if err != nil {
		log.Warnln("There was an issue getting the check deadline:", err.Error())
	}

	// Build the final configuration.
	cfg := &CheckConfig{
		ReportFailure: reportFailure,
		ReportDelay:   reportDelay,
		TimeLimit:     timeLimit,
	}

	return cfg, nil
}

// parseReportFailure reads the REPORT_FAILURE environment variable.
func parseReportFailure() (bool, error) {
	// Read the raw environment value.
	reportFailureStr := os.Getenv(reportFailureEnv)

	// Default to false when unset.
	if len(reportFailureStr) == 0 {
		return false, nil
	}

	// Parse the boolean value.
	reportFailure, err := strconv.ParseBool(reportFailureStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse %s env var: %w", reportFailureEnv, err)
	}

	return reportFailure, nil
}

// parseReportDelay reads the REPORT_DELAY environment variable.
func parseReportDelay() (time.Duration, error) {
	// Read the raw environment value.
	reportDelayStr := os.Getenv(reportDelayEnv)

	// Default to the standard delay when unset.
	if len(reportDelayStr) == 0 {
		return defaultReportDelay, nil
	}

	// Parse the duration value.
	reportDelay, err := time.ParseDuration(reportDelayStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s env var: %w", reportDelayEnv, err)
	}

	return reportDelay, nil
}

// getTimeLimit reads the check deadline and returns a safe time limit.
func getTimeLimit() (time.Duration, error) {
	// Start with the default limit when the deadline is missing.
	timeLimit := defaultTimeLimit

	// Pull the deadline from the environment.
	deadline, err := checkclient.GetDeadline()
	if err != nil {
		return timeLimit, err
	}

	// Subtract a buffer to avoid reporting right at the deadline.
	timeLimit = deadline.Sub(time.Now().Add(timeLimitSkew))

	// Fall back to the default if the deadline is too close.
	if timeLimit <= 0 {
		return defaultTimeLimit, fmt.Errorf("check deadline is too soon to honor")
	}

	return timeLimit, nil
}
