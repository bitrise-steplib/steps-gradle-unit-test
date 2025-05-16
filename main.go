package main

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-android/v2/cache"
	utilscache "github.com/bitrise-io/go-steputils/cache"
	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/kballard/go-shellquote"
)

type Inputs struct {
	GradleFile    string `env:"gradle_file"`
	UnitTestTasks string `env:"unit_test_task,required"`
	GradlewPath   string `env:"gradlew_file_path,required"`
	UnitTestFlags string `env:"unit_test_flags"`
	DeployDir     string `env:"BITRISE_DEPLOY_DIR"`
	CacheLevel    string `env:"cache_level,opt[all,only deps,none]"`
}

func main() {
	// Setup dependencies
	logger := log.NewLogger()
	envRepo := env.NewRepository()
	stepInputParser := stepconf.NewInputParser(envRepo)
	pathChecker := pathutil.NewPathChecker()
	cmdFactory := command.NewFactory(envRepo)
	outputExporter := export.NewExporter(cmdFactory)

	// Parse inputs
	var inputs Inputs
	if err := stepInputParser.Parse(&inputs); err != nil {
		failF(logger, "issue with input: %s", err)
	}

	stepconf.Print(inputs)
	logger.Println()

	if err := validateInputs(pathChecker, inputs); err != nil {
		failF(logger, "issue with input: %s", err)
	}

	// Run gradle task
	err := os.Chmod(inputs.GradlewPath, 0770)
	if err != nil {
		failF(logger, "Failed to add executable permission on gradlew file (%s): %s", inputs.GradlewPath, err)
	}

	fmt.Println()
	logger.Infof("Running gradle task...")
	if err := runGradleTask(cmdFactory, logger, inputs.GradlewPath, inputs.GradleFile, inputs.UnitTestTasks, inputs.UnitTestFlags); err != nil {
		logger.Errorf("Gradle task failed: %s", err)

		if err := outputExporter.ExportOutput("BITRISE_GRADLE_TEST_RESULT", "failed"); err != nil {
			logger.Warnf("Failed to export environment: %s: %s", "BITRISE_GRADLE_TEST_RESULT", err)
		}

		os.Exit(1)
	}

	if err := outputExporter.ExportOutput("BITRISE_GRADLE_TEST_RESULT", "succeeded"); err != nil {
		logger.Warnf("Failed to export environment: %s: %s", "BITRISE_GRADLE_TEST_RESULT", err)
	}

	// Collecting caches
	logger.Infof("Collecting cache:")
	const defaultProjectRoot = "."

	if err := cache.Collect(defaultProjectRoot, utilscache.Level(inputs.CacheLevel), cmdFactory); err != nil {
		logger.Warnf("Failed to collect cache: %s", err)
	}
}

func validateInputs(pathChecker pathutil.PathChecker, inputs Inputs) error {
	if inputs.GradleFile != "" {
		if exist, err := pathChecker.IsPathExists(inputs.GradleFile); err != nil {
			return fmt.Errorf("failed to check if GradleFile exist at %s: %w", inputs.GradleFile, err)
		} else if !exist {
			return fmt.Errorf("gradle file not exist at: %s", inputs.GradleFile)
		}
	}

	return nil
}

func runGradleTask(cmdFactory command.Factory, logger log.Logger, gradleTool, buildFile, tasks, options string) error {
	optionSlice, err := shellquote.Split(options)
	if err != nil {
		return err
	}

	taskSlice, err := shellquote.Split(tasks)
	if err != nil {
		return err
	}

	var args []string
	if buildFile != "" {
		args = append(args, "--build-file", buildFile)
	}
	args = append(args, taskSlice...)
	args = append(args, optionSlice...)

	cmd := cmdFactory.Create(gradleTool, args, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	logger.Donef("$ %s", cmd.PrintableCommandArgs())
	fmt.Println()

	return cmd.Run()
}

func failF(logger log.Logger, format string, args ...interface{}) {
	logger.Errorf(format, args...)
	os.Exit(1)
}
