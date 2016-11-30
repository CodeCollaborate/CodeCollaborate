#!/usr/bin/env bash

PACKAGES=$(go list ./... | grep -v /vendor/)
FILES=$(find . -type f -name '*.go' -not -path "./vendor/*")

RESULT_GOLINT="  PASSED: GoLint"
RESULT_GOVET="  PASSED: GoVet"
RESULT_GOFMT="  PASSED: GoFmt"
RESULT_GOIMPORTS="  PASSED: Imports"

STATUS=0

printf -- "Checking Formatting:\n--------------------------------------------------------------------------------\n"

# Check GoLint
if [[ "$(for p in $PACKAGES; do golint $p 2>&1; done)" ]]; then
    printf -- "GoLint errors:\n"
    RESULT_GOLINT="->Failed: GoLint"
    for p in $PACKAGES; do golint $p 2>&1; done
    printf -- "\n"
    STATUS=1
fi

# Check GoVet
if [[ "$(go vet $PACKAGES 2>&1)" ]]; then
    printf -- "GoVet errors:\n"
    RESULT_GOVET="->Failed: GoVet"
    go vet $PACKAGES 2>&1
    printf -- "\n"
    STATUS=1
fi

# Check GoFmt
if [[ "$(gofmt -s -l $FILES 2>&1)" ]]; then
    printf -- "GoFmt reformatting code.\n"
    gofmt -s -w $FILES 2>&1
    RESULT_GOFMT="  REFORMATTED: GoFmt"
fi

# Check GoImports
if [[ "$($GOPATH/bin/goimports -l $FILES 2>&1)" ]]; then
    printf -- "GoImports reformatting code.\n"
    $GOPATH/bin/goimports -w $FILES 2>&1
    RESULT_GOIMPORTS="  REFORMATTED: GoImports"
fi

printf -- "\nSUMMARY:\n"
printf -- "--------------------------------------------------------------------------------\n"

echo $RESULT_GOLINT
echo $RESULT_GOVET
echo $RESULT_GOFMT
echo $RESULT_GOIMPORTS

printf -- "--------------------------------------------------------------------------------\n"

exit $STATUS
