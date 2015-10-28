#!/bin/bash
set -e

if [ ! -f "gradlew" ]; then
    echo "\e[31mDidn't find a gradlew file in the root directory\e[0m"
    exit 0
fi


./gradlew test --continue