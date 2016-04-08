#!/usr/bin/env bash

STATUS=0

printf -- "Checking Formatting:\n--------------------------------------------------------------------------------\n"

# Check GoLint
if [[ "$($GOPATH/bin/golint ./... 2>&1)" ]]; then
    echo "->FAILED: GoLint - failed linting checks"
    STATUS=1
else
    echo "  PASSED: GoLint"
fi

# Check GoVet
if [[ "$(go vet ./... 2>&1)" ]]; then
    echo "->FAILED: GoVet - failed vetting checks"
    STATUS=1
else
    echo "  PASSED: GoVet"
fi

# Check GoFmt
if [[ "$(gofmt -s -l . 2>&1)" ]]; then
    echo "->FAILED: GoFmt - failed formatting checks"
    STATUS=1
else
    echo "  PASSED: GoFmt"
fi

# Check GoImports
if [[ "$($GOPATH/bin/goimports -l . 2>&1)" ]]; then
    echo "->FAILED: GoImports - failed imports checks"
    STATUS=1
else
    echo "  PASSED: GoImports"
fi

printf -- "--------------------------------------------------------------------------------\n"

exit $STATUS