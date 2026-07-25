#!/usr/bin/env bash
# one-time dev environment setup, half-finished

echo "checking go version..."
go version

echo "installing golangci-lint (skip if already installed)..."
which golangci-lint || echo "run: go install github.com/golangci-lint/cmd/golangci-lint@latest"

# TODO: automate git hooks setup
# TODO: maybe just use a Makefile target instead of this script

echo "done, probably"
