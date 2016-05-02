#!/usr/bin/env bash

PACKAGES=$(go list ./... | grep -v /vendor/)
FILES=$(find . -type f -name '*.go' -not -path "./vendor/*")

STATUS=0

printf -- "Checking Formatting:\n--------------------------------------------------------------------------------\n"

# Check GoLint
if [[ "$(for p in $PACKAGES; do golint $p 2>&1; done)" ]]; then
    echo "->FAILED: GoLint - failed linting checks"
    STATUS=1
else
    echo "  PASSED: GoLint"
fi

# Check GoVet
if [[ "$(go vet $PACKAGES 2>&1)" ]]; then
    echo "->FAILED: GoVet - failed vetting checks"
    STATUS=1
else
    echo "  PASSED: GoVet"
fi

# Check GoFmt
if [[ "$(gofmt -s -l $FILES 2>&1)" ]]; then
    echo "->FAILED: GoFmt - failed formatting checks"
    STATUS=1
else
    echo "  PASSED: GoFmt"
fi

# Check GoImports
if [[ "$($GOPATH/bin/goimports -l $FILES 2>&1)" ]]; then
    echo "->FAILED: GoImports - failed imports checks"
    STATUS=1
else
    echo "  PASSED: GoImports"
fi

printf -- "--------------------------------------------------------------------------------\n"

exit $STATUS