#!/usr/bin/env bash

show_failed(){
    STATUS=0

    while read data; do
        OUTPUT=$(echo $data | cut -c 5-)

        if [[ "$OUTPUT" =~ PASS.* ]]; then
            printf -- "  %s\n" "$OUTPUT"
        else
            printf -- "->%s\n" "$OUTPUT"
            STATUS=1
        fi
    done

    return "$STATUS"
}

printf -- "Running Tests:\n--------------------------------------------------------------------------------\n"
go test -v $(go list ./... | grep -v /vendor/) | grep -E "\--- .*?:" | show_failed
RESULT=$?

printf -- "--------------------------------------------------------------------------------\n"

if [ "$RESULT" == 0 ]; then
    printf -- "Tests Passed\n"
else
    printf -- "Tests Failed\n"
fi

exit "$RESULT"