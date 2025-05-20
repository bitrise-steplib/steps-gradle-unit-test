package main

import (
	"fmt"
	"os"
	"path/filepath"

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
}

type Configs struct {
	ProjectRootDir                string
	GradlewPath                   string
	TestTasks                     []string
	GradlewCommandFlags           []string
	GradleBuildScriptRelativePath string
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
	if err := runGradleTask(cmdFactory, logger, config.ProjectRootDir, config.TestTasks, config.GradleBuildScriptRelativePath, config.GradlewCommandFlags); err != nil {
		logger.Errorf("Gradle task failed: %s", err)

		if err := outputExporter.ExportOutput("BITRISE_GRADLE_TEST_RESULT", "failed"); err != nil {
			logger.Warnf("Failed to export environment: %s: %s", "BITRISE_GRADLE_TEST_RESULT", err)
		}

		os.Exit(1)
	}

	if err := outputExporter.ExportOutput("BITRISE_GRADLE_TEST_RESULT", "succeeded"); err != nil {
		logger.Warnf("Failed to export environment: %s: %s", "BITRISE_GRADLE_TEST_RESULT", err)
	}
}

func processConfig(inputParser stepconf.InputParser, pathChecker pathutil.PathChecker, logger log.Logger) (*Configs, error) {
	var inputs Inputs
	if err := inputParser.Parse(&inputs); err != nil {
		return nil, fmt.Errorf("issue with input: %s", err)
	}

	stepconf.Print(inputs)
	logger.Println()

	var gradleBuildScriptPath string
	if inputs.GradleBuildScriptPath != "" {
		gradleBuildScriptPath = filepath.Join(inputs.ProjectRootDir, inputs.GradleBuildScriptPath)
		if exist, err := pathChecker.IsPathExists(gradleBuildScriptPath); err != nil {
			return nil, fmt.Errorf("failed to check if gradle build file exist at %s: %w", gradleBuildScriptPath, err)
		} else if !exist {
			return nil, fmt.Errorf("gradle build file not exist at: %s", gradleBuildScriptPath)
		}
	}

	gradlewPath := filepath.Join(inputs.ProjectRootDir, "gradlew")
	if exist, err := pathChecker.IsPathExists(gradlewPath); err != nil {
		return nil, fmt.Errorf("failed to check if gradlew exist at %s: %w", gradlewPath, err)
	} else if !exist {
		return nil, fmt.Errorf("gradlew file not exist at: %s", gradlewPath)
	}

	taskSlice, err := shellquote.Split(inputs.TestTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	flagSlice, err := shellquote.Split(inputs.GradlewCommandFlags)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gradlew flags: %w", err)
	}

	return &Configs{
		ProjectRootDir:                inputs.ProjectRootDir,
		GradlewPath:                   gradlewPath,
		TestTasks:                     taskSlice,
		GradlewCommandFlags:           flagSlice,
		GradleBuildScriptRelativePath: inputs.GradleBuildScriptPath,
	}, nil
}

func runGradleTask(cmdFactory command.Factory, logger log.Logger, workDir string, tasks []string, buildScriptPth string, flags []string) error {
	var args []string
	if buildScriptPth != "" {
		args = append(args, "--build-file", buildScriptPth)
	}
	args = append(args, tasks...)
	args = append(args, flags...)

	cmd := cmdFactory.Create("./gradlew", args, &command.Opts{
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
