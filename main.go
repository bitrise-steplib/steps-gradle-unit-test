package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	ProjectRootDir        string `env:"project_root_dir,dir"`
	TestTasks             string `env:"test_task,required"`
	GradlewCommandFlags   string `env:"gradlew_command_flags"`
	GradleBuildScriptPath string `env:"gradle_build_script_path"`
	CacheLevel            string `env:"cache_level,opt[all,only deps,none]"`
}

type Configs struct {
	ProjectRootDir        string
	GradlewPath           string
	TestTasks             string
	GradlewCommandFlags   string
	GradleBuildScriptPath string
	CacheLevel            string
}

func main() {
	// Setup dependencies
	logger := log.NewLogger()
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)
	pathChecker := pathutil.NewPathChecker()
	cmdFactory := command.NewFactory(envRepo)
	outputExporter := export.NewExporter(cmdFactory)

	// Parse inputs
	config, err := processConfig(inputParser, pathChecker, logger)
	if err != nil {
		failF(logger, "Failed to process config: %s", err)
	}

	// Run gradle task
	if err := os.Chmod(config.GradlewPath, 0770); err != nil {
		failF(logger, "Failed to add executable permission on gradlew file (%s): %s", config.GradlewPath, err)
	}

	fmt.Println()
	logger.Infof("Running gradle task...")
	if err := runGradleTask(cmdFactory, logger, config.ProjectRootDir, config.GradlewPath, config.TestTasks, config.GradleBuildScriptPath, config.GradlewCommandFlags); err != nil {
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

	if err := cache.Collect(defaultProjectRoot, utilscache.Level(config.CacheLevel), cmdFactory); err != nil {
		logger.Warnf("Failed to collect cache: %s", err)
	}
}

func processConfig(inputParser stepconf.InputParser, pathChecker pathutil.PathChecker, logger log.Logger) (*Configs, error) {
	var inputs Inputs
	if err := inputParser.Parse(&inputs); err != nil {
		return nil, fmt.Errorf("issue with input: %s", err)
	}

	stepconf.Print(inputs)
	logger.Println()

	if inputs.GradleBuildScriptPath != "" {
		if exist, err := pathChecker.IsPathExists(inputs.GradleBuildScriptPath); err != nil {
			return nil, fmt.Errorf("failed to check if gradle build file exist at %s: %w", inputs.GradleBuildScriptPath, err)
		} else if !exist {
			return nil, fmt.Errorf("gradle build file not exist at: %s", inputs.GradleBuildScriptPath)
		}
	}

	gradlewPath := filepath.Join(inputs.ProjectRootDir, "gradlew")
	if exist, err := pathChecker.IsPathExists(gradlewPath); err != nil {
		return nil, fmt.Errorf("failed to check if gradlew exist at %s: %w", gradlewPath, err)
	} else if !exist {
		return nil, fmt.Errorf("gradlew file not exist at: %s", gradlewPath)
	}

	return &Configs{
		ProjectRootDir:        inputs.ProjectRootDir,
		GradlewPath:           gradlewPath,
		TestTasks:             inputs.TestTasks,
		GradlewCommandFlags:   inputs.GradlewCommandFlags,
		GradleBuildScriptPath: inputs.GradleBuildScriptPath,
		CacheLevel:            inputs.CacheLevel,
	}, nil
}

func runGradleTask(cmdFactory command.Factory, logger log.Logger, workDir, gradlewPth, tasks, buildScriptPth, flags string) error {
	flagSlice, err := shellquote.Split(flags)
	if err != nil {
		return err
	}

	taskSlice, err := shellquote.Split(tasks)
	if err != nil {
		return err
	}

	var args []string
	if buildScriptPth != "" {
		args = append(args, "--build-file", buildScriptPth)
	}
	args = append(args, taskSlice...)
	args = append(args, flagSlice...)

	cmd := cmdFactory.Create(gradlewPth, args, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    workDir,
	})

	logger.Donef("$ %s", cmd.PrintableCommandArgs())
	fmt.Println()

	return cmd.Run()
}

func failF(logger log.Logger, format string, args ...interface{}) {
	logger.Errorf(format, args...)
	os.Exit(1)
}
