# Run Gradle Tests

[![Step changelog](https://shields.io/github/v/release/bitrise-io/steps-gradle-unit-test?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-io/steps-gradle-unit-test/releases)

Runs tests with `gradlew`.

<details>
<summary>Description</summary>

This Step runs tests with `gradlew`.
You can specify the test tasks to run and set flags to the executed `gradlew` command.

### Configuring the Step

To use this Step, you need at least two things:
* [Gradle Wrapper](https://docs.gradle.org/current/userguide/gradle_wrapper.html).
* A Gradle task that is correctly configured in your Gradle project.

To configure the Step:
1. Set the Gradle project root directory in the **Project root directory** input.
1. Add the task you want to run in the **Test task** input.
1. (Optional) Add any flags you want to pass to the executed gradlew command in the **Additional flags** input. For example, you can use `--tests='*.MyTestClass'` to run a specific test class.
1. (Optional) You can set the file path to a `build.gradle` file for the Step in the **Path to the Gradle build script to use** input.

### Troubleshooting

If you receive an error that Gradle Wrapper (gradlew) is required, make sure to generate one if you don't have one already. You can read
more about it in the [official Gradle guide](https://docs.gradle.org/current/userguide/gradle_wrapper.html).

### Useful links

* [Gradle Wrapper](https://docs.gradle.org/current/userguide/gradle_wrapper.html)
* [Caching Gradle](https://devcenter.bitrise.io/builds/caching/caching-gradle/)
* [Running a unit test](https://devcenter.bitrise.io/en/testing/android-unit-tests.html#running-a-unit-test)

### Related Steps

* [Generate Gradle Wrapper](https://www.bitrise.io/integrations/steps/generate-gradle-wrapper)
* [Gradle Runner](https://www.bitrise.io/integrations/steps/gradle-runner)
* [Android Build](https://www.bitrise.io/integrations/steps/android-build)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `project_root_dir` | The root directory of the Gradle project. This is the directory which contains all source files from your project, as well as Gradle files, including the `gradlew` file. | required | `$BITRISE_SOURCE_DIR` |
| `test_task` | The test task to be executed. | required | `test` |
| `gradlew_command_flags` | Flags to pass to the executed gradlew command. For example, you can use `--tests='*.MyTestClass'` to run a specific test class. |  | `--continue` |
| `gradle_build_script_path` | Path to the Gradle build script to use. The path should be relative to the **Project root directory**. |  |  |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `BITRISE_GRADLE_TEST_RESULT` | The result of the tests executed by this Step. |
| `BITRISE_FLAKY_TEST_CASES` | A test case is considered flaky if it has failed at least once, but passed at least once as well.  The list contains the test cases in the following format: ``` - TestSuit_1.TestClass_1.TestName_1 - TestSuit_1.TestClass_1.TestName_2 - TestSuit_1.TestClass_2.TestName_1 - TestSuit_2.TestClass_1.TestName_1 ... ``` |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-io/steps-gradle-unit-test/pulls) and [issues](https://github.com/bitrise-io/steps-gradle-unit-test/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
