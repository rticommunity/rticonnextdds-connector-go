#!/usr/bin/env bash
set -euo pipefail

LINTER_VER=${LINTER_VER:-v2.15.0}

GOTOOLCHAIN=go1.26.10 go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${LINTER_VER}"