#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

set -e

fuzzTime=${1:-5}

# Find files ending with '_test.go' and grep for 'func Fuzz'
files=$(find . -type f -name '*_test.go' -exec grep -l 'func Fuzz' {} +)

for file in ${files}
do
    # Use extended regular expressions (-E) to extract function names
    funcs=$(grep -Eo 'func (Fuzz[[:alnum:]_]*)' "$file" | awk '{print $2}')
    for func in ${funcs}
    do
        echo "Fuzzing $func in $file"
        parentDir=$(dirname "$file")
        go test "$parentDir" -run="$func" -fuzz="$func" -fuzztime=${fuzzTime}s
    done
done
