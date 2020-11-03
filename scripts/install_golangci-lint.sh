#!/usr/bin/env bash

declare LINTER_VER=${LINTER_VER:-'v1.31.0'}
declare GO_ENV_PATH=$(go env GOPATH)

declare -r INSTALL_SCRIPT="https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh"

curl -sfL ${INSTALL_SCRIPT} | sh -s -- -b "${GO_ENV_PATH}/bin" "${LINT_VER}"