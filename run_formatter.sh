#!/usr/bin/env bash

PACKAGES=$(go list ./... | grep -v /vendor/)
FILES=$(find . -type f -name '*.go' -not -path "./vendor/*")

RESULT_GOLINT="  PASSED: GoLint"
RESULT_GOVET="  PASSED: GoVet"
RESULT_GOFMT="  PASSED: GoFmt"
RESULT_GOIMPORTS="  PASSED: Imports"

STATUS=0

printf -- "Checking Formatting:\n"
printf -- "--------------------------------------------------------------------------------\n"

# Check GoLint
while read data; do
    if [[ ${data} ]] ; then
        printf -- "GoLint errors:\n"
        RESULT_GOLINT="->Failed: GoLint"
        echo ${data}
        printf -- "\n"
        STATUS=1
    fi
done <<< "$(for p in ${PACKAGES}; do golint ${p} 2>&1; done)"

# Check GoVet
while read data; do
    if [[ ${data} ]]; then
        printf -- "GoVet errors:\n"
        RESULT_GOVET="->Failed: GoVet"
        echo ${data}
        printf -- "\n"
        STATUS=1
    fi
done <<< "$(go vet ${PACKAGES} 2>&1)"

# Check GoFmt
while read data; do
    if [[ ${data} ]]; then
        printf -- "GoFmt reformatting code.\n"
        echo ${data}
        RESULT_GOFMT="  REFORMATTED: GoFmt"
    fi
done <<< "$(gofmt -s -l -w ${FILES} 2>&1)"

# Check GoImports
while read data; do
    if [[ ${data} ]]; then
        printf -- "GoImports reformatting code.\n"
        echo ${data}
        RESULT_GOIMPORTS="  REFORMATTED: GoImports"
    fi
done <<< "$(${GOPATH}/bin/goimports -l -w ${FILES} 2>&1)"

printf -- "\nSUMMARY:\n"
printf -- "--------------------------------------------------------------------------------\n"

echo ${RESULT_GOLINT}
echo ${RESULT_GOVET}
echo ${RESULT_GOFMT}
echo ${RESULT_GOIMPORTS}

printf -- "--------------------------------------------------------------------------------\n"

exit ${STATUS}
