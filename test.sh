#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

go test -v ./...
go test -tags=integration -v ./...
cd api && go test -fuzztime 5s -fuzz FuzzPostGreeting
