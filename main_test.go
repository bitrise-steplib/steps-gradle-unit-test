package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_testResultName(t *testing.T) {
	tests := []struct {
		name           string
		testResultPath string
		projectRootDir string
		want           string
	}{
		{
			name:           "Absolute path with module and task",
			testResultPath: "/users/vagrant/Taskman/composeApp/build/test-results/testDebugUnitTest/TEST-io.bitrise.taskman.AppTest.xml",
			projectRootDir: "/users/vagrant/Taskman",
			want:           "composeApp-testDebugUnitTest-TEST-io.bitrise.taskman.AppTest.xml",
		},
		{
			name:           "Relative path with module and task",
			testResultPath: "./composeApp/build/test-results/testDebugUnitTest/TEST-io.bitrise.taskman.AppTest.xml",
			projectRootDir: "./",
			want:           "composeApp-testDebugUnitTest-TEST-io.bitrise.taskman.AppTest.xml",
		},
		{
			name:           "Embedded module and task",
			testResultPath: "./server/composeApp/build/test-results/testDebugUnitTest/TEST-io.bitrise.taskman.AppTest.xml",
			projectRootDir: "./",
			want:           "server-composeApp-testDebugUnitTest-TEST-io.bitrise.taskman.AppTest.xml",
		},
		{
			name:           "Relative project root dir and test result path with different syntax",
			testResultPath: "_tmp/composeApp/build/test-results/testDebugUnitTest/TEST-io.bitrise.taskman.AppTest.xml",
			projectRootDir: "./_tmp",
			want:           "composeApp-testDebugUnitTest-TEST-io.bitrise.taskman.AppTest.xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testResultName(tt.testResultPath, tt.projectRootDir)
			require.Equal(t, tt.want, got)
		})
	}
}
