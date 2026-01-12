package main

import (
	"context"
	"os"
	"time"

	"github.com/kuberhealthy/kuberhealthy/v3/pkg/checkclient"
	nodecheck "github.com/kuberhealthy/kuberhealthy/v3/pkg/nodecheck"
	log "github.com/sirupsen/logrus"
)

const reportingURLEnv = "KH_REPORTING_URL"

// main loads configuration, waits for the delay, and reports the test result.
func main() {
	// Enable debug logging on the check client for parity with v2 behavior.
	checkclient.Debug = true

	// Enable nodecheck debug output for troubleshooting readiness.
	nodecheck.EnableDebugOutput()

	// Load configuration from environment variables.
	cfg, err := loadConfig()
	if err != nil {
		reportFailureAndExit(err)
		return
	}

	// Log the reporting URL for visibility.
	reportingURL := os.Getenv(reportingURLEnv)
	log.Infoln("Using kuberhealthy reporting url", reportingURL)

	// Wait before reporting to simulate delayed checks.
	log.Infoln("Waiting", cfg.ReportDelay, "before reporting...")
	time.Sleep(cfg.ReportDelay)

	// Enforce the check timeout based on the deadline.
	startTimeoutWatcher(cfg.TimeLimit)

	// Wait briefly for Kuberhealthy to be reachable before reporting.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = nodecheck.WaitForKuberhealthy(ctx)
	if err != nil {
		log.Warnln("Error waiting for kuberhealthy endpoint to be contactable by checker pod with error:", err.Error())
	}

	// Report the configured test result.
	err = reportCheckResult(cfg)
	if err != nil {
		log.Errorln("Error reporting to Kuberhealthy servers:", err.Error())
		return
	}

	// Log success for operators.
	log.Infoln("Successfully reported to Kuberhealthy servers")
}

// startTimeoutWatcher schedules a forced exit when the check time limit expires.
func startTimeoutWatcher(timeLimit time.Duration) {
	// Guard against invalid time limits.
	if timeLimit <= 0 {
		log.Warnln("Time limit is non-positive, skipping timeout watcher")
		return
	}

	// Log the enforced time limit.
	log.Infoln("Check time limit set to:", timeLimit)

	// Schedule a timeout exit to avoid hanging past the deadline.
	time.AfterFunc(timeLimit, exitOnTimeout)
}

// exitOnTimeout logs a timeout message and exits with failure.
func exitOnTimeout() {
	// Emit a timeout log entry.
	log.Errorln("Check took too long and timed out.")

	// Exit with a non-zero status to signal failure.
	os.Exit(1)
}

// reportCheckResult reports success or failure based on the configuration.
func reportCheckResult(cfg *CheckConfig) error {
	// Report a failure when configured to do so.
	if cfg.ReportFailure {
		log.Infoln("Reporting failure...")
		return checkclient.ReportFailure([]string{"Test has failed!"})
	}

	// Report success for default behavior.
	log.Infoln("Reporting success...")
	return checkclient.ReportSuccess()
}

// reportFailureAndExit reports a configuration error and exits.
func reportFailureAndExit(err error) {
	// Log the configuration error.
	log.Errorln(err)

	// Attempt to report the failure to Kuberhealthy.
	reportErr := checkclient.ReportFailure([]string{err.Error()})
	if reportErr != nil {
		log.Errorln("error when reporting to kuberhealthy:", reportErr.Error())
	}

	// Exit after reporting the failure.
	os.Exit(1)
}
