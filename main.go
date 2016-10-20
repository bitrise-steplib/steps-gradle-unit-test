package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	shellquote "github.com/kballard/go-shellquote"
)

// ConfigsModel ...
type ConfigsModel struct {
	GradleFile    string
	UnitTestTasks string
	GradlewPath   string
	UnitTestFlags string

	DeployDir string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		GradleFile:    os.Getenv("gradle_file"),
		UnitTestTasks: os.Getenv("unit_test_task"),
		GradlewPath:   os.Getenv("gradlew_file_path"),
		UnitTestFlags: os.Getenv("unit_test_flags"),

		DeployDir: os.Getenv("BITRISE_DEPLOY_DIR"),
	}
}

func (configs ConfigsModel) print() {
	log.Info("Configs:")
	log.Detail("- GradleFile: %s", configs.GradleFile)
	log.Detail("- UnitTestTasks: %s", configs.UnitTestTasks)
	log.Detail("- GradlewPath: %s", configs.GradlewPath)
	log.Detail("- UnitTestFlags: %s", configs.UnitTestFlags)

	log.Detail("- DeployDir: %s", configs.DeployDir)
}

func (configs ConfigsModel) validate() (string, error) {
	// required
	if configs.GradleFile == "" {
		return "", errors.New("No GradleFile parameter specified!")
	}
	if exist, err := pathutil.IsPathExists(configs.GradleFile); err != nil {
		return "", fmt.Errorf("Failed to check if GradleFile exist at: %s, error: %s", configs.GradleFile, err)
	} else if !exist {
		return "", fmt.Errorf("GradleFile not exist at: %s", configs.GradleFile)
	}

	if configs.UnitTestTasks == "" {
		return "", errors.New("No UnitTestTasks parameter specified!")
	}

	if configs.GradlewPath == "" {
		explanation := `
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
that the right Gradle version is installed and used for the build.

You can find more information about the Gradle Wrapper (gradlew),
and about how you can generate one (if you would not have one already
in the official guide at: https://docs.gradle.org/current/userguide/gradle_wrapper.html`

		return explanation, errors.New("No GradlewPath parameter specified!")
	}
	if exist, err := pathutil.IsPathExists(configs.GradlewPath); err != nil {
		return "", fmt.Errorf("Failed to check if GradlewPath exist at: %s, error: %s", configs.GradlewPath, err)
	} else if !exist {
		return "", fmt.Errorf("GradlewPath not exist at: %s", configs.GradlewPath)
	}

	return "", nil
}

func runGradleTask(gradleTool, buildFile, tasks, options string) error {
	optionSlice, err := shellquote.Split(options)
	if err != nil {
		return err
	}

	taskSlice, err := shellquote.Split(tasks)
	if err != nil {
		return err
	}

	cmdSlice := []string{gradleTool, "--build-file", buildFile}
	cmdSlice = append(cmdSlice, taskSlice...)
	cmdSlice = append(cmdSlice, optionSlice...)

	log.Done("$ %s", cmdex.PrintableCommandArgs(false, cmdSlice))
	fmt.Println()

	cmd, err := cmdex.NewCommandFromSlice(cmdSlice)
	if err != nil {
		return fmt.Errorf("failed to create command, error: %s", err)
	}

	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	return cmd.Run()
}

func exportEnvironmentWithEnvman(keyStr, valueStr string) error {
	cmd := cmdex.NewCommand("envman", "add", "--key", keyStr)
	cmd.SetStdin(strings.NewReader(valueStr))
	return cmd.Run()
}

func main() {
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()

	if explanation, err := configs.validate(); err != nil {
		fmt.Println()
		log.Error("Issue with input: %s", err)
		fmt.Println()

		if explanation != "" {
			fmt.Println(explanation)
			fmt.Println()
		}

		os.Exit(1)
	}

	err := os.Chmod(configs.GradlewPath, 0770)
	if err != nil {
		log.Error("Failed to add executable permission on gradlew file (%s), error: %s", configs.GradlewPath, err)
		os.Exit(1)
	}

	fmt.Println()
	log.Info("Running gradle task...")
	if err := runGradleTask(configs.GradlewPath, configs.GradleFile, configs.UnitTestTasks, configs.UnitTestFlags); err != nil {
		log.Error("Gradle task failed, error: %s", err)

		if err := exportEnvironmentWithEnvman("BITRISE_GRADLE_TEST_RESULT", "failed"); err != nil {
			log.Warn("Failed to export environment: %s, error: %s", "BITRISE_GRADLE_TEST_RESULT", err)
		}

		os.Exit(1)
	}

	if err := exportEnvironmentWithEnvman("BITRISE_GRADLE_TEST_RESULT", "succeeded"); err != nil {
		log.Warn("Failed to export environment: %s, error: %s", "BITRISE_GRADLE_TEST_RESULT", err)
	}
}
