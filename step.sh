#!/bin/bash

if [ ! -f "gradlew" ] ; then
    echo "\e[31mDidn't find a gradlew file in the root directory\e[0m"
    exit 0
fi

./gradlew ${unit_test_task} ${unit_test_flags}
return_code=$?

if [ "${return_code}" -eq "0" ] ; then
	envman add --key "BITRISE_GRADLE_TEST_RESULT" --value "succeeded"
else
	envman add --key "BITRISE_GRADLE_TEST_RESULT" --value "failed"
fi

echo "BITRISE_GRADLE_TEST_RESULT added to the environment with value ${return_code}"
exit ${return_code}

