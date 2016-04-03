#!/usr/bin/env bash

STATUS=0

# Check GoLint
if [[ "$(golint ./...)" ]]; then
    echo "FAILED: GoLint - failed linting checks"
    STATUS=1
else
    echo "PASSED: GoLint"
fi

# Check GoFmt
if [[ "$(gofmt -s -l .)" ]]; then
    echo "FAILED: GoFmt - failed formatting checks"
    STATUS=1
else
    echo "PASSED: GoFmt"
fi

exit $STATUS