#!/usr/bin/env bash

STATUS=0

show_failed(){
    re="\--- .*?:"
    ra="\=== RUN"
    while read data; do
        if [[ "$data" =~ $re ]] ; then
            OUTPUT=$(echo ${data} | cut -c 5-)
        else
            # catch junk
            if [[ "$data" =~ $ra ]] || [[ "$data" == ok* ]] || [[ "$data" == PASS* ]] || [[ "$data" == FAIL* ]] || [[ "$data" == \?* ]]; then
                continue
            else
                OUTPUT=$(echo ${data})
            fi
        fi


        if [[ "$OUTPUT" =~ FAIL.* ]]; then
            printf -- "\e[1;31m->%s \e[0m\n" "$OUTPUT"
            STATUS=1
        else
            printf -- "  %s\n" "$OUTPUT"
        fi
    done
}

printf -- "Running Tests:\n--------------------------------------------------------------------------------\n"
go test -v $(go list ./... | grep -v /vendor/) | show_failed
#go test -v $(go list ./... | grep -v /vendor/) | grep -E "\--- .*?:" | show_failed

printf -- "--------------------------------------------------------------------------------\n"

if [ "$STATUS" == 0 ]; then
    printf -- "Tests Passed\n"
else
    printf -- "Tests Failed\n"
fi

exit "$STATUS"
