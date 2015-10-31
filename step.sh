#!/bin/bash

if [ -z "${gradlew_file_path}" ] ; then
	printf "\e[31gradlew_file_path was not defined\e[0m\n"
    exit 1
fi

if [ ! -f "${gradlew_file_path}" ] ; then
    printf "\e[31mDidn't find a gradlew file in the root directory\e[0m\n"
    exit 1
fi

if [ -z "${unit_test_task}" ] ; then
    printf "\e[31munit_test_task was not defined\e[0m\n"
    exit 1
fi

echo "$" ${gradlew_file_path} ${unit_test_task} ${unit_test_flags}
${gradlew_file_path} ${unit_test_task} ${unit_test_flags}
return_code=$?

if [ "${return_code}" -eq "0" ] ; then
	envman add --key "BITRISE_GRADLE_TEST_RESULT" --value "succeeded"
	echo "BITRISE_GRADLE_TEST_RESULT added to the environment with value succeeded"
else
	envman add --key "BITRISE_GRADLE_TEST_RESULT" --value "failed"
	echo "BITRISE_GRADLE_TEST_RESULT added to the environment with value failed"
fi

exit ${return_code}

