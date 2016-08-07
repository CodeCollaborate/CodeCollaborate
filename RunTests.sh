#!/usr/bin/env bash

show_failed(){
    STATUS=0

    while read data; do
        OUTPUT=$(echo $data | cut -c 5-)


        if [[ "$OUTPUT" =~ FAIL.* ]]; then
            printf -- "\e[1;31m->%s \e[0m\n" "$OUTPUT"
            STATUS=1
        else
            printf -- "  %s\n" "$OUTPUT"
        fi
    done

    return "$STATUS"
}

printf -- "Running Tests:\n--------------------------------------------------------------------------------\n"
go test -v $(go list ./... | grep -v /vendor/) | grep -E "\--- .*?:" | show_failed
RESULT=${PIPESTATUS[0]}

printf -- "--------------------------------------------------------------------------------\n"

if [ "$RESULT" == 0 ]; then
    printf -- "Tests Passed\n"
else
    printf -- "Tests Failed\n"
fi

exit "$RESULT"
