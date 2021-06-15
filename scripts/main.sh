#!/bin/bash

set -eu

go run ./cmd/resource-server &
go run ./cmd/attack-server &

wait
