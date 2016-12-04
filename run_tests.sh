#!/usr/bin/env bash

FAILED=0

printf -- "Running Tests:\n"
printf -- "--------------------------------------------------------------------------------\n"

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
        FAILED=1
    else
        printf -- "  %s\n" "$OUTPUT"
    fi
done < <(go test -v $(go list ./... | grep -v /vendor/))

printf -- "--------------------------------------------------------------------------------\n"

[ "$FAILED" == 0 ] && printf -- "\e[1;32mTests Passed\e[0m\n" || printf -- "\e[1;31mTests Failed\e[0m\n"
printf -- "\n"
exit "$FAILED"
