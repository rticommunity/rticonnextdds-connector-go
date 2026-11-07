#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd "$(dirname "$0")" && pwd)
work=$(mktemp -d)
trap 'rm -f "$work/go" "$work/curl" "$work/calls"; rmdir "$work"' EXIT

cat > "$work/go" <<'EOF'
#!/usr/bin/env bash
if [ "$1" = env ]; then
    echo /tmp/mock-gopath
    exit 0
fi
printf '%s\n' "${GOTOOLCHAIN:-}" "$@" > "$GO_INSTALL_RECORD"
exit "${GO_INSTALL_STATUS:-0}"
EOF
cat > "$work/curl" <<'EOF'
#!/usr/bin/env bash
printf '#!/bin/sh\nexit 0\n'
EOF
chmod +x "$work/go" "$work/curl"

export PATH="$work:$PATH"
export GO_INSTALL_RECORD="$work/calls"

env -u LINTER_VER bash "$script_dir/install_golangci-lint.sh"
if [ ! -f "$work/calls" ] ||
    ! grep -Fxq 'go1.26.10' "$work/calls" ||
    ! grep -Fxq 'install' "$work/calls" ||
    ! grep -Fxq 'github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0' "$work/calls"; then
    echo 'FAIL: installer must build pinned v2.15.0 with patched Go 1.26.10' >&2
    exit 1
fi

LINTER_VER=v2.14.0 bash "$script_dir/install_golangci-lint.sh"
grep -Fxq 'github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0' "$work/calls"

if GO_INSTALL_STATUS=42 bash "$script_dir/install_golangci-lint.sh"; then
    echo 'FAIL: installer must propagate Go installation failure' >&2
    exit 1
else
    status=$?
    test "$status" -eq 42
fi

echo 'PASS: pinned linter/toolchain, explicit override, and installation failure propagation'
