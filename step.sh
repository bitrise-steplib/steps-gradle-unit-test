#!/bin/bash
if [ ! -f "gradlew" ]; then
    gradle build connectedCheck
else
	./gradlew build connectedCheck
fi

if [ "$?" = "0" ]; then
	envman aadd --key BITRISE_GRADLE_TEST_RESULT --value "succeeded"
	exit 0
else
	envman aadd --key BITRISE_GRADLE_TEST_RESULT --value "failed"
	exit 1
fi